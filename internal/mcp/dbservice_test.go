package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/mbaquerizo/dagger/internal/publish"
	"github.com/pashagolub/pgxmock/v5"
)

func TestDBService_GetIssue(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	// Main issue query (QueryRow: displayID, workspaceID)
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-42", "story", "Test issue", nil, "open", nil, 1, 1))

	// Linked docs (Query: issueID, workspaceID)
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}))

	// Children (Query: parentID, workspaceID)
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	// Related issues (Query: issueID, workspaceID)
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}))

	svc := NewDBService(mock, "http://localhost:8080")
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	result, err := svc.GetIssue(ctx, "DGR-42")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "DGR-42") {
		t.Errorf("content missing issue ID: %s", result.Content[0].Text)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDBService_GetDoc(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(1, "DOC-1", "adr", "Test doc", nil, "approved", 1, 1, nil, nil, nil))

	svc := NewDBService(mock, "http://localhost:8080")
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	result, err := svc.GetDoc(ctx, "DOC-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "DOC-1") {
		t.Errorf("content missing doc ID: %s", result.Content[0].Text)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDBService_ListIssues(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "status", "type_name", "parent_display_id"}).
			AddRow("DGR-42", "Test issue", "open", "story", nil))

	svc := NewDBService(mock, "http://localhost:8080")
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	result, err := svc.ListIssues(ctx, "open")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "DGR-42") {
		t.Errorf("content missing issue ID: %s", result.Content[0].Text)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDBService_Publish(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(".*").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mock.ExpectQuery(".*").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{""}).AddRow(42))
	mock.ExpectQuery(".*").WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(99))
	mock.ExpectCommit()

	svc := NewDBService(mock, "https://api.dagger.sh")
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	req := publish.PublishRequest{
		Type:      "adr",
		Title:     "Test ADR",
		Body:      "# Test",
		ProjectID: 1,
	}

	result, err := svc.Publish(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "DGR-42") {
		t.Errorf("content missing display ID: %s", result.Content[0].Text)
	}
	if !strings.Contains(result.Content[0].Text, "99") {
		t.Errorf("content missing internal ID: %s", result.Content[0].Text)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDBService_UpdateIssueStatus(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	svc := NewDBService(mock, "http://localhost:8080")
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	result, err := svc.UpdateIssueStatus(ctx, "DGR-42", "in-review")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "DGR-42") {
		t.Errorf("content missing issue ID: %s", result.Content[0].Text)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
