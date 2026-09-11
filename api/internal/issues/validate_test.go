package issues

import (
	"errors"
	"testing"
)

func Test_ValidateStatus_ValidStatus(t *testing.T) {
	validStatuses := []string{"open", "in-progress", "in-review", "done", "closed"}

	for _, status := range validStatuses {
		t.Run("valid status "+status, func(t *testing.T) {
			got := ValidateStatus(status)

			if got != nil {
				t.Errorf("Validate() = %q, want nil", got)
			}
		})
	}
}

func Test_ValidateStatus_InvalidStatus(t *testing.T) {
	err := ValidateStatus("i like turtles")

	if err == nil {
		t.Errorf("Validate() = %q, want nil", err)
	}

	if !errors.Is(err, ErrInvalidStatus) {
		t.Errorf("expected ErrInvalidStatus, but got %v", err)
	}
}
