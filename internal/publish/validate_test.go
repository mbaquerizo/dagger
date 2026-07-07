package publish

import "testing"

func TestValidate_InvalidType(t *testing.T) {
	request := PublishRequest{
		Type:      "",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for empty type, got none")
	}
}

func TestValidate_MissingTitle(t *testing.T) {
	request := PublishRequest{
		Type:      "adr",
		Body:      "Test Body",
		ProjectID: 1,
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for missing title, got none")
	}
}

func TestValidate_MissingBody(t *testing.T) {
	request := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		ProjectID: 1,
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for missing body, got none")
	}
}

func TestValidate_InvalidProjectID(t *testing.T) {
	request := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 0,
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for invalid project id, got none")
	}
}

func TestValidate_MissingIssueType(t *testing.T) {
	request := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for missing issue type, got none")
	}
}

func TestValidate_InvalidIssueType(t *testing.T) {
	issueType := "lemon"

	request := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType: &issueType,
		},
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for invalid issue type, got none")
	}
}

func TestValidate_ValidIssueRequest(t *testing.T) {
	issueType := "story"

	request := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 3,
		Metadata: Metadata{
			IssueType: &issueType,
		},
	}

	errs := Validate(request)

	if len(errs) > 0 {
		t.Error("expected no validation errors for valid issue request")
	}
}

func TestValidate_ValidDocRequest(t *testing.T) {
	request := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 3,
	}

	errs := Validate(request)

	if len(errs) > 0 {
		t.Error("expected no validation errors for valid doc request")
	}
}

func TestValidate_IssueRelationsOnNonIssue(t *testing.T) {
	request := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueRelations: []IssueRelation{
				{TargetID: 5, RelationType: "blocks"},
			},
		},
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for issue_relations on non-issue type, got none")
	}
}

func TestValidate_InvalidRelationType(t *testing.T) {
	issueType := "story"

	request := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType: &issueType,
			IssueRelations: []IssueRelation{
				{TargetID: 5, RelationType: "invalid"},
			},
		},
	}

	errs := Validate(request)

	if len(errs) == 0 {
		t.Error("expected validation error for invalid relation type, got none")
	}
}

func TestValidate_ValidRelationTypes(t *testing.T) {
	validTypes := []string{"blocks", "blocked_by", "duplicates", "duplicated_from", "relates_to", "causes", "caused_by"}
	issueType := "story"

	var relations []IssueRelation

	for i, rt := range validTypes {
		relations = append(relations, IssueRelation{
			TargetID:     i + 1,
			RelationType: rt,
		})
	}

	request := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType:      &issueType,
			IssueRelations: relations,
		},
	}

	errs := Validate(request)

	if len(errs) > 0 {
		t.Errorf("expected no validation errors for valid relation types, got %d: %+v", len(errs), errs)
	}
}
