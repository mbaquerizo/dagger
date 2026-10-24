package docs

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type poolIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
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

func UpdateDoc(ctx context.Context, pool poolIface, req UpdateDocRequest, displayID string, workspaceID int, authProjectID *int) error {
	var docProjectID int

	checkQuery := `
		SELECT project_id FROM docs
		WHERE display_id = $1 AND workspace_id = $2
	`

	err := pool.QueryRow(ctx, checkQuery, displayID, workspaceID).Scan(&docProjectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDocNotFound
		}
		return fmt.Errorf("checking doc project: %w", err)
	}

	if authProjectID != nil && docProjectID != *authProjectID {
		return ErrProjectIDMismatch
	}

	var result pgconn.CommandTag

	baseQuery := `
		UPDATE docs
		SET title = COALESCE($1, title),
		    body = COALESCE($2, body),
		    status = COALESCE($3, status)
		WHERE display_id = $4
		AND workspace_id = $5
	`

	var titlePtr, bodyPtr, statusPtr any
	if req.Title != nil {
		titlePtr = *req.Title
	}
	if req.Body != nil {
		bodyPtr = *req.Body
	}
	if req.Status != nil {
		statusPtr = *req.Status
	}

	if authProjectID != nil {
		result, err = pool.Exec(ctx, baseQuery+" AND project_id = $6", titlePtr, bodyPtr, statusPtr, displayID, workspaceID, *authProjectID)
	} else {
		result, err = pool.Exec(ctx, baseQuery, titlePtr, bodyPtr, statusPtr, displayID, workspaceID)
	}

	if err != nil {
		return fmt.Errorf("updating doc: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrDocNotFound
	}

	return nil
}
