// Package httpx is the JSON HTTP layer of the API.
package httpx

import (
	"context"
	"net/http"

	"latihseleksi/api/internal/auth"
	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/insight"
)

// Config holds the HTTP settings read from the environment.
type Config struct {
	// CORSOrigin is the single allowed origin; empty disables CORS headers.
	CORSOrigin string
	// Auth verifies ID tokens for /v1/attempts*, /v1/exams* and /v1/me. Nil
	// means auth is not configured: those routes answer 503.
	Auth *auth.Verifier
}

type api struct {
	svc     *drill.Service
	exam    *exam.Service
	insight *insight.Service
	auth    *auth.Verifier
	users   userStore
}

// userStore maps an auth subject to a local user id.
type userStore interface {
	UpsertUser(ctx context.Context, sub string) (string, error)
	// DeleteUser removes the user and everything it owns.
	DeleteUser(ctx context.Context, id string) error
}

// NewHandler builds the router with request id, access log and CORS middleware.
func NewHandler(svc *drill.Service, ex *exam.Service, in *insight.Service, users userStore, cfg Config) http.Handler {
	a := &api{svc: svc, exam: ex, insight: in, auth: cfg.Auth, users: users}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/patterns", a.patterns)
	mux.HandleFunc("POST /v1/attempts", a.requireAuth(a.createAttempt))
	mux.HandleFunc("GET /v1/attempts/{id}", a.requireAuth(a.getAttempt))
	mux.HandleFunc("PUT /v1/attempts/{id}/responses/{pos}", a.requireAuth(a.putResponse))
	mux.HandleFunc("POST /v1/attempts/{id}/submit", a.requireAuth(a.submit))
	mux.HandleFunc("GET /v1/specs", a.specs)
	mux.HandleFunc("POST /v1/exams", a.requireAuth(a.createExam))
	mux.HandleFunc("GET /v1/exams", a.requireAuth(a.listExams))
	mux.HandleFunc("GET /v1/exams/{id}", a.requireAuth(a.getExam))
	mux.HandleFunc("PUT /v1/exams/{id}/responses/{pos}", a.requireAuth(a.putExamResponse))
	mux.HandleFunc("POST /v1/exams/{id}/submit", a.requireAuth(a.submitExam))
	mux.HandleFunc("GET /v1/exams/{id}/result", a.requireAuth(a.examResult))
	mux.HandleFunc("POST /v1/exams/{id}/items/{pos}/report", a.requireAuth(a.reportItem))
	mux.HandleFunc("GET /v1/me/dashboard", a.requireAuth(a.dashboard))
	mux.HandleFunc("DELETE /v1/me", a.requireAuth(a.deleteMe))
	return requestID(cors(cfg.CORSOrigin, mux))
}
