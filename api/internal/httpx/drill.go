package httpx

import (
	"net/http"
	"strconv"

	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/question"
)

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
		writeError(w, r, question.ErrNotFound)
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
