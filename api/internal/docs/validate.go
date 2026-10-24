package docs

import "errors"

var ErrNoFieldsToUpdate = errors.New("at least one of title, body, or status must be provided")
var ErrInvalidDocStatus = errors.New("invalid value for doc status")

var DocStatuses = map[string]bool{
	"proposed":   true,
	"approved":   true,
	"superseded": true,
	"rejected":   true,
}

func ValidateDocStatus(status string) error {
	if !DocStatuses[status] {
		return ErrInvalidDocStatus
	}
	return nil
}

func ValidateUpdateDocRequest(req UpdateDocRequest) error {
	if req.Title == nil && req.Body == nil && req.Status == nil {
		return ErrNoFieldsToUpdate
	}

	if req.Title != nil && *req.Title == "" {
		return errors.New("title cannot be empty")
	}

	if req.Body != nil && *req.Body == "" {
		return errors.New("body cannot be empty")
	}

	if req.Status != nil {
		if err := ValidateDocStatus(*req.Status); err != nil {
			return err
		}
	}

	return nil
}
