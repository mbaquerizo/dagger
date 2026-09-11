package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/pashagolub/pgxmock/v5"
)

func TestMCPHandler_ToolsList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewMCPHandler(mock, "http://localhost:8080")

	body := `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	req = req.WithContext(auth.WithWorkspaceID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
}

func TestMCPHandler_ToolsCallGetIssue(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-42", "story", "Test issue", nil, "open", nil, 1, 1))
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}))
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}))

	handler := NewMCPHandler(mock, "http://localhost:8080")

	body := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_issue","arguments":{"display_id":"DGR-42"}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	req = req.WithContext(auth.WithWorkspaceID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
}

func TestMCPHandler_Initialize(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewMCPHandler(mock, "http://localhost:8080")

	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	req = req.WithContext(auth.WithWorkspaceID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	if resp.Result == nil {
		t.Fatal("result should not be nil")
	}
}

func TestMCPHandler_NotificationsInitialized(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewMCPHandler(mock, "http://localhost:8080")

	body := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	req = req.WithContext(auth.WithWorkspaceID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body for notification, got %q", rec.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestMCPHandler_BadJSON(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewMCPHandler(mock, "http://localhost:8080")

	body := `not json`
	req := httptest.NewRequest("POST", "/mcp", strings.NewReader(body))
	req = req.WithContext(auth.WithWorkspaceID(req.Context(), 1))
	rec := httptest.NewRecorder()

	handler(rec, req)

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error == nil {
		t.Fatal("expected parse error, got nil")
	}
	if resp.Error.Code != ErrCodeParse {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeParse)
	}
}
