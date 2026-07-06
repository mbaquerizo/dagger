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
	if len(tools) != 5 {
		t.Fatalf("got %d tools, want 5", len(tools))
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
	want := []string{"get_issue", "get_doc", "list_issues", "update_issue_status", "publish"}
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

	t.Run("publish", func(t *testing.T) {
		tool, ok := byName["publish"]
		if !ok {
			t.Fatal("missing")
		}

		// Check top-level required params
		for _, name := range []string{"type", "title", "body"} {
			p, ok := tool.InputSchema.Properties[name]
			if !ok {
				t.Fatalf("missing %s param", name)
			}
			if !contains(tool.InputSchema.Required, name) {
				t.Errorf("%s should be required", name)
			}
			if p.Type != "string" {
				t.Errorf("%s type = %q, want string", name, p.Type)
			}
		}

		t.Run("project_id", func(t *testing.T) {
			p, ok := tool.InputSchema.Properties["project_id"]
			if !ok {
				t.Fatal("missing project_id param")
			}
			if !contains(tool.InputSchema.Required, "project_id") {
				t.Error("project_id should be required")
			}
			if p.Type != "number" {
				t.Errorf("project_id type = %q, want number", p.Type)
			}
		})

		t.Run("parent_id", func(t *testing.T) {
			p, ok := tool.InputSchema.Properties["parent_id"]
			if !ok {
				t.Fatal("missing parent_id param")
			}
			if contains(tool.InputSchema.Required, "parent_id") {
				t.Error("parent_id should be optional")
			}
			if p.Type != "number" {
				t.Errorf("parent_id type = %q, want number", p.Type)
			}
		})

		t.Run("metadata", func(t *testing.T) {
			p, ok := tool.InputSchema.Properties["metadata"]
			if !ok {
				t.Fatal("missing metadata param")
			}
			if contains(tool.InputSchema.Required, "metadata") {
				t.Error("metadata should be optional")
			}
			if p.Type != "object" {
				t.Fatalf("metadata type = %q, want object", p.Type)
			}
			if p.Properties == nil {
				t.Fatal("metadata.Properties is nil")
			}

			meta := p.Properties

			// issue_type
			ip, ok := meta.Properties["issue_type"]
			if !ok {
				t.Fatal("missing metadata.issue_type")
			}
			if contains(meta.Required, "issue_type") {
				t.Error("metadata.issue_type should be optional")
			}
			if ip.Type != "string" {
				t.Errorf("metadata.issue_type type = %q, want string", ip.Type)
			}

			// status
			sp, ok := meta.Properties["status"]
			if !ok {
				t.Fatal("missing metadata.status")
			}
			if sp.Type != "string" {
				t.Errorf("metadata.status type = %q, want string", sp.Type)
			}

			// tags
			tp, ok := meta.Properties["tags"]
			if !ok {
				t.Fatal("missing metadata.tags")
			}
			if tp.Type != "array" {
				t.Errorf("metadata.tags type = %q, want array", tp.Type)
			}
			if tp.Items == nil || tp.Items.Type != "string" {
				t.Error("metadata.tags.Items should be {Type: string}")
			}

			// relationships
			rp, ok := meta.Properties["relationships"]
			if !ok {
				t.Fatal("missing metadata.relationships")
			}
			if rp.Type != "array" {
				t.Errorf("metadata.relationships type = %q, want array", rp.Type)
			}
			if rp.Items == nil {
				t.Fatal("metadata.relationships.Items is nil")
			}
			if rp.Items.Type != "object" {
				t.Errorf("metadata.relationships.Items.Type = %q, want object", rp.Items.Type)
			}
			for _, fn := range []string{"target_id", "type"} {
				fp, ok := rp.Items.Properties[fn]
				if !ok {
					t.Fatalf("missing relationships.Items.Properties[%s]", fn)
				}
				if !contains(rp.Items.Required, fn) {
					t.Errorf("relationships.Items.Required missing %s", fn)
				}
				wantType := "string"
				if fn == "target_id" {
					wantType = "number"
				}
				if fp.Type != wantType {
					t.Errorf("relationships.Items.Properties[%s] type = %q, want %s", fn, fp.Type, wantType)
				}
			}
		})
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
