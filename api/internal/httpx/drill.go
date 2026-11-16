package httpx

import (
	"net/http"
	"strconv"

	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/question"
)

// deleteMe removes the caller and all their attempts. Repeating it with the
// same token recreates the user through requireAuth and deletes it again.
func (a *api) deleteMe(w http.ResponseWriter, r *http.Request, userID string) {
	if err := a.users.DeleteUser(r.Context(), userID); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *api) patterns(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"patterns": a.svc.Patterns()})
}

func (a *api) createAttempt(w http.ResponseWriter, r *http.Request, userID string) {
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
	v, err := a.svc.Create(r.Context(), userID, body.PatternID, n, band)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Location", "/v1/attempts/"+v.ID)
	writeJSON(w, http.StatusCreated, v)
}

func (a *api) getAttempt(w http.ResponseWriter, r *http.Request, userID string) {
	v, err := a.svc.Get(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) putResponse(w http.ResponseWriter, r *http.Request, userID string) {
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
	v, err := a.svc.Answer(r.Context(), userID, r.PathValue("id"), pos, body.Answer, body.TimeMS)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) submit(w http.ResponseWriter, r *http.Request, userID string) {
	res, err := a.svc.Submit(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *api) reportAttemptItem(w http.ResponseWriter, r *http.Request, userID string) {
	pos, err := strconv.Atoi(r.PathValue("pos"))
	if err != nil {
		writeError(w, r, question.ErrNotFound)
		return
	}
	var body struct {
		Reason string `json:"reason"`
		Note   string `json:"note"`
	}
	if err := decode(w, r, &body); err != nil {
		writeError(w, r, err)
		return
	}
	if err := a.svc.Report(r.Context(), userID, r.PathValue("id"), pos, body.Reason, body.Note); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
