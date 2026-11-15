package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/question"
)

const maxBodyBytes = 1 << 16

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, question.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, question.ErrNotFound):
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
		return errors.Join(question.ErrInvalid, errors.New("malformed JSON body"))
	}
	if dec.More() {
		return errors.Join(question.ErrInvalid, errors.New("unexpected data after JSON body"))
	}
	return nil
}
