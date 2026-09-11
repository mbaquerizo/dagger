package issues

type Issue struct {
	ID          int
	DisplayID   string
	TypeName    string
	Title       string
	Body        *string
	Status      string
	ParentID    *int
	ProjectID   int
	WorkspaceID int
}

type LinkedDoc struct {
	ID        int
	DisplayID string
	DocType   string
	Title     string
	Body      *string
	Status    string
}

type RelatedIssue struct {
	DisplayID    string
	Title        string
	RelationType string
}

type IssueContext struct {
	Issue         Issue
	LinkedDocs    []LinkedDoc
	Parent        *Issue
	Children      []Issue
	RelatedIssues []RelatedIssue
}

type IssueSummary struct {
	ID              int     `json:"id"`
	DisplayID       string  `json:"displayId"`
	Title           string  `json:"title"`
	Status          string  `json:"status"`
	TypeName        string  `json:"type"`
	ParentDisplayID *string `json:"parentDisplayId,omitempty"`
}

type UpdateStatusRequest struct {
	// Valid values: "open", "in-progress", "in-review", "done", "closed"
	Status string `json:"status"`
}

type AddIssueRelationRequest struct {
	SourceID     int    `json:"sourceId"`
	TargetID     int    `json:"targetId"`
	RelationType string `json:"relationType"`
}
