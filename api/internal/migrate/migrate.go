package migrate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/mbaquerizo/dagger/internal/db"
)

func Run(databaseURL, migrationsDir string) error {
	pool, err := db.Connect(databaseURL)

	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	defer pool.Close()

	ctx := context.Background()

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
    	version TEXT PRIMARY KEY,
    	applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	files, err := os.ReadDir(migrationsDir)

	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	var sqlFiles []string

	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}

	sort.Strings(sqlFiles)

	for _, file := range sqlFiles {
		var exists bool

		err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)", file).Scan(&exists)

		if err != nil {
			return fmt.Errorf("checking migration %s: %w", file, err)
		}

		if exists {
			fmt.Printf("SKIP    %s (already applied)\n", file)
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, file))

		if err != nil {
			return fmt.Errorf("reading migration %s: %w", file, err)
		}

		_, err = pool.Exec(ctx, string(content))

		if err != nil {
			return fmt.Errorf("applying migration %s: %w", file, err)
		}

		_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", file)

		if err != nil {
			return fmt.Errorf("recording migration %s: %w", file, err)
		}

		fmt.Printf("OK    %s\n", file)
	}

	return nil
}
