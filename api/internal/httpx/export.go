package httpx

import "net/http"

func (a *api) exportMe(w http.ResponseWriter, r *http.Request, userID string) {
	d, err := a.export.Build(r.Context(), userID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="latih-seleksi-data.json"`)
	writeJSON(w, http.StatusOK, d)
}
