package publish

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
)

func TestPublish_Doc(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectBegin()
	mockPool.ExpectQuery(`SELECT slug FROM projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery(`UPDATE projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`INSERT INTO docs`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectCommit()

	req := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
	}

	resp, err := Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err != nil {
		t.Fatalf("Publish returned an error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if resp.ID != 84 {
		t.Errorf("expected ID 84, but got %d", resp.ID)
	}

	if resp.DisplayID != "DGR-47" {
		t.Errorf("expected display ID DGR-47, but got %s", resp.DisplayID)
	}
}

func TestPublish_Issue(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	issueType := "story"

	mockPool.ExpectBegin()
	mockPool.ExpectQuery(`SELECT slug FROM projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery(`UPDATE projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`SELECT id FROM issue_types`).
		WithArgs(issueType).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))
	mockPool.ExpectQuery(`INSERT INTO issues`).
		WithArgs(pgxmock.AnyArg(), 2, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectCommit()

	req := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType: &issueType,
		},
	}

	resp, err := Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err != nil {
		t.Fatalf("Publish returned an error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if resp.ID != 84 {
		t.Errorf("expected ID 84, but got %d", resp.ID)
	}

	if resp.DisplayID != "DGR-47" {
		t.Errorf("expected display ID DGR-47, but got %s", resp.DisplayID)
	}
}

func TestPublish_DocWithParent(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	parentID := 5

	mockPool.ExpectBegin()
	mockPool.ExpectQuery(`SELECT slug FROM projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery(`UPDATE projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`SELECT EXISTS.*FROM docs`).
		WithArgs(&parentID, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mockPool.ExpectQuery(`INSERT INTO docs`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectCommit()

	req := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ParentID:  &parentID,
		ProjectID: 1,
	}

	resp, err := Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err != nil {
		t.Fatalf("Publish returned an error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if resp.ID != 84 {
		t.Errorf("expected ID 84, but got %d", resp.ID)
	}

	if resp.DisplayID != "DGR-47" {
		t.Errorf("expected display ID DGR-47, but got %s", resp.DisplayID)
	}
}

func TestPublish_InvalidIssueType(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectBegin()
	mockPool.ExpectQuery(`SELECT slug FROM projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery(`UPDATE projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`SELECT id FROM issue_types`).
		WithArgs("moogle").
		WillReturnError(pgx.ErrNoRows)

	issueType := "moogle"

	req := PublishRequest{
		Type:      "issue",
		Title:     "Test Title",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType: &issueType,
		},
	}

	_, err = Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err == nil {
		t.Errorf("expected error for missing issue type, but got none")
	}
}

func TestPublish_MissingParent(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	parentID := 5

	mockPool.ExpectBegin()
	mockPool.ExpectQuery(`SELECT slug FROM projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery(`UPDATE projects`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`SELECT EXISTS.*FROM docs`).
		WithArgs(&parentID, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	req := PublishRequest{
		Type:      "adr",
		Title:     "Test Title",
		Body:      "Test Body",
		ParentID:  &parentID,
		ProjectID: 1,
	}

	_, err = Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err == nil {
		t.Errorf("expected error for missing parent, but got none")
	}
}

func TestPublish_IssueWithRelationships(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	issueType := "story"

	mockPool.ExpectBegin()
	mockPool.ExpectQuery("SELECT slug FROM projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery("UPDATE projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery("SELECT id FROM issue_types").
		WithArgs(issueType).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))
	mockPool.ExpectQuery("INSERT INTO issues").
		WithArgs(pgxmock.AnyArg(), 2, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectQuery("SELECT EXISTS.*FROM docs").
		WithArgs(42, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mockPool.ExpectExec("INSERT INTO doc_issues").
		WithArgs(42, 84, "motivates").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mockPool.ExpectCommit()

	req := PublishRequest{
		Type:      "issue",
		Title:     "Test Issue",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType: &issueType,
			Relationships: []Relationship{
				{TargetID: 42, Type: "motivates"},
			},
		},
	}

	resp, err := Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err != nil {
		t.Fatalf("Publish returned an error: %v", err)
	}

	if resp.ID != 84 {
		t.Errorf("expected ID 84, but got %d", resp.ID)
	}

	if resp.DisplayID != "DGR-47" {
		t.Errorf("expected DisplayID DGR-47, but got %s", resp.DisplayID)
	}
}

func TestPublish_IssueWithMissingRelationshipTarget(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	issueType := "story"

	mockPool.ExpectBegin()
	mockPool.ExpectQuery("SELECT slug FROM projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery("UPDATE projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery("SELECT id FROM issue_types").
		WithArgs(issueType).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(2))
	mockPool.ExpectQuery("INSERT INTO issues").
		WithArgs(pgxmock.AnyArg(), 2, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectQuery("SELECT EXISTS.*FROM docs").
		WithArgs(999, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	req := PublishRequest{
		Type:      "issue",
		Title:     "Test Issue",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			IssueType: &issueType,
			Relationships: []Relationship{
				{TargetID: 999, Type: "motivates"},
			},
		},
	}

	_, err = Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err == nil {
		t.Errorf("expected error for missing relationship target, but got none")
	}
}

func TestPublish_DocWithRelationships(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectBegin()
	mockPool.ExpectQuery("SELECT slug FROM projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery("UPDATE projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery(`INSERT INTO docs`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectQuery("SELECT EXISTS.*FROM issues").
		WithArgs(42, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mockPool.ExpectExec("INSERT INTO doc_issues").
		WithArgs(84, 42, "motivates").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mockPool.ExpectCommit()

	req := PublishRequest{
		Type:      "adr",
		Title:     "Test Issue",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			Relationships: []Relationship{
				{TargetID: 42, Type: "motivates"},
			},
		},
	}

	resp, err := Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err != nil {
		t.Fatalf("Publish returned an error: %v", err)
	}

	if resp.ID != 84 {
		t.Errorf("expected ID 84, but got %d", resp.ID)
	}

	if resp.DisplayID != "DGR-47" {
		t.Errorf("expected DisplayID DGR-47, but got %s", resp.DisplayID)
	}
}

func TestPublish_DocWithMissingRelationshipTarget(t *testing.T) {
	mockPool, err := pgxmock.NewPool()

	if err != nil {
		t.Errorf("failed to create mock pool: %v", err)
	}

	t.Cleanup(func() { mockPool.Close() })

	mockPool.ExpectBegin()
	mockPool.ExpectQuery("SELECT slug FROM projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"slug"}).AddRow("DGR"))
	mockPool.ExpectQuery("UPDATE projects").
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"next_display_number"}).AddRow(47))
	mockPool.ExpectQuery("INSERT INTO docs").
		WithArgs(pgxmock.AnyArg(), 2, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(84))
	mockPool.ExpectQuery("SELECT EXISTS.*FROM issue").
		WithArgs(999, 1).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	req := PublishRequest{
		Type:      "adr",
		Title:     "Test Issue",
		Body:      "Test Body",
		ProjectID: 1,
		Metadata: Metadata{
			Relationships: []Relationship{
				{TargetID: 999, Type: "motivates"},
			},
		},
	}

	_, err = Publish(context.Background(), mockPool, req, 1, "localhost:8080")

	if err == nil {
		t.Errorf("expected error for missing relationship target, but got none")
	}
}
