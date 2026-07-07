package publish

type PublishRequest struct {
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	ParentID  *int     `json:"parentId"`
	ProjectID int      `json:"projectId"`
	Metadata  Metadata `json:"metadata"`
}

type Metadata struct {
	Status         *string         `json:"status"`
	IssueType      *string         `json:"issueType"`
	Tags           []string        `json:"tags"`
	Relationships  []Relationship  `json:"relationships"`
	IssueRelations []IssueRelation `json:"issueRelations"`
}

type Relationship struct {
	TargetID int    `json:"targetId"`
	Type     string `json:"type"`
}

type IssueRelation struct {
	TargetID     int    `json:"targetId"`
	RelationType string `json:"relationType"`
}

type PublishResponse struct {
	ID        int    `json:"id"`
	DisplayID string `json:"displayId"`
	URL       string `json:"url"`
}
