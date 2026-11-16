package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/question"
)

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
			writeError(w, r, errors.Join(question.ErrInvalid, errors.New("limit must be a number")))
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
		writeError(w, r, question.ErrNotFound)
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
	v, err := a.exam.Answer(r.Context(), userID, r.PathValue("id"), exam.AnswerInput{Pos: pos, Answer: body.Answer, TimeMS: body.TimeMS, Flagged: body.Flagged, Seq: body.Seq})
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

func (a *api) reportItem(w http.ResponseWriter, r *http.Request, userID string) {
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
	if err := a.exam.Report(r.Context(), userID, r.PathValue("id"), pos, body.Reason, body.Note); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
