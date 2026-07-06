package mcp

import (
	"context"

	"github.com/mbaquerizo/dagger/internal/publish"
)

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
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Items       *InputSchema `json:"items,omitempty"`
	Properties  *InputSchema `json:"properties,omitempty"`
}

type ToolService interface {
	GetIssue(ctx context.Context, displayID string) (ToolResult, error)
	GetDoc(ctx context.Context, displayID string) (ToolResult, error)
	ListIssues(ctx context.Context, status string) (ToolResult, error)
	UpdateIssueStatus(ctx context.Context, displayID string, newStatus string) (ToolResult, error)
	Publish(ctx context.Context, req publish.PublishRequest) (ToolResult, error)
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
		{
			Name:        "publish",
			Description: "Publish a document (including issues)",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"type": {
						Type:        "string",
						Description: "Type of document to be published, one of: adr, pitch, ce, issue",
					},
					"title": {
						Type:        "string",
						Description: "Title of document",
					},
					"body": {
						Type:        "string",
						Description: "Body of document in markdown format",
					},
					"project_id": {
						Type:        "number",
						Description: "ID of Dagger project to publish the document to",
					},
					"parent_id": {
						Type:        "number",
						Description: "ID of parent document",
					},
					"metadata": {
						Type:        "object",
						Description: "Optional document metadata",
						Properties: &InputSchema{
							Type: "object",
							Properties: map[string]PropertySchema{
								"issue_type": {
									Type:        "string",
									Description: "Type of issue, one of: epic, story, task, bug, spike. Required when document type is 'issue'",
								},
								"status": {
									Type:        "string",
									Description: "Optional initial document status. Defaults to 'open' if issue, 'proposed' if document",
								},
								"tags": {
									Type:        "array",
									Description: "Array of strings representing keywords related to the document",
									Items:       &InputSchema{Type: "string"},
								},
								"relationships": {
									Type:        "array",
									Description: "Array of {target_id,type} objects describing related cross-type documents. Not for doc↔doc or issue↔issue relationships",
									Items: &InputSchema{
										Type: "object",
										Properties: map[string]PropertySchema{
											"target_id": {
												Type:        "number",
												Description: "ID of related document",
											},
											"type": {
												Type:        "string",
												Description: "Type of relationship",
											},
										},
										Required: []string{"target_id", "type"},
									},
								},
							},
						},
					},
				},
				Required: []string{"type", "title", "body", "project_id"},
			},
		},
	}
}
