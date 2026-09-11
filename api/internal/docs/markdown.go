package docs

import (
	"fmt"
	"strings"
)

func RenderDoc(doc *Doc) string {
	var b strings.Builder

	fmt.Fprintf(&b, "---\n")
	fmt.Fprintf(&b, "id: %d\n", doc.ID)
	fmt.Fprintf(&b, "display_id: %s\n", doc.DisplayID)
	fmt.Fprintf(&b, "status: %s\n", doc.Status)
	fmt.Fprintf(&b, "type: %s\n", doc.DocType)

	if doc.Parent != nil {
		fmt.Fprintf(&b, "parent_id: %d\n", doc.Parent.ID)
		fmt.Fprintf(&b, "parent_display_id: %s\n", doc.Parent.DisplayID)
	}
	fmt.Fprintf(&b, "---\n\n")

	fmt.Fprintf(&b, "# %s: %s\n\n", doc.DisplayID, doc.Title)
	fmt.Fprintf(&b, "**Status:** %s  |  **Type:** %s\n", doc.Status, doc.DocType)

	if doc.Parent != nil {
		fmt.Fprintf(&b, "**Parent:** %s (%s)\n", doc.Parent.DisplayID, doc.Parent.Title)
	}

	if doc.Body != nil {
		fmt.Fprintf(&b, "\n%s\n", *doc.Body)
	}

	return b.String()
}
