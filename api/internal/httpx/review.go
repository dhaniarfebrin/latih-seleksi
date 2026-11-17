package httpx

import "net/http"

// reviews answers how many facts and patterns are due now; no keys or snapshots.
func (a *api) reviews(w http.ResponseWriter, r *http.Request, userID string) {
	c, err := a.review.Counts(r.Context(), userID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}
