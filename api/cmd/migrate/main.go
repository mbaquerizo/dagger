package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mbaquerizo/dagger/internal/migrate"
)

func main() {
	godotenv.Load()

	databaseURL := os.Getenv("DB_URL")

	if databaseURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	if err := migrate.Run(databaseURL, "db/migrations"); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("migrations complete")
}
