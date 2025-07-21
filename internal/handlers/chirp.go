package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PeterRonne/chirpy/internal/models"
	"github.com/PeterRonne/chirpy/pkg/response"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	var req models.ChirpRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		response.WithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Body) > 140 {
		response.WithError(w, http.StatusBadRequest, "chirp is too long")
		return
	}

	words := strings.Split(req.Body, " ")
	msgSlice := make([]string, 0, len(words))
	for _, word := range words {
		if val, ok := models.BadWords[strings.ToLower(word)]; ok {
			word = val
		}
		msgSlice = append(msgSlice, word)
	}
	cleanedMsg := strings.Join(msgSlice, " ")

	response.WithJSON(w, http.StatusOK, models.ValidateChirpResponse{CleanedBody: cleanedMsg})
}
