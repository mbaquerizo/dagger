package docs

import "errors"

var ErrDocNotFound = errors.New("doc not found")

type ParentDoc struct {
	DisplayID string
	Title     string
}

type Doc struct {
	ID        int
	DisplayID string
	DocType   string
	Title     string
	Body      *string
	Status    string
	Parent    *ParentDoc
}
