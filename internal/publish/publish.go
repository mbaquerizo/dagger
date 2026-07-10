package publish

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mbaquerizo/dagger/internal/issues"
)

type poolIface interface {
	Begin(context.Context) (pgx.Tx, error)
}

func Publish(ctx context.Context, pool poolIface, req PublishRequest, workspaceID int, baseURL string, authProjectID *int) (*PublishResponse, error) {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var slug string

	err = tx.QueryRow(ctx, `
		SELECT slug FROM projects
		WHERE id = $1
		AND workspace_id = $2
	`, req.ProjectID, workspaceID).Scan(&slug)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("project %d not found", req.ProjectID)
		}

		return nil, fmt.Errorf("retrieving project: %w", err)
	}

	var displayNumber int

	err = tx.QueryRow(ctx, `
		UPDATE projects
		SET next_display_number = next_display_number + 1
		WHERE id = $1
		AND workspace_id = $2
		RETURNING next_display_number - 1
	`, req.ProjectID, workspaceID).Scan(&displayNumber)

	if err != nil {
		return nil, fmt.Errorf("allocating display id: %w", err)
	}

	displayID := fmt.Sprintf("%s-%d", slug, displayNumber)

	var entityType string

	switch req.Type {
	case "issue":
		entityType = "issues"
	default:
		entityType = "docs"
	}

	if req.ParentID != nil {
		var parentExists bool

		parentQuery := `
			SELECT EXISTS(
				SELECT 1 FROM %s
				WHERE id = $1
				AND workspace_id = $2
			)
		`

		if authProjectID != nil {
			err = tx.QueryRow(ctx, fmt.Sprintf(parentQuery, entityType)+" AND project_id = $3", req.ParentID, workspaceID, *authProjectID).Scan(&parentExists)
		} else {
			err = tx.QueryRow(ctx, fmt.Sprintf(parentQuery, entityType), req.ParentID, workspaceID).Scan(&parentExists)
		}

		if err != nil {
			return nil, fmt.Errorf("looking up parent: %w", err)
		}

		if !parentExists {
			return nil, fmt.Errorf("parent %s %d not found", entityType, *req.ParentID)
		}
	}

	var insertedID int

	if req.Type == "issue" {
		var issueTypeID int

		err = tx.QueryRow(ctx, `
		SELECT id FROM issue_types
		WHERE name = $1
		`, *req.Metadata.IssueType).Scan(&issueTypeID)

		if err != nil {
			return nil, fmt.Errorf("looking up issue type: %w", err)
		}

		status := "open"

		if req.Metadata.Status != nil {
			status = *req.Metadata.Status
		}

		err = tx.QueryRow(ctx, `
			INSERT INTO issues (display_id, issue_type_id, title, body, status, parent_id, project_id, workspace_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`, displayID, issueTypeID, req.Title, req.Body, status, req.ParentID, req.ProjectID, workspaceID).Scan(&insertedID)

		if err != nil {
			return nil, fmt.Errorf("inserting issue: %w", err)
		}
	} else {
		status := "proposed"

		if req.Metadata.Status != nil {
			status = *req.Metadata.Status
		}

		err = tx.QueryRow(ctx, `
			INSERT INTO docs (display_id, type, title, body, status, parent_id, project_id, workspace_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`, displayID, req.Type, req.Title, req.Body, status, req.ParentID, req.ProjectID, workspaceID).Scan(&insertedID)

		if err != nil {
			return nil, fmt.Errorf("inserting doc: %w", err)
		}
	}

	var targetRelationshipTable string

	if entityType == "issues" {
		targetRelationshipTable = "docs"
	} else {
		targetRelationshipTable = "issues"
	}

	for _, rel := range req.Metadata.Relationships {
		var targetExists bool

		relationshipQuery := `
			SELECT EXISTS(
				SELECT 1 FROM %s
				WHERE id = $1
				AND workspace_id = $2
		`

		if authProjectID != nil {
			err = tx.QueryRow(ctx, fmt.Sprintf(relationshipQuery, targetRelationshipTable)+" AND project_id = $3)", rel.TargetID, workspaceID, *authProjectID).Scan(&targetExists)
		} else {
			err = tx.QueryRow(ctx, fmt.Sprintf(relationshipQuery, targetRelationshipTable)+")", rel.TargetID, workspaceID).Scan(&targetExists)
		}

		if err != nil {
			return nil, fmt.Errorf("checking relationship target: %w", err)
		}

		if !targetExists {
			return nil, fmt.Errorf("relationship target %s %d not found", targetRelationshipTable, rel.TargetID)
		}

		if entityType == "issues" {
			_, err = tx.Exec(ctx, `
				INSERT INTO doc_issues (doc_id, issue_id, relationship_type)
				VALUES ($1, $2, $3)
			`, rel.TargetID, insertedID, rel.Type)
		} else {
			_, err = tx.Exec(ctx, `
				INSERT INTO doc_issues (doc_id, issue_id, relationship_type)
				VALUES ($1, $2, $3)
			`, insertedID, rel.TargetID, rel.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("inserting relationship: %w", err)
		}
	}

	if entityType == "issues" {
		var targetExists bool

		for _, rel := range req.Metadata.IssueRelations {
			issueRelationQuery := `
				SELECT EXISTS(
					SELECT 1 FROM issues
					WHERE id = $1
					AND workspace_id = $2
			`

			if authProjectID != nil {
				err = tx.QueryRow(ctx, issueRelationQuery+" AND project_id = $3)", rel.TargetID, workspaceID, *authProjectID).Scan(&targetExists)
			} else {
				err = tx.QueryRow(ctx, issueRelationQuery+")", rel.TargetID, workspaceID).Scan(&targetExists)
			}

			if err != nil {
				return nil, fmt.Errorf("checking issue relation target: %w", err)
			}

			if !targetExists {
				return nil, fmt.Errorf("issue relation target issue %d not found", rel.TargetID)
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO issue_relations (source_issue_id, target_issue_id, relation_id)
				VALUES
					($1, $2, (SELECT id FROM relations WHERE name = $3)),
					($2, $1, (SELECT id FROM relations WHERE name = $4))
			`, insertedID, rel.TargetID, rel.RelationType, issues.RelationInverse[rel.RelationType])

			if err != nil {
				return nil, fmt.Errorf("inserting issue relations: %w", err)
			}
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &PublishResponse{
		ID:        insertedID,
		DisplayID: displayID,
		URL:       fmt.Sprintf("%s/%s/%s", baseURL, entityType, displayID),
	}, nil
}
