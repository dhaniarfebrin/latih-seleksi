// Package httpx is the JSON HTTP layer of the API.
package httpx

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"latihseleksi/api/internal/auth"
	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
)

const maxBodyBytes = 1 << 16

type ctxKey struct{}

// Config holds the HTTP settings read from the environment.
type Config struct {
	// CORSOrigin is the single allowed origin; empty disables CORS headers.
	CORSOrigin string
	// Auth verifies ID tokens for /v1/exams*. Nil means auth is not
	// configured: those routes answer 503 and drill routes keep working.
	Auth *auth.Verifier
}

type api struct {
	svc   *drill.Service
	exam  *exam.Service
	auth  *auth.Verifier
	users userStore
}

// userStore maps an auth subject to a local user id.
type userStore interface {
	UpsertUser(ctx context.Context, sub string) (string, error)
}

// NewHandler builds the router with request id, access log and CORS middleware.
func NewHandler(svc *drill.Service, ex *exam.Service, users userStore, cfg Config) http.Handler {
	a := &api{svc: svc, exam: ex, auth: cfg.Auth, users: users}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/patterns", a.patterns)
	mux.HandleFunc("POST /v1/attempts", a.createAttempt)
	mux.HandleFunc("GET /v1/attempts/{id}", a.getAttempt)
	mux.HandleFunc("PUT /v1/attempts/{id}/responses/{pos}", a.putResponse)
	mux.HandleFunc("POST /v1/attempts/{id}/submit", a.submit)
	mux.HandleFunc("GET /v1/specs", a.specs)
	mux.HandleFunc("POST /v1/exams", a.requireAuth(a.createExam))
	mux.HandleFunc("GET /v1/exams", a.requireAuth(a.listExams))
	mux.HandleFunc("GET /v1/exams/{id}", a.requireAuth(a.getExam))
	mux.HandleFunc("PUT /v1/exams/{id}/responses/{pos}", a.requireAuth(a.putExamResponse))
	mux.HandleFunc("POST /v1/exams/{id}/submit", a.requireAuth(a.submitExam))
	mux.HandleFunc("GET /v1/exams/{id}/result", a.requireAuth(a.examResult))
	return requestID(cors(cfg.CORSOrigin, mux))
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b [8]byte
		_, _ = rand.Read(b[:])
		id := hex.EncodeToString(b[:])
		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), ctxKey{}, id)
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(sw, r.WithContext(ctx))
		slog.Info("request", "id", id, "method", r.Method, "path", r.URL.Path, "status", sw.status, "ms", time.Since(start).Milliseconds())
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func cors(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin != "" {
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Add("Vary", "Origin")
			if r.Method == http.MethodOptions {
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
				h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				h.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// requireAuth runs next for a caller with a valid Bearer ID token, passing the
// local user id. The token itself is never logged.
func (a *api) requireAuth(next func(w http.ResponseWriter, r *http.Request, userID string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.auth == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
			return
		}
		scheme, token, _ := strings.Cut(r.Header.Get("Authorization"), " ")
		token = strings.TrimSpace(token)
		if !strings.EqualFold(scheme, "Bearer") || token == "" {
			unauthorized(w)
			return
		}
		claims, err := a.auth.Verify(r.Context(), token)
		switch {
		case errors.Is(err, auth.ErrEmailNotVerified):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "email not verified"})
			return
		case errors.Is(err, auth.ErrKeysUnavailable):
			id, _ := r.Context().Value(ctxKey{}).(string)
			slog.Error("auth keys unavailable", "id", id, "err", err)
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth unavailable"})
			return
		case err != nil:
			unauthorized(w)
			return
		}
		userID, err := a.users.UpsertUser(r.Context(), claims.Sub)
		if err != nil {
			writeError(w, r, err)
			return
		}
		next(w, r, userID)
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, drill.ErrInvalid), errors.Is(err, exam.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, drill.ErrNotFound), errors.Is(err, exam.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, drill.ErrSubmitted), errors.Is(err, exam.ErrClosed), errors.Is(err, exam.ErrInProgress):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	default:
		id, _ := r.Context().Value(ctxKey{}).(string)
		slog.Error("internal error", "id", id, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

// decode reads a JSON body capped at maxBodyBytes; unknown fields are rejected.
// An empty body leaves dst untouched.
func decode(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.Join(drill.ErrInvalid, errors.New("malformed JSON body"))
	}
	if dec.More() {
		return errors.Join(drill.ErrInvalid, errors.New("unexpected data after JSON body"))
	}
	return nil
}

func (a *api) patterns(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"patterns": a.svc.Patterns()})
}

func (a *api) createAttempt(w http.ResponseWriter, r *http.Request) {
	var body struct {
		PatternID string `json:"pattern_id"`
		N         *int   `json:"n"`
		Band      *int   `json:"band"`
	}
	if err := decode(w, r, &body); err != nil {
		writeError(w, r, err)
		return
	}
	n, band := drill.DefaultN, drill.DefaultBand
	if body.N != nil {
		n = *body.N
	}
	if body.Band != nil {
		band = *body.Band
	}
	v, err := a.svc.Create(r.Context(), body.PatternID, n, band)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Location", "/v1/attempts/"+v.ID)
	writeJSON(w, http.StatusCreated, v)
}

func (a *api) getAttempt(w http.ResponseWriter, r *http.Request) {
	v, err := a.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) putResponse(w http.ResponseWriter, r *http.Request) {
	pos, err := strconv.Atoi(r.PathValue("pos"))
	if err != nil {
		writeError(w, r, drill.ErrNotFound)
		return
	}
	var body struct {
		Answer string `json:"answer"`
		TimeMS int    `json:"time_ms"`
	}
	if err := decode(w, r, &body); err != nil {
		writeError(w, r, err)
		return
	}
	v, err := a.svc.Answer(r.Context(), r.PathValue("id"), pos, body.Answer, body.TimeMS)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) submit(w http.ResponseWriter, r *http.Request) {
	res, err := a.svc.Submit(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *api) specs(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"specs": a.exam.Specs()})
}

func (a *api) createExam(w http.ResponseWriter, r *http.Request, userID string) {
	var body struct {
		SpecID   string   `json:"spec_id"`
		Track    string   `json:"track"`
		Sections []string `json:"sections"`
	}
	if err := decode(w, r, &body); err != nil {
		writeError(w, r, err)
		return
	}
	v, err := a.exam.Create(r.Context(), userID, body.SpecID, body.Track, body.Sections)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Location", "/v1/exams/"+v.ID)
	writeJSON(w, http.StatusCreated, v)
}

func (a *api) listExams(w http.ResponseWriter, r *http.Request, userID string) {
	limit := exam.DefaultListLimit
	if q := r.URL.Query().Get("limit"); q != "" {
		n, err := strconv.Atoi(q)
		if err != nil {
			writeError(w, r, errors.Join(exam.ErrInvalid, errors.New("limit must be a number")))
			return
		}
		limit = n
	}
	list, err := a.exam.List(r.Context(), userID, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"exams": list})
}

func (a *api) getExam(w http.ResponseWriter, r *http.Request, userID string) {
	v, err := a.exam.Get(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) putExamResponse(w http.ResponseWriter, r *http.Request, userID string) {
	pos, err := strconv.Atoi(r.PathValue("pos"))
	if err != nil {
		writeError(w, r, exam.ErrNotFound)
		return
	}
	var body struct {
		Answer  string `json:"answer"`
		TimeMS  int    `json:"time_ms"`
		Flagged bool   `json:"flagged"`
		Seq     int    `json:"seq"`
	}
	if err := decode(w, r, &body); err != nil {
		writeError(w, r, err)
		return
	}
	v, err := a.exam.Answer(r.Context(), userID, r.PathValue("id"), pos, body.Answer, body.TimeMS, body.Flagged, body.Seq)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) submitExam(w http.ResponseWriter, r *http.Request, userID string) {
	res, err := a.exam.Submit(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *api) examResult(w http.ResponseWriter, r *http.Request, userID string) {
	res, err := a.exam.Result(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}
