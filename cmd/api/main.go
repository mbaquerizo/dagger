package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/mbaquerizo/dagger/internal/db"
	"github.com/mbaquerizo/dagger/internal/issues"
	"github.com/mbaquerizo/dagger/internal/publish"
)

func main() {
	godotenv.Load()

	databaseURL := os.Getenv("DB_URL")

	if databaseURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	baseURL := os.Getenv("BASE_URL")

	if baseURL == "" {
		log.Fatalf("BASE_URL environment variable is required")
	}

	pool, err := db.Connect(databaseURL)

	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}

	defer pool.Close()

	r := chi.NewRouter()

	r.Use(auth.NewMiddleware(pool))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/api/v1/publish", publish.NewHandler(pool, baseURL))

	r.Get("/api/v1/agent/issues/{displayId}", issues.NewGetIssueHandler(pool))

	r.Get("/api/v1/issues", issues.NewListIssuesHandler(pool))

	log.Println("server starting on :8080")

	http.ListenAndServe(":8080", r)
}
