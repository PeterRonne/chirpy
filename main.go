package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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

	// mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {

	// })

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server kører på http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
