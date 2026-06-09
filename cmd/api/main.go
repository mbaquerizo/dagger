package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"log"
)

func main() {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("server starting on :8080")

	http.ListenAndServe(":8080", r)
}