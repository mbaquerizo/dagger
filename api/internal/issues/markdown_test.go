package issues

import (
	"strings"
	"testing"
)

func TestRenderIssueContext_Full(t *testing.T) {
	body := "Test body content\n\nWith multiple lines."
	adrBody := "ADR decision content"
	parentEpic := Issue{
		ID:        1,
		DisplayID: "EPIC-1",
		Title:     "API Hardening Epic",
	}

	ctx := &IssueContext{
		Issue: Issue{
			ID:        2,
			DisplayID: "DGR-42",
			TypeName:  "story",
			Title:     "Add rate limiting",
			Body:      &body,
			Status:    "open",
			ParentID:  &parentEpic.ID,
		},
		Parent: &parentEpic,
		LinkedDocs: []LinkedDoc{
			{
				ID:        1,
				DisplayID: "ADR-1",
				DocType:   "adr",
				Title:     "Rate limiting strategy",
				Body:      &adrBody,
				Status:    "approved",
			},
		},
		Children: []Issue{
			{DisplayID: "DGR-43", Title: "Implement token bucket"},
		},
		RelatedIssues: []RelatedIssue{
			{
				DisplayID:    "DGR-50",
				Title:        "Add Redis config",
				RelationType: "blocks",
			},
		},
	}

	result := RenderIssueContext(ctx)

	if !strings.HasPrefix(result, "---\n") {
		t.Errorf("expected frontmatter prefix, but got:\n%s", result)
	}

	if !strings.Contains(result, "id: 2") {
		t.Errorf("expected frontmatter id, got:\n%s", result)
	}
	if !strings.Contains(result, "display_id: DGR-42") {
		t.Errorf("expected frontmatter display_id, got:\n%s", result)
	}
	if !strings.Contains(result, "status: open") {
		t.Errorf("expected frontmatter status, got:\n%s", result)
	}
	if !strings.Contains(result, "type: story") {
		t.Errorf("expected frontmatter type, got:\n%s", result)
	}
	if !strings.Contains(result, "parent_id: 1") {
		t.Errorf("expected frontmatter parent_id, got:\n%s", result)
	}
	if !strings.Contains(result, "parent_display_id: EPIC-1") {
		t.Errorf("expected frontmatter parent_display_id, got:\n%s", result)
	}
	// Header
	if !strings.Contains(result, "# DGR-42: Add rate limiting") {
		t.Errorf("expected header line, got:\n%s", result)
	}
	if !strings.Contains(result, "**Status:** open") {
		t.Errorf("expected status line")
	}
	if !strings.Contains(result, "**Type:** story") {
		t.Errorf("expected type line")
	}
	if !strings.Contains(result, "**Parent:** EPIC-1 (API Hardening Epic)") {
		t.Errorf("expected parent line")
	}

	// Body
	if !strings.Contains(result, "Test body content") {
		t.Errorf("expected issue body")
	}

	// Separator
	if !strings.Contains(result, "\n---\n") {
		t.Errorf("expected separator")
	}

	// Linked Context
	if !strings.Contains(result, "## Linked Context") {
		t.Errorf("expected Linked Context section")
	}
	if !strings.Contains(result, "### ADR-1: Rate limiting strategy") {
		t.Errorf("expected linked doc header")
	}
	if !strings.Contains(result, "ADR decision content") {
		t.Errorf("expected linked doc body")
	}
	if !strings.Contains(result, "**Status:** approved") {
		t.Errorf("expected linked doc status")
	}

	// Subtasks
	if !strings.Contains(result, "## Subtasks") {
		t.Errorf("expected Subtasks section")
	}
	if !strings.Contains(result, "- DGR-43: Implement token bucket") {
		t.Errorf("expected subtask line")
	}

	// Related Issues
	if !strings.Contains(result, "## Related Issues") {
		t.Errorf("expected Related Issues section")
	}
	if !strings.Contains(result, "- **blocks** DGR-50: Add Redis config") {
		t.Errorf("expected related issue line")
	}
}

func TestRenderIssueContext_Minimal(t *testing.T) {
	ctx := &IssueContext{
		Issue: Issue{
			DisplayID: "DGR-1",
			TypeName:  "bug",
			Title:     "Simple bug",
			Body:      nil,
			Status:    "open",
		},
	}

	result := RenderIssueContext(ctx)

	if !strings.HasPrefix(result, "---\n") {
		t.Errorf("expected frontmatter prefix, but got:\n%s", result)
	}

	// Header present
	if !strings.Contains(result, "# DGR-1: Simple bug") {
		t.Errorf("expected header")
	}
	if !strings.Contains(result, "**Status:** open") {
		t.Errorf("expected status")
	}
	if !strings.Contains(result, "**Type:** bug") {
		t.Errorf("expected type")
	}

	// No optional sections
	if strings.Contains(result, "## Subtasks") {
		t.Errorf("expected no Subtasks section")
	}
	if strings.Contains(result, "## Related Issues") {
		t.Errorf("expected no Related Issues section")
	}
	if strings.Contains(result, "## Linked Context") {
		t.Errorf("expected no Linked Context section")
	}

	// No parent
	if strings.Contains(result, "**Parent:**") {
		t.Errorf("expected no parent line")
	}
}
