package middleware

import (
	"net/http"

	"github.com/PeterRonne/chirpy/internal/config"
)

func MiddlewareMetricsInc(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
