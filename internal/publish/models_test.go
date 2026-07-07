package publish

import (
	"encoding/json"
	"testing"
)

func TestIssueRelation_JSONTags(t *testing.T) {
	rel := IssueRelation{
		TargetID:     5,
		RelationType: "blocks",
	}

	data, err := json.Marshal(rel)

	if err != nil {
		t.Fatalf("failed to marshal IssueRelation: %v", err)
	}

	var decoded IssueRelation
	err = json.Unmarshal(data, &decoded)

	if err != nil {
		t.Fatalf("failed to unmarshal IssueRelation: %v", err)
	}

	if decoded.TargetID != 5 {
		t.Errorf("expected TargetID 5, got %d", decoded.TargetID)
	}

	if decoded.RelationType != "blocks" {
		t.Errorf("expected RelationType 'blocks', got '%s'", decoded.RelationType)
	}
}

func TestMetadata_IssueRelationsField(t *testing.T) {
	md := Metadata{
		IssueRelations: []IssueRelation{
			{TargetID: 5, RelationType: "blocks"},
			{TargetID: 3, RelationType: "relates_to"},
		},
	}

	data, err := json.Marshal(md)

	if err != nil {
		t.Fatalf("failed to marshal Metadata: %v", err)
	}

	var decoded Metadata
	err = json.Unmarshal(data, &decoded)

	if err != nil {
		t.Fatalf("failed to unmarshal Metadata: %v", err)
	}

	if len(decoded.IssueRelations) != 2 {
		t.Fatalf("expected 2 issue relations, got %d", len(decoded.IssueRelations))
	}

	if decoded.IssueRelations[0].TargetID != 5 || decoded.IssueRelations[0].RelationType != "blocks" {
		t.Errorf("first relation mismatch: got %+v", decoded.IssueRelations[0])
	}

	if decoded.IssueRelations[1].TargetID != 3 || decoded.IssueRelations[1].RelationType != "relates_to" {
		t.Errorf("second relation mismatch: got %+v", decoded.IssueRelations[1])
	}
}

func TestMetadata_ExistingRelationshipsUntouched(t *testing.T) {
	relType := "story"
	status := "open"

	md := Metadata{
		Status:    &status,
		IssueType: &relType,
		Tags:      []string{"mcp", "publish"},
		Relationships: []Relationship{
			{TargetID: 10, Type: "motivates"},
		},
		IssueRelations: []IssueRelation{
			{TargetID: 5, RelationType: "blocks"},
		},
	}

	data, err := json.Marshal(md)

	if err != nil {
		t.Fatalf("failed to marshal Metadata: %v", err)
	}

	var decoded Metadata
	err = json.Unmarshal(data, &decoded)

	if err != nil {
		t.Fatalf("failed to unmarshal Metadata: %v", err)
	}

	if len(decoded.Relationships) != 1 {
		t.Errorf("expected 1 relationship, got %d", len(decoded.Relationships))
	}

	if decoded.Relationships[0].TargetID != 10 || decoded.Relationships[0].Type != "motivates" {
		t.Errorf("relationship mismatch: got %+v", decoded.Relationships[0])
	}

	if len(decoded.IssueRelations) != 1 {
		t.Errorf("expected 1 issue relation, got %d", len(decoded.IssueRelations))
	}
}
