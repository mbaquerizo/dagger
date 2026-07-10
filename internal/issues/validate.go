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
