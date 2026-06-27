package docs

import (
	"fmt"
	"strings"
)

func RenderDoc(doc *Doc) string {
	var b strings.Builder

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
