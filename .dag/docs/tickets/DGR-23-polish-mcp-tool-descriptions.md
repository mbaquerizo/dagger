---
id: DGR-23
issueType: task
status: in-review
tags:
  - mcp
blockedBy: []
parent: DGR-5
---

# Polish MCP tool descriptions

## Description

Some MCP tool descriptions in `ListTools()` are terse or ambiguous. Better descriptions help AI agents choose the right tool and use correct parameters.

Update the `Description` fields in `internal/mcp/mcp.go` to accurately reflect the tool's behavior after all other tickets in this epic are implemented.

## Acceptance criteria

1. `tools/list` returns the updated descriptions for all tools
2. Each description accurately reflects the tool's behavior
3. The `status` parameter on `list_issues` clarifies that omission returns all issues

## Updated descriptions

| Tool | Description |
|------|-------------|
| `get_issue` | `"Fetch an issue by display ID with full context (linked docs, parent, children, related issues) as markdown. Includes YAML frontmatter with internal id, display_id, status, type, and parent info."` |
| `get_doc` | `"Fetch a document by display ID as markdown. Includes YAML frontmatter with internal id, display_id, status, and type."` |
| `list_issues` | `"List all issues in the project as JSON. Optionally filter by status (open, in-progress, in-review, done, closed). When status is omitted, returns all issues. Each issue includes id, displayId, title, status, type, and parentDisplayId."` |
| `update_issue_status` | `"Update an issue's status by display ID. Valid status values: open, in-progress, in-review, done, closed."` |
| `publish` | `"Create a new document or issue. Accepts type (issue, adr, pitch, ce), title, body, project_id, optional parent_id, and nested metadata (issue_type, status, tags, relationships, issue_relations). Returns the created id, displayId, and url."` |

Also update the `status` param description on `list_issues`:
- Before: `"Filter by status (open, in-progress, in-review, done, closed)"`
- After: `"Filter by status (open, in-progress, in-review, done, closed). Omit to return all issues."`

## Related docs

- ADR: (none)
- Code exploration: (none)
- Parent epic: DGR-5

---

*This ticket was created by opencode and reviewed by a human before publishing.*
