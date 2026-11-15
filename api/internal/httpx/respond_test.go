package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/question"
)

func TestWriteError_status(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"shared invalid", fmt.Errorf("%w: bad", question.ErrInvalid), http.StatusBadRequest},
		{"exam invalid", fmt.Errorf("%w: bad", exam.ErrInvalid), http.StatusBadRequest},
		{"drill not found", drill.ErrNotFound, http.StatusNotFound},
		{"shared not found", question.ErrNotFound, http.StatusNotFound},
		{"submitted", drill.ErrSubmitted, http.StatusConflict},
		{"closed", exam.ErrClosed, http.StatusConflict},
		{"generator failure", fmt.Errorf("pick items: %w", errors.New("boom")), http.StatusInternalServerError},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		writeError(rec, httptest.NewRequest(http.MethodGet, "/", nil), c.err)
		if rec.Code != c.want {
			t.Errorf("%s: status %d, want %d", c.name, rec.Code, c.want)
		}
	}
}
