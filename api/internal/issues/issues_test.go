package issues

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
)

func TestGetIssueContext_FullSuccess(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	issueBody := "Test Body"
	adrBody := "ADR Body"
	parentID := 10

	// Issue query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-42", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-42", "story", "Test Issue", &issueBody, "open", &parentID, 1, 1))

	// Linked docs query
	mockPool.ExpectQuery(`SELECT d.id, d.display_id, d.type`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}).
			AddRow(5, "ADR-1", "adr", "Test ADR", &adrBody, "approved"))

	// Parent query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(10, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(10, "EPIC-1", "epic", "Parent Epic", nil, "open", nil, 1, 1))

	childParentID := 1

	// Children query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(2, "DGR-43", "task", "Child Task", nil, "open", &childParentID, 1, 1))

	// Related issues query
	mockPool.ExpectQuery(`SELECT i.display_id, i.title, r.name`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}).
			AddRow("DGR-50", "Related Issue", "blocks"))

	ctx, err := GetIssueContext(context.Background(), mockPool, "DGR-42", 1, nil)
	if err != nil {
		t.Fatalf("GetIssueContext returned error: %v", err)
	}
	if ctx == nil {
		t.Fatal("expected non-nil IssueContext")
	}

	// Issue
	if ctx.Issue.DisplayID != "DGR-42" {
		t.Errorf("expected DisplayID DGR-42, got %s", ctx.Issue.DisplayID)
	}
	if ctx.Issue.TypeName != "story" {
		t.Errorf("expected TypeName story, got %s", ctx.Issue.TypeName)
	}
	if ctx.Issue.Title != "Test Issue" {
		t.Errorf("expected Title 'Test Issue', got %s", ctx.Issue.Title)
	}
	if ctx.Issue.Body == nil || *ctx.Issue.Body != "Test Body" {
		t.Errorf("expected Body 'Test Body', got %v", ctx.Issue.Body)
	}
	if ctx.Issue.Status != "open" {
		t.Errorf("expected Status open, got %s", ctx.Issue.Status)
	}
	if ctx.Issue.ParentID == nil || *ctx.Issue.ParentID != 10 {
		t.Errorf("expected ParentID 10, got %v", ctx.Issue.ParentID)
	}

	// Linked docs
	if len(ctx.LinkedDocs) != 1 {
		t.Fatalf("expected 1 linked doc, got %d", len(ctx.LinkedDocs))
	}
	if ctx.LinkedDocs[0].DisplayID != "ADR-1" {
		t.Errorf("expected ADR-1, got %s", ctx.LinkedDocs[0].DisplayID)
	}
	if ctx.LinkedDocs[0].DocType != "adr" {
		t.Errorf("expected type adr, got %s", ctx.LinkedDocs[0].DocType)
	}

	// Parent
	if ctx.Parent == nil {
		t.Fatal("expected non-nil parent")
	}
	if ctx.Parent.DisplayID != "EPIC-1" {
		t.Errorf("expected parent EPIC-1, got %s", ctx.Parent.DisplayID)
	}

	// Children
	if len(ctx.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(ctx.Children))
	}
	if ctx.Children[0].DisplayID != "DGR-43" {
		t.Errorf("expected child DGR-43, got %s", ctx.Children[0].DisplayID)
	}

	// Related issues
	if len(ctx.RelatedIssues) != 1 {
		t.Fatalf("expected 1 related issue, got %d", len(ctx.RelatedIssues))
	}
	if ctx.RelatedIssues[0].DisplayID != "DGR-50" {
		t.Errorf("expected related DGR-50, got %s", ctx.RelatedIssues[0].DisplayID)
	}
	if ctx.RelatedIssues[0].RelationType != "blocks" {
		t.Errorf("expected relation type 'blocks', got %s", ctx.RelatedIssues[0].RelationType)
	}
}

func TestGetIssueContext_NotFound(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-999", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	ctx, err := GetIssueContext(context.Background(), mockPool, "DGR-999", 1, nil)
	if err == nil {
		t.Fatal("expected ErrIssueNotFound, got nil")
	}
	if ctx != nil {
		t.Errorf("expected nil context, got %v", ctx)
	}
}

func TestGetIssueContext_NoExtras(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	// Issue query — no parent_id
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-1", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-1", "bug", "Standalone Bug", nil, "open", nil, 1, 1))

	// Linked docs — empty
	mockPool.ExpectQuery(`SELECT d.id, d.display_id, d.type`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}))

	// No parent query (parent_id is nil)

	// Children — empty
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	// Related issues — empty
	mockPool.ExpectQuery(`SELECT i.display_id, i.title, r.name`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}))

	ctx, err := GetIssueContext(context.Background(), mockPool, "DGR-1", 1, nil)
	if err != nil {
		t.Fatalf("GetIssueContext returned error: %v", err)
	}
	if ctx == nil {
		t.Fatal("expected non-nil IssueContext")
	}

	if ctx.Issue.DisplayID != "DGR-1" {
		t.Errorf("expected DGR-1, got %s", ctx.Issue.DisplayID)
	}
	if ctx.Issue.TypeName != "bug" {
		t.Errorf("expected TypeName bug, got %s", ctx.Issue.TypeName)
	}

	if len(ctx.LinkedDocs) != 0 {
		t.Errorf("expected 0 linked docs, got %d", len(ctx.LinkedDocs))
	}
	if ctx.Parent != nil {
		t.Errorf("expected nil parent, got %v", ctx.Parent)
	}
	if len(ctx.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(ctx.Children))
	}
	if len(ctx.RelatedIssues) != 0 {
		t.Errorf("expected 0 related issues, got %d", len(ctx.RelatedIssues))
	}
}

func TestGetIssueContext_ParentQueryUsesQualifiedID(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	parentID := 10

	// Issue with parent_id = 10
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-TEST", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(5, "DGR-TEST", "story", "Child", nil, "open", &parentID, 1, 1))

	// Linked docs — empty
	mockPool.ExpectQuery(`SELECT d.id, d.display_id, d.type`).
		WithArgs(5, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}))

	// Parent query — must use i.id =, not bare id =, to avoid ambiguity
	mockPool.ExpectQuery(`WHERE i\.id = \$1 AND i\.workspace_id = \$2`).
		WithArgs(10, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(10, "EPIC-1", "epic", "Parent", nil, "open", nil, 1, 1))

	// Children — empty
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(5, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}))

	// Related issues — empty
	mockPool.ExpectQuery(`SELECT i.display_id, i.title, r.name`).
		WithArgs(5, 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}))

	ctx, err := GetIssueContext(context.Background(), mockPool, "DGR-TEST", 1, nil)
	if err != nil {
		t.Fatalf("GetIssueContext returned error: %v", err)
	}
	if ctx == nil {
		t.Fatal("expected non-nil IssueContext")
	}
	if ctx.Parent == nil {
		t.Fatal("expected non-nil parent")
	}
	if ctx.Parent.DisplayID != "EPIC-1" {
		t.Errorf("expected parent EPIC-1, got %s", ctx.Parent.DisplayID)
	}
}

func TestGetIssueContext_ScopedToken(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	issueBody := "Test Body"
	adrBody := "ADR Body"
	parentID := 10
	projectID := 1

	// Issue query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-42", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-42", "story", "Test Issue", &issueBody, "open", &parentID, 1, 1))

	// Linked docs query
	mockPool.ExpectQuery(`SELECT d.id, d.display_id, d.type`).
		WithArgs(1, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type", "title", "body", "status"}).
			AddRow(5, "ADR-1", "adr", "Test ADR", &adrBody, "approved"))

	// Parent query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(10, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(10, "EPIC-1", "epic", "Parent Epic", nil, "open", nil, 1, 1))

	childParentID := 1

	// Children query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs(1, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(2, "DGR-43", "task", "Child Task", nil, "open", &childParentID, 1, 1))

	// Related issues query
	mockPool.ExpectQuery(`SELECT i.display_id, i.title, r.name`).
		WithArgs(1, 1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"display_id", "title", "relation_type"}).
			AddRow("DGR-50", "Related Issue", "blocks"))

	ctx, err := GetIssueContext(context.Background(), mockPool, "DGR-42", 1, &projectID)

	if err != nil {
		t.Fatalf("GetIssueContext returned error: %v", err)
	}
	if ctx == nil {
		t.Fatal("expected non-nil IssueContext")
	}

	// Issue
	if ctx.Issue.DisplayID != "DGR-42" {
		t.Errorf("expected DisplayID DGR-42, got %s", ctx.Issue.DisplayID)
	}
	if ctx.Issue.TypeName != "story" {
		t.Errorf("expected TypeName story, got %s", ctx.Issue.TypeName)
	}
	if ctx.Issue.Title != "Test Issue" {
		t.Errorf("expected Title 'Test Issue', got %s", ctx.Issue.Title)
	}
	if ctx.Issue.Body == nil || *ctx.Issue.Body != "Test Body" {
		t.Errorf("expected Body 'Test Body', got %v", ctx.Issue.Body)
	}
	if ctx.Issue.Status != "open" {
		t.Errorf("expected Status open, got %s", ctx.Issue.Status)
	}
	if ctx.Issue.ParentID == nil || *ctx.Issue.ParentID != 10 {
		t.Errorf("expected ParentID 10, got %v", ctx.Issue.ParentID)
	}

	// Linked docs
	if len(ctx.LinkedDocs) != 1 {
		t.Fatalf("expected 1 linked doc, got %d", len(ctx.LinkedDocs))
	}
	if ctx.LinkedDocs[0].DisplayID != "ADR-1" {
		t.Errorf("expected ADR-1, got %s", ctx.LinkedDocs[0].DisplayID)
	}
	if ctx.LinkedDocs[0].DocType != "adr" {
		t.Errorf("expected type adr, got %s", ctx.LinkedDocs[0].DocType)
	}

	// Parent
	if ctx.Parent == nil {
		t.Fatal("expected non-nil parent")
	}
	if ctx.Parent.DisplayID != "EPIC-1" {
		t.Errorf("expected parent EPIC-1, got %s", ctx.Parent.DisplayID)
	}

	// Children
	if len(ctx.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(ctx.Children))
	}
	if ctx.Children[0].DisplayID != "DGR-43" {
		t.Errorf("expected child DGR-43, got %s", ctx.Children[0].DisplayID)
	}

	// Related issues
	if len(ctx.RelatedIssues) != 1 {
		t.Fatalf("expected 1 related issue, got %d", len(ctx.RelatedIssues))
	}
	if ctx.RelatedIssues[0].DisplayID != "DGR-50" {
		t.Errorf("expected related DGR-50, got %s", ctx.RelatedIssues[0].DisplayID)
	}
	if ctx.RelatedIssues[0].RelationType != "blocks" {
		t.Errorf("expected relation type 'blocks', got %s", ctx.RelatedIssues[0].RelationType)
	}
}

func TestGetIssueContext_ScopedTokenProjectIDMismatch(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	issueBody := "Test Body"
	parentID := 10
	projectID := 1

	// Issue query
	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("DGR-42", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "type_name", "title", "body", "status", "parent_id", "project_id", "workspace_id"}).
			AddRow(1, "DGR-42", "story", "Test Issue", &issueBody, "open", &parentID, 2, 1))

	ctx, err := GetIssueContext(context.Background(), mockPool, "DGR-42", 1, &projectID)

	if err != nil {
		if !errors.Is(err, ErrProjectIDMismatch) {
			t.Errorf("expected project id mismatch error, but got %v", err)
		}
	}

	if ctx != nil {
		t.Error("expected no context, but got context")
	}
}

func TestListIssues_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	parentDisplayID := "EPIC-1"

	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("open", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "title", "status", "type_name", "parent_display_id"}).
			AddRow(1, "DGR-1", "First issue", "open", "story", nil).
			AddRow(2, "DGR-2", "Second issue", "open", "bug", &parentDisplayID))

	issues, err := ListIssues(context.Background(), mockPool, "open", 1, nil)
	if err != nil {
		t.Fatalf("ListIssues returned error: %v", err)
	}
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}

	if issues[0].DisplayID != "DGR-1" {
		t.Errorf("expected DGR-1, got %s", issues[0].DisplayID)
	}
	if issues[0].ParentDisplayID != nil {
		t.Errorf("expected nil parent for first issue, got %v", issues[0].ParentDisplayID)
	}

	if issues[1].DisplayID != "DGR-2" {
		t.Errorf("expected DGR-2, got %s", issues[1].DisplayID)
	}
	if issues[1].ParentDisplayID == nil || *issues[1].ParentDisplayID != "EPIC-1" {
		t.Errorf("expected parent EPIC-1, got %v", issues[1].ParentDisplayID)
	}
}

func TestListIssues_NoResults(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery(`SELECT i.id, i.display_id, it.name`).
		WithArgs("closed", 1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "display_id", "title", "status", "type_name", "parent_display_id"}))

	issues, err := ListIssues(context.Background(), mockPool, "closed", 1, nil)
	if err != nil {
		t.Fatalf("ListIssues returned error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestUpdateIssueStatus_ValidStatus(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectExec(`UPDATE issues`).
		WithArgs("in-progress", "DGR-45", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	req := UpdateStatusRequest{
		Status: "in-progress",
	}

	err = UpdateIssueStatus(context.Background(), mockPool, req, "DGR-45", 1, nil)

	if err != nil {
		t.Errorf("expected success, but got error: %v", err)
	}
}

func TestAddIssueRelation_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery("SELECT EXISTS.*FROM issues").
		WithArgs(5, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mockPool.ExpectQuery("SELECT EXISTS.*FROM issues").
		WithArgs(6, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mockPool.ExpectExec("INSERT INTO issue_relations").
		WithArgs(5, 6, "blocks", "blocked_by").
		WillReturnResult(pgxmock.NewResult("INSERT", 2))

	err = AddIssueRelation(context.Background(), mockPool, 5, 6, "blocks", 1, nil)

	if err != nil {
		t.Errorf("expected success, but got error: %v", err)
	}
}

func TestAddIssueRelation_SelfRelation(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	err = AddIssueRelation(context.Background(), mockPool, 5, 5, "blocks", 1, nil)

	if err == nil {
		t.Error("expected error for self-relation, got none")
	}
}

func TestAddIssueRelation_SourceNotFound(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery("SELECT EXISTS.*FROM issues").
		WithArgs(999, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	err = AddIssueRelation(context.Background(), mockPool, 999, 5, "blocks", 1, nil)

	if !errors.Is(err, ErrIssueNotFound) {
		t.Errorf("expected ErrIssueNotFound, got: %v", err)
	}
}

func TestAddIssueRelation_TargetNotFound(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectQuery("SELECT EXISTS.*FROM issues").
		WithArgs(5, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mockPool.ExpectQuery("SELECT EXISTS.*FROM issues").
		WithArgs(999, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	err = AddIssueRelation(context.Background(), mockPool, 5, 999, "blocks", 1, nil)

	if !errors.Is(err, ErrIssueNotFound) {
		t.Errorf("expected ErrIssueNotFound, got: %v", err)
	}
}
