package docs

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

func TestHandler_GetDoc_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	body := "ADR content"

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("DGR-3", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(5, "DGR-3", "adr", "Test ADR", &body, "approved", 1, 1, nil, nil, nil, nil))

	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/docs/DGR-3", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-3")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetDocHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "text/markdown" {
		t.Errorf("expected Content-Type text/markdown, got %s", ct)
	}

	bodyStr := w.Body.String()
	if bodyStr == "" {
		t.Error("expected non-empty body")
	}
	if !strings.Contains(bodyStr, "# DGR-3: Test ADR") {
		t.Errorf("expected header in response, got:\n%s", bodyStr)
	}
	if !strings.Contains(bodyStr, "ADR content") {
		t.Errorf("expected body content in response, got:\n%s", bodyStr)
	}
}

func TestHandler_GetDoc_NotFound(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT d\.id, d\.display_id, d\.type, d\.title, d\.body, d\.status`).
		WithArgs("DGR-999", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_id", "p_project_id", "p_display_id", "p_title"}))

	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/docs/DGR-999", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-999")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetDocHandler(mockPool)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_GetDoc_NoWorkspace(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/docs/DGR-3", nil)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "DGR-3")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetDocHandler(nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandler_GetDoc_MissingDisplayID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/agent/docs/", nil)
	r = r.WithContext(auth.WithWorkspaceID(r.Context(), 1))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("displayId", "")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	handler := NewGetDocHandler(nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
