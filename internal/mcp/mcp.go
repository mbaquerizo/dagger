package mcp

import "context"

const (
	ErrCodeParse          = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternal       = -32603
)

type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]PropertySchema `json:"properties"`
	Required   []string                  `json:"required"`
}

type PropertySchema struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolService interface {
	GetIssue(ctx context.Context, displayID string) (ToolResult, error)
	GetDoc(ctx context.Context, displayID string) (ToolResult, error)
	ListIssues(ctx context.Context, status string) (ToolResult, error)
	UpdateIssueStatus(ctx context.Context, displayID string, newStatus string) (ToolResult, error)
}

type ToolResult struct {
	Content []ContentItem `json:"content"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func ListTools() []ToolDefinition {
	return []ToolDefinition{
		{
			Name:        "get_issue",
			Description: "Fetch an issue with full context as markdown",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"display_id": {
						Type:        "string",
						Description: "Issue display ID (e.g. DGR-42)",
					},
				},
				Required: []string{"display_id"},
			},
		},
		{
			Name:        "list_issues",
			Description: "List issues with optional status filter",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"status": {
						Type:        "string",
						Description: "Filter by status (open, in-progress, in-review, done, closed)",
					},
				},
			},
		},
		{
			Name:        "update_issue_status",
			Description: "Update an issue's status",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"display_id": {
						Type:        "string",
						Description: "Issue display ID (e.g. DGR-42)",
					},
					"status": {
						Type:        "string",
						Description: "New status (open, in-progress, in-review, done, closed)",
					},
				},
				Required: []string{"display_id", "status"},
			},
		},
		{
			Name:        "get_doc",
			Description: "Fetch a document as markdown",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"display_id": {
						Type:        "string",
						Description: "Document display ID (e.g. DGR-36)",
					},
				},
				Required: []string{"display_id"},
			},
		},
	}
}
