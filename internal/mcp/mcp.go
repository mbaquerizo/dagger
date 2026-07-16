package mcp

import (
	"context"

	"github.com/mbaquerizo/dagger/internal/issues"
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

type InitializeResult struct {
	ProtocolVersion    string             `json:"protocolVersion"`
	ServerCapabilities ServerCapabilities `json:"capabilities"`
	ServerInfo         ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools struct{} `json:"tools"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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
	Properties map[string]PropertySchema `json:"properties,omitempty"`
	Required   []string                  `json:"required,omitempty"`
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
	AddIssueRelation(ctx context.Context, req issues.AddIssueRelationRequest) (ToolResult, error)
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
			Description: "Fetch an issue by display ID with full context (linked docs, parent, children, related issues) as markdown. Includes YAML frontmatter with internal id, display_id, status, type, and parent info.",
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
			Description: "List all issues in the project as JSON. Optionally filter by status (open, in-progress, in-review, done, closed). When status is omitted, returns all issues. Each issue includes id, displayId, title, status, type, and parentDisplayId.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"status": {
						Type:        "string",
						Description: "Filter by status (open, in-progress, in-review, done, closed). Omit to return all issues.",
					},
				},
			},
		},
		{
			Name:        "update_issue_status",
			Description: "Update an issue's status by display ID. Valid status values: open, in-progress, in-review, done, closed.",
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
			Name:        "add_issue_relation",
			Description: "Create a bi-directional relationship between two issues by internal ID. Accepts source_id, target_id, and relation_type (blocks, blocked_by, duplicates, duplicated_from, relates_to, causes, caused_by). Returns the created relationship.",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]PropertySchema{
					"source_id": {
						Type:        "number",
						Description: "Internal ID of the source issue",
					},
					"target_id": {
						Type:        "number",
						Description: "Internal ID of the target issue",
					},
					"relation_type": {
						Type:        "string",
						Description: "Relation type (blocks, blocked_by, duplicates, duplicated_from, relates_to, causes, caused_by)",
					},
				},
				Required: []string{"source_id", "target_id", "relation_type"},
			},
		},
		{
			Name:        "get_doc",
			Description: "Fetch a document by display ID as markdown. Includes YAML frontmatter with internal id, display_id, status, and type.",
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
			Description: "Create a new document or issue. Accepts type (issue, adr, pitch, ce), title, body, project_id, optional parent_id, and nested metadata (issue_type, status, tags, relationships, issue_relations). Returns the created id, displayId, and url.",
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
