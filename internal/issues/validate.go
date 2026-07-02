package issues

import "errors"

var ErrInvalidStatus = errors.New("invalid value for status")

func ValidateStatus(status string) error {
	switch status {
	case "open", "in-progress", "in-review", "done", "closed":
		return nil
	default:
		return ErrInvalidStatus
	}
}
