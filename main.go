package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PeterRonne/chirpy/internal/config"
	"github.com/PeterRonne/chirpy/internal/database"
	"github.com/PeterRonne/chirpy/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	cfg := config.NewAPIConfig(dbQueries)

	mux := setupRoutes(cfg)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server kører på http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func setupRoutes(cfg *config.APIConfig) *http.ServeMux {
	mux := http.NewServeMux()

	// File server with metrics middleware
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	))

	// Admin routes
	mux.HandleFunc("GET /admin/metrics", handlers.GetMetrics(cfg))
	mux.HandleFunc("POST /admin/reset", handlers.Reset(cfg))

	// API routes
	mux.HandleFunc("GET /api/healthz", handlers.Healthz)
	mux.HandleFunc("POST /api/validate_chirp", handlers.ValidateChirp)

	return mux
}
