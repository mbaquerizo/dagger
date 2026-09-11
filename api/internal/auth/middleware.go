package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type keyQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewMiddleware(q keyQuerier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			rawKey := strings.TrimPrefix(authHeader, "Bearer ")

			if !strings.HasPrefix(rawKey, KeyPrefix) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			incomingHash := HashKey(rawKey)

			ctx := r.Context()

			var (
				keyID       int
				workspaceID int
				projectID   *int
			)

			err := q.QueryRow(ctx,
				`SELECT id, workspace_id, project_id
				 FROM api_keys
				 WHERE key_hash = $1
				 	 AND (expires_at IS NULL OR expires_at > NOW())`,
				incomingHash).Scan(&keyID, &workspaceID, &projectID)

			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ctx = WithKeyID(ctx, keyID)
			ctx = WithWorkspaceID(ctx, workspaceID)

			if projectID != nil {
				ctx = WithProjectID(ctx, *projectID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
