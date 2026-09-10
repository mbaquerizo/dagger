---
id: DGR-19
issueType: story
status: done
tags:
  - mcp
  - api
blockedBy: []
parent: DGR-5
---

# Show internal DB ID in rendered markdown outputs

## Description

As an AI agent using the MCP tools, I want the rendered markdown from `get_issue` and `get_doc` to include YAML frontmatter with internal DB IDs, so that I can reference issues and docs by their internal ID when calling `publish`'s `parent_id` parameter or the `add_issue_relation` tool.

Currently, agents receive markdown with only display IDs (e.g. `DGR-42`), but `parent_id` and the new `add_issue_relation` tool require internal DB IDs. Similarly, `list_issues` returns JSON without `id`, forcing agents to guess or use display IDs which won't work.

## Acceptance criteria

1. `get_issue` output starts with YAML frontmatter containing `id`, `display_id`, `status`, `type`, and optionally `parent_id`/`parent_display_id`
2. `get_doc` output starts with same frontmatter
3. `list_issues` JSON includes an `id` field for each issue
4. All existing tests pass with updated assertions

## Scenarios

```gherkin
Scenario 1: get_issue includes frontmatter with id
  Given an issue with internal ID 12 and display ID "DGR-42"
  When I call get_issue with display_id "DGR-42"
  Then the output starts with "---\nid: 12\ndisplay_id: DGR-42\nstatus: open\ntype: story"

Scenario 2: get_issue frontmatter includes parent info
  Given the issue has a parent with internal ID 8 and display ID "EPIC-1"
  When I call get_issue with display_id "DGR-42"
  Then the frontmatter includes "parent_id: 8" and "parent_display_id: EPIC-1"

Scenario 3: get_doc includes frontmatter
  Given a doc with internal ID 5 and display ID "DGR-3"
  When I call get_doc with display_id "DGR-3"
  Then the output starts with frontmatter containing "id: 5"

Scenario 4: list_issues includes id field
  Given issues exist in the database
  When I call list_issues
  Then each JSON object includes an "id" field with the internal DB ID

Scenario 5: No parent info in frontmatter when no parent
  Given an issue with no parent
  When I call get_issue
  Then the frontmatter does not contain "parent_id" or "parent_display_id"
```

## Technical notes

- Add YAML frontmatter to `RenderIssueContext` in `internal/issues/markdown.go` and `RenderDoc` in `internal/docs/markdown.go`
- Use `Issue.ID` and `IssueContext.Parent` fields (already available) to populate frontmatter
- Add `ID int` to `IssueSummary` in `internal/issues/models.go`
- Add `i.id` to `ListIssues` SELECT and scan in `internal/issues/issues.go`
- Update test assertions in `markdown_test.go` and `issues_test.go`

## Related docs

- ADR: (none)
- Code exploration: (none)
- Parent epic: DGR-5

---

*This ticket was created by opencode and reviewed by a human before publishing.*
