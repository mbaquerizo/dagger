package publish

import (
	"fmt"

	"github.com/mbaquerizo/dagger/internal/issues"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(req PublishRequest) []ValidationError {
	var errs []ValidationError

	switch req.Type {
	case "adr", "pitch", "ce", "issue":
	default:
		errs = append(errs, ValidationError{
			Field:   "type",
			Message: "type must be one of: adr, pitch, ce, issue",
		})
	}

	if req.Title == "" {
		errs = append(errs, ValidationError{
			Field:   "title",
			Message: "title is required",
		})
	}

	if req.Body == "" {
		errs = append(errs, ValidationError{
			Field:   "body",
			Message: "body is required",
		})
	}

	if req.ProjectID == 0 {
		errs = append(errs, ValidationError{
			Field:   "projectId",
			Message: "invalid projectId",
		})
	}

	if req.Type == "issue" {
		issueType := req.Metadata.IssueType

		if issueType == nil {
			errs = append(errs, ValidationError{
				Field:   "metadata.issueType",
				Message: "issue type is required for issues",
			})
		} else {
			switch *issueType {
			case "epic", "story", "task", "bug", "spike":
			default:
				errs = append(errs, ValidationError{
					Field:   "metadata.issueType",
					Message: "issue type must be one of: epic, story, task, bug, spike",
				})
			}
		}

		for i, relation := range req.Metadata.IssueRelations {
			if !issues.RelationTypes[relation.RelationType] {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("metadata.issueType[%d].relationType", i),
					Message: "issue relation must be one of: blocks, blocked_by, duplicates, duplicated_from, relates_to, causes, caused_by",
				})
			}
		}

	}

	if req.Type != "issue" && len(req.Metadata.IssueRelations) > 0 {
		errs = append(errs, ValidationError{
			Field:   "metadata.issueRelations",
			Message: "metadata.issueRelations only valid for type 'issue'",
		})
	}

	return errs
}
