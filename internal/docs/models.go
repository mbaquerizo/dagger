package docs

import "errors"

var ErrDocNotFound = errors.New("doc not found")
var ErrProjectIDMismatch = errors.New("doc project_id does not match auth context")

type ParentDoc struct {
	DisplayID string
	Title     string
}

type Doc struct {
	ID          int
	DisplayID   string
	DocType     string
	Title       string
	Body        *string
	Status      string
	ProjectID   int
	WorkspaceID int
	Parent      *ParentDoc
}
