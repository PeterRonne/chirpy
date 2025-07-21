package config

import (
	"net/http"
	"sync/atomic"

	"github.com/PeterRonne/chirpy/internal/database"
)

type APIConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
}

func NewAPIConfig(dbQueries *database.Queries) *APIConfig {
	return &APIConfig{
		DBQueries: dbQueries,
	}
}

func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
