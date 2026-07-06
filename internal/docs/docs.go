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

func GetDoc(ctx context.Context, pool poolIface, displayID string, workspaceID int, authProjectID *int) (*Doc, error) {
	var doc Doc
	var parentDisplayID, parentTitle *string
	var parentID, parentProjectID *int

	err := pool.QueryRow(ctx, `
		SELECT d.id, d.display_id, d.type, d.title, d.body, d.status, d.workspace_id, d.project_id, p.id, p.project_id, p.display_id, p.title
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
		&doc.WorkspaceID,
		&doc.ProjectID,
		&parentID,
		&parentProjectID,
		&parentDisplayID,
		&parentTitle,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDocNotFound
		}

		return nil, fmt.Errorf("querying doc: %w", err)
	}

	if authProjectID != nil && *authProjectID != doc.ProjectID {
		return nil, ErrProjectIDMismatch
	}

	if authProjectID != nil && parentProjectID != nil && *authProjectID == *parentProjectID && parentDisplayID != nil {
		doc.Parent = &ParentDoc{
			ID:        *parentID,
			DisplayID: *parentDisplayID,
			Title:     *parentTitle,
		}
	}

	return &doc, nil
}
