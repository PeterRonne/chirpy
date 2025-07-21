package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func WithError(w http.ResponseWriter, code int, msg string) {
	WithJSON(w, code, ErrorResponse{ErrorMessage: msg})
}

func WithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
