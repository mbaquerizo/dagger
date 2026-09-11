package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mbaquerizo/dagger/internal/auth"
	"github.com/pashagolub/pgxmock/v5"
)

func TestServer_ToolsList(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)

	resp := server.HandleRequest(context.Background(), Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("result should be a map")
	}

	tools, ok := result["tools"].([]ToolDefinition)
	if !ok {
		t.Fatal("result.tools should be []ToolDefinition")
	}

	if len(tools) != 6 {
		t.Fatalf("got %d tools, want 6", len(tools))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_ToolsCallGetIssue(t *testing.T) {
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

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	resp := callTool(ctx, server, "get_issue", map[string]interface{}{
		"display_id": "DGR-42",
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}

	result, ok := resp.Result.(ToolResult)
	if !ok {
		t.Fatal("result should be a ToolResult")
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}
	if result.Content[0].Type != "text" {
		t.Errorf("content type = %q, want text", result.Content[0].Type)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_ToolsCallGetDoc(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status", "workspace_id", "project_id", "p_id", "p_project_id", "p_display_id", "p_title"}).
			AddRow(1, "DOC-1", "adr", "Test doc", nil, "approved", 1, 1, nil, nil, nil, nil))

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	resp := callTool(ctx, server, "get_doc", map[string]interface{}{
		"display_id": "DOC-1",
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	result, ok := resp.Result.(ToolResult)
	if !ok {
		t.Fatal("result should be a ToolResult")
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_ToolsCallListIssues(t *testing.T) {
	t.Run("with status", func(t *testing.T) {
		mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()

		mock.ExpectQuery(".*").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "title", "status", "type_name", "parent_display_id"}).
				AddRow(1, "DGR-42", "Test issue", "open", "story", nil))

		svc := NewDBService(mock, "http://localhost:8080")
		server := NewServer(svc)
		ctx := auth.WithWorkspaceID(context.Background(), 1)

		resp := callTool(ctx, server, "list_issues", map[string]interface{}{
			"status": "open",
		})
		if resp.Error != nil {
			t.Fatalf("unexpected error: %+v", resp.Error)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("without status", func(t *testing.T) {
		mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()

		mock.ExpectQuery(".*").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "title", "status", "type_name", "parent_display_id"}).
				AddRow(1, "DGR-43", "Another issue", "open", "task", nil))

		svc := NewDBService(mock, "http://localhost:8080")
		server := NewServer(svc)
		ctx := auth.WithWorkspaceID(context.Background(), 1)

		resp := callTool(ctx, server, "list_issues", map[string]interface{}{})
		if resp.Error != nil {
			t.Fatalf("unexpected error: %+v", resp.Error)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestServer_ToolsCallUpdateStatus(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	resp := callTool(ctx, server, "update_issue_status", map[string]interface{}{
		"display_id": "DGR-42",
		"status":     "in-review",
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	result, ok := resp.Result.(ToolResult)
	if !ok {
		t.Fatal("result should be a ToolResult")
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_ToolsCallAddIssueRelation(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherAny))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectExec(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 2))

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	resp := callTool(ctx, server, "add_issue_relation", map[string]interface{}{
		"source_id":     float64(5),
		"target_id":     float64(6),
		"relation_type": "blocks",
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}

	result, ok := resp.Result.(ToolResult)
	if !ok {
		t.Fatal("result should be a ToolResult")
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_ToolsCallPublish(t *testing.T) {
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

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)
	ctx := auth.WithWorkspaceID(context.Background(), 1)

	resp := callTool(ctx, server, "publish", map[string]interface{}{
		"type":       "adr",
		"title":      "Test ADR",
		"body":       "# Test",
		"project_id": float64(1),
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	result, ok := resp.Result.(ToolResult)
	if !ok {
		t.Fatal("result should be a ToolResult")
	}
	if len(result.Content) != 1 {
		t.Fatalf("got %d content items, want 1", len(result.Content))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_Initialize(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)

	tests := []struct {
		name        string
		params      map[string]interface{}
		wantVersion string
	}{
		{"2024-11-05", map[string]interface{}{"protocolVersion": "2024-11-05"}, "2024-11-05"},
		{"2025-03-26", map[string]interface{}{"protocolVersion": "2025-03-26"}, "2025-03-26"},
		{"2025-06-18", map[string]interface{}{"protocolVersion": "2025-06-18"}, "2025-06-18"},
		{"2025-11-25", map[string]interface{}{"protocolVersion": "2025-11-25"}, "2025-11-25"},
		{"unknown version defaults to latest", map[string]interface{}{"protocolVersion": "1.0.0"}, "2025-11-25"},
		{"empty version defaults to latest", map[string]interface{}{"protocolVersion": ""}, "2025-11-25"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := server.HandleRequest(context.Background(), Request{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "initialize",
				Params:  tt.params,
			})

			if resp.Error != nil {
				t.Fatalf("unexpected error: %+v", resp.Error)
			}

			result, ok := resp.Result.(InitializeResult)
			if !ok {
				t.Fatal("result should be InitializeResult")
			}
			if result.ProtocolVersion != tt.wantVersion {
				t.Errorf("ProtocolVersion = %q, want %q", result.ProtocolVersion, tt.wantVersion)
			}
			if result.ServerInfo.Name == "" {
				t.Error("ServerInfo.Name should not be empty")
			}
			if result.ServerInfo.Version == "" {
				t.Error("ServerInfo.Version should not be empty")
			}
		})
	}

	t.Run("missing params", func(t *testing.T) {
		resp := server.HandleRequest(context.Background(), Request{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
		})

		if resp.Error == nil {
			t.Fatal("expected error for missing params, got nil")
		}
		if resp.Error.Code != ErrCodeInvalidParams {
			t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeInvalidParams)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_NotificationsInitialized(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)

	resp := server.HandleRequest(context.Background(), Request{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	})

	if resp.JSONRPC != "" {
		t.Error("expected empty JSONRPC for notification")
	}
	if resp.Result != nil {
		t.Error("expected nil result for notification")
	}
	if resp.Error != nil {
		t.Errorf("unexpected error: %+v", resp.Error)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_InitializeResult_Marshal(t *testing.T) {
	result := InitializeResult{
		ProtocolVersion: "2025-03-26",
		ServerCapabilities: ServerCapabilities{
			Tools: struct{}{},
		},
		ServerInfo: ServerInfo{
			Name:    "dagger-mcp",
			Version: "1.0.0",
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["protocolVersion"] != "2025-03-26" {
		t.Errorf("protocolVersion = %v, want 2025-03-26", m["protocolVersion"])
	}

	caps, ok := m["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatal("capabilities should be an object")
	}

	tools, ok := caps["tools"]
	if !ok {
		t.Fatal("capabilities.tools should be present")
	}

	toolsObj, ok := tools.(map[string]interface{})
	if !ok {
		t.Fatalf("capabilities.tools should be an object, got %T", tools)
	}
	if len(toolsObj) != 0 {
		t.Errorf("capabilities.tools should be empty, got %v", toolsObj)
	}

	info, ok := m["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("serverInfo should be an object")
	}
	if info["name"] != "dagger-mcp" {
		t.Errorf("serverInfo.name = %v, want dagger-mcp", info["name"])
	}
}

func TestServer_MethodNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)

	resp := server.HandleRequest(context.Background(), Request{
		JSONRPC: "2.0",
		ID:      42,
		Method:  "unknown_method",
	})

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeMethodNotFound)
	}
	if resp.ID != 42 {
		t.Errorf("response id = %d, want 42", resp.ID)
	}
	if resp.Result != nil {
		t.Error("result should be nil for error responses")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_MissingRequiredParams(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)

	tests := []struct {
		name     string
		toolName string
		args     map[string]interface{}
	}{
		{"get_issue missing display_id", "get_issue", map[string]interface{}{}},
		{"get_doc missing display_id", "get_doc", map[string]interface{}{}},
		{"update_issue_status missing all", "update_issue_status", map[string]interface{}{}},
		{"update_issue_status missing status", "update_issue_status", map[string]interface{}{"display_id": "DGR-42"}},
		{"publish missing type", "publish", map[string]interface{}{"title": "t", "body": "b", "project_id": float64(1)}},
		{"publish missing project_id", "publish", map[string]interface{}{"type": "adr", "title": "t", "body": "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := callTool(context.Background(), server, tt.toolName, tt.args)
			if resp.Error == nil {
				t.Fatal("expected error for missing params, got nil")
			}
			if resp.Error.Code != ErrCodeInvalidParams {
				t.Errorf("error code = %d, want %d", resp.Error.Code, ErrCodeInvalidParams)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestServer_IDPropagation(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	svc := NewDBService(mock, "http://localhost:8080")
	server := NewServer(svc)

	resp := server.HandleRequest(context.Background(), Request{
		JSONRPC: "2.0",
		ID:      99,
		Method:  "tools/list",
	})

	if resp.ID != 99 {
		t.Errorf("response ID = %d, want 99", resp.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func callTool(ctx context.Context, server *Server, name string, args map[string]interface{}) Response {
	return server.HandleRequest(ctx, Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      name,
			"arguments": args,
		},
	})
}
