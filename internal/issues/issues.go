package issues

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrIssueNotFound = errors.New("issue not found")

type poolIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func queryLinkedDocs(ctx context.Context, pool poolIface, issueID int, workspaceID int) ([]LinkedDoc, error) {
	var linkedDocs []LinkedDoc

	linkedRows, err := pool.Query(ctx, `
		SELECT d.id, d.display_id, d.type, d.title, d.body, d.status FROM docs d
		JOIN doc_issues di ON d.id = di.doc_id
		WHERE di.issue_id = $1 AND d.workspace_id = $2
	`, issueID, workspaceID)

	if err != nil {
		return nil, fmt.Errorf("querying for linked docs: %w", err)
	}

	defer linkedRows.Close()

	for linkedRows.Next() {
		var doc LinkedDoc

		err := linkedRows.Scan(
			&doc.ID,
			&doc.DisplayID,
			&doc.DocType,
			&doc.Title,
			&doc.Body,
			&doc.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning linked doc: %w", err)
		}

		linkedDocs = append(linkedDocs, doc)
	}

	if err := linkedRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating linked doc: %w", err)
	}

	return linkedDocs, nil
}

func queryParent(ctx context.Context, pool poolIface, parentID *int, workspaceID int) (*Issue, error) {
	var parent *Issue

	if parentID != nil {
		var p Issue

		err := pool.QueryRow(ctx, `
		SELECT i.id, i.display_id, it.name, i.title, i.body, i.status, i.parent_id, i.project_id, i.workspace_id
		FROM issues i
		JOIN issue_types it ON i.issue_type_id = it.id 
		WHERE i.id = $1 AND i.workspace_id = $2
	`, *parentID, workspaceID).
			Scan(
				&p.ID,
				&p.DisplayID,
				&p.TypeName,
				&p.Title,
				&p.Body,
				&p.Status,
				&p.ParentID,
				&p.ProjectID,
				&p.WorkspaceID,
			)

		if err != nil {
			return nil, fmt.Errorf("querying for parent issue: %w", err)
		}

		parent = &p
	}

	return parent, nil
}

func queryChildren(ctx context.Context, pool poolIface, issueID int, workspaceID int) ([]Issue, error) {
	var children []Issue

	childRows, err := pool.Query(ctx, `
		SELECT i.id, i.display_id, it.name, i.title, i.body, i.status, i.parent_id, i.project_id, i.workspace_id
		FROM issues i
		JOIN issue_types it ON i.issue_type_id = it.id
		WHERE i.parent_id = $1 AND workspace_id = $2
	`, issueID, workspaceID)

	if err != nil {
		return nil, fmt.Errorf("querying for children: %w", err)
	}

	defer childRows.Close()

	for childRows.Next() {
		var child Issue

		err := childRows.Scan(
			&child.ID,
			&child.DisplayID,
			&child.TypeName,
			&child.Title,
			&child.Body,
			&child.Status,
			&child.ParentID,
			&child.ProjectID,
			&child.WorkspaceID,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning child row: %w", err)
		}

		children = append(children, child)
	}

	if err := childRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating children: %w", err)
	}

	return children, nil
}

func queryRelatedIssues(ctx context.Context, pool poolIface, issueID int, workspaceID int) ([]RelatedIssue, error) {
	var relatedIssues []RelatedIssue

	relatedRows, err := pool.Query(ctx, `
		SELECT i.display_id, i.title, r.name
		FROM issue_relations ir
		JOIN relations r ON r.id = ir.relation_id 
		JOIN issues i ON i.id = ir.target_issue_id
		WHERE ir.source_issue_id = $1 AND i.workspace_id = $2
	`, issueID, workspaceID)

	if err != nil {
		return nil, fmt.Errorf("querying for related rows: %w", err)
	}

	defer relatedRows.Close()

	for relatedRows.Next() {
		var relatedIssue RelatedIssue

		err := relatedRows.Scan(&relatedIssue.DisplayID, &relatedIssue.Title, &relatedIssue.RelationType)

		if err != nil {
			return nil, fmt.Errorf("scanning related row: %w", err)
		}

		relatedIssues = append(relatedIssues, relatedIssue)
	}

	if err := relatedRows.Err(); err != nil {
		return nil, fmt.Errorf("iterating related rows: %w", err)
	}

	return relatedIssues, nil
}

func GetIssueContext(ctx context.Context, pool poolIface, displayID string, workspaceID int) (*IssueContext, error) {
	var issue Issue

	err := pool.QueryRow(ctx, `
		SELECT i.id, i.display_id, it.name, i.title, i.body, i.status, i.parent_id, i.project_id, i.workspace_id
		FROM issues i
		JOIN issue_types it ON i.issue_type_id = it.id 
		WHERE display_id = $1 AND i.workspace_id = $2
	`, displayID, workspaceID).
		Scan(
			&issue.ID,
			&issue.DisplayID,
			&issue.TypeName,
			&issue.Title,
			&issue.Body,
			&issue.Status,
			&issue.ParentID,
			&issue.ProjectID,
			&issue.WorkspaceID,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIssueNotFound
		}

		return nil, fmt.Errorf("querying for issue: %w", err)
	}

	linkedDocs, err := queryLinkedDocs(ctx, pool, issue.ID, workspaceID)

	if err != nil {
		return nil, err
	}

	parent, err := queryParent(ctx, pool, issue.ParentID, workspaceID)

	if err != nil {
		return nil, err
	}

	children, err := queryChildren(ctx, pool, issue.ID, workspaceID)

	if err != nil {
		return nil, err
	}

	relatedIssues, err := queryRelatedIssues(ctx, pool, issue.ID, workspaceID)

	if err != nil {
		return nil, err
	}

	return &IssueContext{
		Issue:         issue,
		LinkedDocs:    linkedDocs,
		Parent:        parent,
		Children:      children,
		RelatedIssues: relatedIssues,
	}, nil
}

func ListIssues(ctx context.Context, pool poolIface, status string, workspaceID int) ([]IssueSummary, error) {
	rows, err := pool.Query(ctx, `
		SELECT i.display_id, it.name, i.title, i.status, p.display_id
		FROM issues i
		JOIN issue_types it ON it.id = i.issue_type_id
		LEFT JOIN issues p ON p.id = i.parent_id
		WHERE i.status = $1 AND i.workspace_id = $2
		ORDER BY i.id
	`, status, workspaceID)

	if err != nil {
		return nil, fmt.Errorf("querying issues: %w", err)
	}

	defer rows.Close()

	var issueSummaries []IssueSummary

	for rows.Next() {
		var issueSummary IssueSummary

		err := rows.Scan(
			&issueSummary.DisplayID,
			&issueSummary.TypeName,
			&issueSummary.Title,
			&issueSummary.Status,
			&issueSummary.ParentDisplayID,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning issue row: %w", err)
		}

		issueSummaries = append(issueSummaries, issueSummary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating issue rows: %w", err)
	}

	return issueSummaries, nil
}
