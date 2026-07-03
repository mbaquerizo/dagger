package mcp

import (
	"encoding/json"
	"testing"
)

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{"ParseError", ErrCodeParse, -32700},
		{"InvalidRequest", ErrCodeInvalidRequest, -32600},
		{"MethodNotFound", ErrCodeMethodNotFound, -32601},
		{"InvalidParams", ErrCodeInvalidParams, -32602},
		{"InternalError", ErrCodeInternal, -32603},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
		}
	}
}

func TestResponseMarshalResult(t *testing.T) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      1,
		Result:  map[string]string{"status": "ok"},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	if m["jsonrpc"] != "2.0" {
		t.Error(`jsonrpc = missing or not "2.0"`)
	}
	if _, ok := m["error"]; ok {
		t.Error("error field present in result response")
	}
	result, ok := m["result"].(map[string]interface{})
	if !ok {
		t.Fatal("result field missing or not an object")
	}
	if result["status"] != "ok" {
		t.Error("result.status missing or wrong")
	}
}

func TestResponseMarshalError(t *testing.T) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      1,
		Error:   &Error{Code: -32602, Message: "Invalid params"},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	if _, ok := m["result"]; ok {
		t.Error("result field present in error response")
	}
	errObj, ok := m["error"].(map[string]interface{})
	if !ok {
		t.Fatal("error field missing or not an object")
	}
	if errObj["code"] != float64(-32602) {
		t.Error("error.code missing or wrong")
	}
	if errObj["message"] != "Invalid params" {
		t.Error("error.message missing or wrong")
	}
}

func TestListTools(t *testing.T) {
	tools := ListTools()
	if len(tools) != 4 {
		t.Fatalf("got %d tools, want 4", len(tools))
	}
	names := make(map[string]bool)
	for _, tool := range tools {
		if names[tool.Name] {
			t.Errorf("duplicate tool name: %s", tool.Name)
		}
		names[tool.Name] = true
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
		if tool.InputSchema.Type != "object" {
			t.Errorf("tool %q InputSchema.Type = %q, want object", tool.Name, tool.InputSchema.Type)
		}
	}
	want := []string{"get_issue", "get_doc", "list_issues", "update_issue_status"}
	for _, name := range want {
		if !names[name] {
			t.Errorf("missing tool: %s", name)
		}
	}
}

func TestToolParams(t *testing.T) {
	tools := ListTools()
	byName := make(map[string]ToolDefinition)
	for _, tool := range tools {
		byName[tool.Name] = tool
	}

	t.Run("get_issue", func(t *testing.T) {
		tool, ok := byName["get_issue"]
		if !ok {
			t.Fatal("missing")
		}
		p, ok := tool.InputSchema.Properties["display_id"]
		if !ok {
			t.Fatal("missing display_id param")
		}
		if p.Type != "string" {
			t.Errorf("display_id type = %q, want string", p.Type)
		}
		if !contains(tool.InputSchema.Required, "display_id") {
			t.Error("display_id should be required")
		}
	})

	t.Run("get_doc", func(t *testing.T) {
		tool, ok := byName["get_doc"]
		if !ok {
			t.Fatal("missing")
		}
		p, ok := tool.InputSchema.Properties["display_id"]
		if !ok {
			t.Fatal("missing display_id param")
		}
		if p.Type != "string" {
			t.Errorf("display_id type = %q, want string", p.Type)
		}
		if !contains(tool.InputSchema.Required, "display_id") {
			t.Error("display_id should be required")
		}
	})

	t.Run("list_issues", func(t *testing.T) {
		tool, ok := byName["list_issues"]
		if !ok {
			t.Fatal("missing")
		}
		_, ok = tool.InputSchema.Properties["status"]
		if !ok {
			t.Fatal("missing status param")
		}
		if contains(tool.InputSchema.Required, "status") {
			t.Error("status should be optional")
		}
	})

	t.Run("update_issue_status", func(t *testing.T) {
		tool, ok := byName["update_issue_status"]
		if !ok {
			t.Fatal("missing")
		}
		for _, name := range []string{"display_id", "status"} {
			p, ok := tool.InputSchema.Properties[name]
			if !ok {
				t.Fatalf("missing %s param", name)
			}
			if p.Type != "string" {
				t.Errorf("%s type = %q, want string", name, p.Type)
			}
			if !contains(tool.InputSchema.Required, name) {
				t.Errorf("%s should be required", name)
			}
		}
	})
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
