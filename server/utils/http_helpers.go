package utils

import (
	"encoding/json"
	"net/http"
)

// string is for the initial label
// interface so we can pass whatever our hearts desire
type JSONResponse map[string]any

// pass a writer, a status, and a any data, to encode a new json response.
func WriteJSON(w http.ResponseWriter, status int, data JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusBadRequest)
	}
}
