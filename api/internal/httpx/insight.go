package httpx

import "net/http"

func (a *api) dashboard(w http.ResponseWriter, r *http.Request, userID string) {
	d, err := a.insight.Dashboard(r.Context(), userID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, d)
}
