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
	Status        *string        `json:"status"`
	IssueType     *string        `json:"issueType"`
	Tags          []string       `json:"tags"`
	Relationships []Relationship `json:"relationships"`
}

type Relationship struct {
	TargetID int    `json:"targetId"`
	Type     string `json:"type"`
}

type PublishResponse struct {
	ID        int    `json:"id"`
	DisplayID string `json:"displayId"`
	URL       string `json:"url"`
}
