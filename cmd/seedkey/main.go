package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/mbaquerizo/dagger/internal/db"
)

func main() {
	var workspaceID int
	var name string

	flag.IntVar(&workspaceID, "workspace-id", 0, "ID of the workspace to associate the key with")
	flag.StringVar(&name, "name", "seed key", "Human-readable name for the API key")

	flag.Parse()

	if workspaceID == 0 {
		log.Fatal("--workspace-id is required")
	}

	godotenv.Load()

	databaseURL := os.Getenv("DB_URL")

	if databaseURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	pool, err := db.Connect(databaseURL)

	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}

	defer pool.Close()

	rawKey, hash, prefix, err := auth.GenerateKey()

	if err != nil {
		log.Fatalf("generating key: %v", err)
	}

	_, err = pool.Exec(context.Background(),
		`INSERT INTO api_keys (key_hash, prefix, name, workspace_id)
		VALUES ($1, $2, $3, $4)`,
		hash, prefix, name, workspaceID)

	if err != nil {
		log.Fatalf("inserting key: %v", err)
	}

	fmt.Println(rawKey)
}
