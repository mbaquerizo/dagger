package issues

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandler_GetIssue_ScopedTokenSuccess(t *testing.T) {
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
		WithArgs(1, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}))

	// No parent (parent_id is nil), skip

	// Children — empty
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(1, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	// Related issues — empty
	mockPool.ExpectQuery(`SELECT i.display_id, i.title, r.name`).
		WithArgs(1, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}))

	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/issues/DGR-42", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))
	r = r.WithContext(auth.WithProjectID(r.Context(), 1))

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

func TestHandler_GetIssue_ScopedTokenProjectIDMismatch(t *testing.T) {
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

	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/issues/DGR-42", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))
	r = r.WithContext(auth.WithProjectID(r.Context(), 2))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-42")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetIssueHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}

}

func TestHandler_ListIssues_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("open", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "title", "status", "type_name", "parent_display_id"}).
			AddRow(1, "DGR-1", "First issue", "open", "story", nil))

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

func TestHandler_ListIssues_DefaultsToAll(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "title", "status", "type_name", "parent_display_id"}))

	// No ?status= query param — handler should default to "all"
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

func TestHandler_UpdateIssueStatus_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	req := UpdateStatusRequest{
		Status: "in-review",
	}

	body, err := json.Marshal(req)

	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	mockPool.ExpectExec(`UPDATE issues`).
		WithArgs(req.Status, "DGR-42", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	r := httptest.NewRequest(http.MethodPatch, "/api/v1/issues/DGR-42/status", bytes.NewReader(body))

	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))
	r.Header.Set("Content-Type", "application/json")

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-42")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewUpdateIssueStatusHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandler_UpdateIssueStatus_ErrInvalidStatus(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	req := UpdateStatusRequest{
		Status: "pegasus",
	}

	body, err := json.Marshal(req)

	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	r := httptest.NewRequest(http.MethodPatch, "/api/v1/issues/DGR-42/status", bytes.NewReader(body))

	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))
	r.Header.Set("Content-Type", "application/json")

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-42")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewUpdateIssueStatusHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}
