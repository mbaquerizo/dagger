package issues

import (
	"fmt"
	"strings"
)

func RenderIssueContext(ctx *IssueContext) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s: %s\n\n", ctx.Issue.DisplayID, ctx.Issue.Title)
	fmt.Fprintf(&b, "**Status:** %s  |  **Type:** %s\n", ctx.Issue.Status, ctx.Issue.TypeName)

	if ctx.Parent != nil {
		fmt.Fprintf(&b, "**Parent:** %s (%s)\n", ctx.Parent.DisplayID, ctx.Parent.Title)
	}

	if ctx.Issue.Body != nil {
		fmt.Fprintf(&b, "\n%s\n", *ctx.Issue.Body)
	}

	fmt.Fprint(&b, "\n---\n")

	if len(ctx.LinkedDocs) > 0 {
		fmt.Fprint(&b, "\n## Linked Context\n\n")

		for _, doc := range ctx.LinkedDocs {
			fmt.Fprintf(&b, "### %s: %s\n", doc.DisplayID, doc.Title)
			fmt.Fprintf(&b, "**Status:** %s", doc.Status)

			if doc.Body != nil {
				fmt.Fprintf(&b, "\n%s\n", *doc.Body)
			}
		}
	}

	if len(ctx.Children) > 0 {
		fmt.Fprint(&b, "\n## Subtasks\n\n")

		for _, child := range ctx.Children {
			fmt.Fprintf(&b, "- %s: %s\n", child.DisplayID, child.Title)
		}
	}

	if len(ctx.RelatedIssues) > 0 {
		fmt.Fprint(&b, "\n## Related Issues\n\n")

		for _, relatedIssue := range ctx.RelatedIssues {
			fmt.Fprintf(&b, "- **%s** %s: %s\n", relatedIssue.RelationType, relatedIssue.DisplayID, relatedIssue.Title)
		}
	}

	return b.String()
}
