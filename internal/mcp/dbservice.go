package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/mbaquerizo/dagger/internal/docs"
	"github.com/mbaquerizo/dagger/internal/issues"
)

type poolIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type DBService struct {
	pool poolIface
}

func NewDBService(pool poolIface) *DBService {
	return &DBService{pool: pool}
}

func (s *DBService) GetIssue(ctx context.Context, displayID string) (ToolResult, error) {
	workspaceID, ok := auth.WorkspaceIDFromContext(ctx)

	if !ok {
		return ToolResult{}, fmt.Errorf("unauthorized")
	}

	projectID, hasProjectID := auth.ProjectIDFromContext(ctx)
	var projectIDPtr *int

	if hasProjectID {
		projectIDPtr = &projectID
	}

	issueContext, err := issues.GetIssueContext(ctx, s.pool, displayID, workspaceID, projectIDPtr)

	if err != nil {
		return ToolResult{}, err
	}

	markdown := issues.RenderIssueContext(issueContext)

	return ToolResult{
		Content: []ContentItem{{Type: "text", Text: markdown}},
	}, nil
}

func (s *DBService) GetDoc(ctx context.Context, displayID string) (ToolResult, error) {
	workspaceID, ok := auth.WorkspaceIDFromContext(ctx)

	if !ok {
		return ToolResult{}, fmt.Errorf("unauthorized")
	}

	projectID, hasProjectID := auth.ProjectIDFromContext(ctx)
	var projectIDPtr *int

	if hasProjectID {
		projectIDPtr = &projectID
	}

	doc, err := docs.GetDoc(ctx, s.pool, displayID, workspaceID, projectIDPtr)

	if err != nil {
		return ToolResult{}, err
	}

	markdown := docs.RenderDoc(doc)

	return ToolResult{
		Content: []ContentItem{{Type: "text", Text: markdown}},
	}, nil
}

func (s *DBService) ListIssues(ctx context.Context, status string) (ToolResult, error) {
	workspaceID, ok := auth.WorkspaceIDFromContext(ctx)

	if !ok {
		return ToolResult{}, fmt.Errorf("unauthorized")
	}

	projectID, hasProjectID := auth.ProjectIDFromContext(ctx)
	var projectIDPtr *int

	if hasProjectID {
		projectIDPtr = &projectID
	}

	if status == "" {
		status = "open"
	}

	issues, err := issues.ListIssues(ctx, s.pool, status, workspaceID, projectIDPtr)

	if err != nil {
		return ToolResult{}, err
	}

	body, err := json.Marshal(issues)

	if err != nil {
		return ToolResult{}, fmt.Errorf("marshaling issues: %w", err)
	}

	return ToolResult{
		Content: []ContentItem{{Type: "text", Text: string(body)}},
	}, nil
}

func (s *DBService) UpdateIssueStatus(ctx context.Context, displayID string, newStatus string) (ToolResult, error) {
	workspaceID, ok := auth.WorkspaceIDFromContext(ctx)

	if !ok {
		return ToolResult{}, fmt.Errorf("unauthorized")
	}

	projectID, hasProjectID := auth.ProjectIDFromContext(ctx)
	var projectIDPtr *int

	if hasProjectID {
		projectIDPtr = &projectID
	}

	if err := issues.ValidateStatus(newStatus); err != nil {
		return ToolResult{}, fmt.Errorf("invalid status: %w", err)
	}

	req := issues.UpdateStatusRequest{Status: newStatus}

	err := issues.UpdateIssueStatus(ctx, s.pool, req, displayID, workspaceID, projectIDPtr)

	if err != nil {
		return ToolResult{}, err
	}

	return ToolResult{
		Content: []ContentItem{{Type: "text", Text: "Updated " + displayID + " to " + newStatus}},
	}, nil
}
