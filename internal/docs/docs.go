package docs

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type poolIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func GetDoc(ctx context.Context, pool poolIface, displayID string, workspaceID int) (*Doc, error) {
	var doc Doc
	var parentDisplayID, parentTitle *string

	err := pool.QueryRow(ctx, `
		SELECT d.id, d.display_id, d.type, d.title, d.body, d.status, p.display_id, p.title
		FROM docs d
		LEFT JOIN docs p ON p.id = d.parent_id
		WHERE d.display_id = $1 AND d.workspace_id = $2
	`, displayID, workspaceID).Scan(
		&doc.ID,
		&doc.DisplayID,
		&doc.DocType,
		&doc.Title,
		&doc.Body,
		&doc.Status,
		&parentDisplayID,
		&parentTitle,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDocNotFound
		}

		return nil, fmt.Errorf("querying doc: %w", err)
	}

	if parentDisplayID != nil {
		doc.Parent = &ParentDoc{
			DisplayID: *parentDisplayID,
			Title:     *parentTitle,
		}
	}

	return &doc, nil
}
