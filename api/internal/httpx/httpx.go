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
	"time"

	"latihseleksi/api/internal/drill"
)

const maxBodyBytes = 1 << 16

type ctxKey struct{}

// Config holds the HTTP settings read from the environment.
type Config struct {
	// CORSOrigin is the single allowed origin; empty disables CORS headers.
	CORSOrigin string
}

type api struct {
	svc *drill.Service
}

// NewHandler builds the router with request id, access log and CORS middleware.
func NewHandler(svc *drill.Service, cfg Config) http.Handler {
	a := &api{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/patterns", a.patterns)
	mux.HandleFunc("POST /v1/attempts", a.createAttempt)
	mux.HandleFunc("GET /v1/attempts/{id}", a.getAttempt)
	mux.HandleFunc("PUT /v1/attempts/{id}/responses/{pos}", a.putResponse)
	mux.HandleFunc("POST /v1/attempts/{id}/submit", a.submit)
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
				h.Set("Access-Control-Allow-Headers", "Content-Type")
				h.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, drill.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, drill.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, drill.ErrSubmitted):
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
