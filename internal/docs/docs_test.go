package docs

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
)

func TestGetDoc_ByDisplayID_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	body := "ADR body content"

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("DGR-3", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(5, "DGR-3", "adr", "Test ADR", &body, "approved", 1, 1, nil, nil, nil))

	doc, err := GetDoc(context.Background(), mockPool, "DGR-3", 1, nil)
	if err != nil {
		t.Fatalf("GetDoc returned error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil Doc")
	}

	if doc.ID != 5 {
		t.Errorf("expected ID 5, got %d", doc.ID)
	}
	if doc.DisplayID != "DGR-3" {
		t.Errorf("expected DisplayID DGR-3, got %s", doc.DisplayID)
	}
	if doc.DocType != "adr" {
		t.Errorf("expected DocType adr, got %s", doc.DocType)
	}
	if doc.Title != "Test ADR" {
		t.Errorf("expected Title 'Test ADR', got %s", doc.Title)
	}
	if doc.Body == nil || *doc.Body != "ADR body content" {
		t.Errorf("expected Body 'ADR body content', got %v", doc.Body)
	}
	if doc.Status != "approved" {
		t.Errorf("expected Status approved, got %s", doc.Status)
	}
	if doc.Parent != nil {
		t.Errorf("expected nil parent, got %v", doc.Parent)
	}
}

func TestGetDoc_NotFound(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("DGR-999", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_project_id", "p_display_id", "p_title"}))

	doc, err := GetDoc(context.Background(), mockPool, "DGR-999", 1, nil)
	if err == nil {
		t.Fatal("expected ErrDocNotFound, got nil")
	}
	if !errors.Is(err, ErrDocNotFound) {
		t.Errorf("expected ErrDocNotFound, got %v", err)
	}
	if doc != nil {
		t.Errorf("expected nil doc, got %v", doc)
	}
}

func TestGetDoc_WithNilBody(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("CE-1", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(10, "CE-1", "adr", "Code Exploration", nil, "proposed", 1, 1, nil, nil, nil))

	doc, err := GetDoc(context.Background(), mockPool, "CE-1", 1, nil)
	if err != nil {
		t.Fatalf("GetDoc returned error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil Doc")
	}
	if doc.Body != nil {
		t.Errorf("expected nil body, got %v", doc.Body)
	}
	if doc.Status != "proposed" {
		t.Errorf("expected Status proposed, got %s", doc.Status)
	}
	if doc.Parent != nil {
		t.Errorf("expected nil parent, got %v", doc.Parent)
	}
}

func TestGetDoc_WithProjectIDMismatch(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	projectID := 2

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("CE-1", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(10, "CE-1", "adr", "Code Exploration", nil, "proposed", 1, 1, nil, nil, nil))

	doc, err := GetDoc(context.Background(), mockPool, "CE-1", 1, &projectID)
	if err == nil {
		t.Fatal("expected ErrProjectIDMismatch, but got nil")
	}
	if !errors.Is(err, ErrProjectIDMismatch) {
		t.Errorf("expected ErrProjectIDMismatch, buy got %v", err)
	}
	if doc != nil {
		t.Errorf("expected nil doc, but got %v", doc)
	}
}

func TestGetDoc_WithParentProjectIDMismatch(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	projectID := 2

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("CE-1", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(10, "CE-1", "adr", "Code Exploration", nil, "proposed", 1, 2, nil, nil, nil))

	doc, err := GetDoc(context.Background(), mockPool, "CE-1", 1, &projectID)

	if err != nil {
		t.Fatalf("GetDoc returned error: %v", err)
	}

	if doc == nil {
		t.Fatal("expected non-nil Doc")
	}

	if doc.ProjectID != 2 {
		t.Errorf("expected project id 2, but got %d", doc.ProjectID)
	}

	if doc.Parent != nil {
		t.Errorf("expected nil parent, but got %v", doc.Parent)
	}
}
