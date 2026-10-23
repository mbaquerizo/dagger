package issues

import "errors"

var ErrInvalidStatus = errors.New("invalid value for status")

var RelationTypes = map[string]bool{
	"blocks":          true,
	"blocked_by":      true,
	"duplicates":      true,
	"duplicated_from": true,
	"relates_to":      true,
	"causes":          true,
	"caused_by":       true,
}

var RelationInverse = map[string]string{
	"blocks":          "blocked_by",
	"blocked_by":      "blocks",
	"duplicates":      "duplicated_from",
	"duplicated_from": "duplicates",
	"relates_to":      "relates_to",
	"causes":          "caused_by",
	"caused_by":       "causes",
}

func ValidateStatus(status string) error {
	switch status {
	case "open", "in-progress", "in-review", "done", "closed":
		return nil
	default:
		return ErrInvalidStatus
	}
}

var ErrNoFieldsToUpdate = errors.New("at least one of title, body, or status must be provided")

func ValidateUpdateIssueRequest(req UpdateIssueRequest) error {
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
		if err := ValidateStatus(*req.Status); err != nil {
			return err
		}
	}

	return nil
}
