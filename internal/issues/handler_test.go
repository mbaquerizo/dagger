package issues

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/pashagolub/pgxmock/v5"
)

func TestHandler_GetIssue_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	issueBody := "Issue description"

	// Issue query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-42", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-42", "story", "Test Issue", &issueBody, "open", nil, 1, 1))

	// Linked docs — empty
	mockPool.ExpectQuery(`SELECT d.id, d.display_id, d.type`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}))

	// No parent (parent_id is nil), skip

	// Children — empty
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	// Related issues — empty
	mockPool.ExpectQuery(`SELECT i.display_id, i.title, r.name`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}))

	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/issues/DGR-42", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-42")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetIssueHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "text/markdown" {
		t.Errorf("expected Content-Type text/markdown, got %s", ct)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("expected non-empty body")
	}
	if !strings.Contains(body, "# DGR-42: Test Issue") {
		t.Errorf("expected header in body, got:\n%s", body)
	}
}

func TestHandler_GetIssue_NotFound(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-999", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/issues/DGR-999", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-999")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetIssueHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_GetIssue_NoWorkspace(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/issues/DGR-42", nil)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-42")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	// handler should fail fast without workspace context
	handler := NewGetIssueHandler(nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestListIssues_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	parentDisplayID := "EPIC-1"

	mockPool.ExpectQuery(`SELECT i.display_id, it.name`).
		WithArgs("open", 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "status", "type_name", "parent_display_id"}).
			AddRow("DGR-1", "First issue", "open", "story", nil).
			AddRow("DGR-2", "Second issue", "open", "bug", &parentDisplayID))

	issues, err := ListIssues(context.Background(), mockPool, "open", 1)
	if err != nil {
		t.Fatalf("ListIssues returned error: %v", err)
	}
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}

	if issues[0].DisplayID != "DGR-1" {
		t.Errorf("expected DGR-1, got %s", issues[0].DisplayID)
	}
	if issues[0].ParentDisplayID != nil {
		t.Errorf("expected nil parent for first issue, got %v", issues[0].ParentDisplayID)
	}

	if issues[1].DisplayID != "DGR-2" {
		t.Errorf("expected DGR-2, got %s", issues[1].DisplayID)
	}
	if issues[1].ParentDisplayID == nil || *issues[1].ParentDisplayID != "EPIC-1" {
		t.Errorf("expected parent EPIC-1, got %v", issues[1].ParentDisplayID)
	}
}

func TestListIssues_NoResults(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.display_id, it.name`).
		WithArgs("closed", 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "status", "type_name", "parent_display_id"}))

	issues, err := ListIssues(context.Background(), mockPool, "closed", 1)
	if err != nil {
		t.Fatalf("ListIssues returned error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestHandler_ListIssues_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.display_id, it.name`).
		WithArgs("open", 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "status", "type_name", "parent_display_id"}).
			AddRow("DGR-1", "First issue", "open", "story", nil))

	r := httptest.NewRequest(http.MethodGet, "/api/v1/issues?status=open", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	w := httptest.NewRecorder()

	handler := NewListIssuesHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	body := w.Body.String()
	if !strings.Contains(body, "DGR-1") {
		t.Errorf("expected DGR-1 in response, got %s", body)
	}
}

func TestHandler_ListIssues_DefaultsToOpen(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.display_id, it.name`).
		WithArgs("open", 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "status", "type_name", "parent_display_id"}))

	// No ?status= query param — handler should default to "open"
	r := httptest.NewRequest(http.MethodGet, "/api/v1/issues", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	w := httptest.NewRecorder()

	handler := NewListIssuesHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandler_ListIssues_NoWorkspace(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/issues", nil)

	w := httptest.NewRecorder()

	handler := NewListIssuesHandler(nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}


