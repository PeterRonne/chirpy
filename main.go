package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

var badwords = map[string]string{
	"kerfuffle": "****",
	"sharbert":  "****",
	"fornax":    "****",
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

type errorResponse struct {
	ErrorMessage string `json:"error"`
}

type validateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type chirpRequest struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	// message := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	html := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(html))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte(""))
}

func main() {
	mux := http.NewServeMux()
	apicfg := apiConfig{}

	mux.Handle("/app/", apicfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	))

	mux.HandleFunc("GET /admin/metrics", apicfg.getMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.reset)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		var req chirpRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if len(req.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "chirp is too long")
			return
		}

		words := strings.Split(req.Body, " ")
		msgSlice := make([]string, 0, len(words))
		for _, word := range words {
			if val, ok := badwords[strings.ToLower(word)]; ok {
				word = val
			}
			msgSlice = append(msgSlice, word)
		}
		cleanedMsg := strings.Join(msgSlice, " ")

		responWithJson(w, http.StatusOK, validateChirpResponse{CleanedBody: cleanedMsg})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server kører på http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

// Helper functions
func respondWithError(w http.ResponseWriter, code int, msg string) {
	responWithJson(w, code, errorResponse{ErrorMessage: msg})
}

func responWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
