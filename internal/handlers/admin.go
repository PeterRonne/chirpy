package handlers

import (
	"fmt"
	"net/http"

	"github.com/PeterRonne/chirpy/internal/config"
)

func GetMetrics(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)

		html := fmt.Sprintf(`
			<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>
		`, cfg.FileserverHits.Load())
		w.Write([]byte(html))
	}
}

func Reset(cfg *config.APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Store(0)
		w.WriteHeader(200)
		w.Write([]byte(""))
	}
}
