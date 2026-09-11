package docs

import (
	"strings"
	"testing"
)

func TestRenderDoc_Full(t *testing.T) {
	body := "## Decision\n\nWe will use PostgreSQL."
	doc := &Doc{
		ID:        5,
		DisplayID: "DGR-3",
		DocType:   "adr",
		Title:     "Deploy on Railway + Supabase",
		Body:      &body,
		Status:    "approved",
	}

	output := RenderDoc(doc)

	if !strings.HasPrefix(output, "---\n") {
		t.Errorf("expected frontmatter prefix, but got:\n%s", output)
	}
	if !strings.Contains(output, "id: 5") {
		t.Errorf("expected frontmatter id, got:\n%s", output)
	}
	if !strings.Contains(output, "display_id: DGR-3") {
		t.Errorf("expected frontmatter display_id, got:\n%s", output)
	}
	if !strings.Contains(output, "status: approved") {
		t.Errorf("expected frontmatter status, got:\n%s", output)
	}
	if !strings.Contains(output, "type: adr") {
		t.Errorf("expected frontmatter type, got:\n%s", output)
	}

	if !strings.Contains(output, "# DGR-3: Deploy on Railway + Supabase") {
		t.Errorf("expected header line, got:\n%s", output)
	}
	if !strings.Contains(output, "**Status:** approved") {
		t.Errorf("expected status line, got:\n%s", output)
	}
	if !strings.Contains(output, "**Type:** adr") {
		t.Errorf("expected type line, got:\n%s", output)
	}
	if !strings.Contains(output, "## Decision\n\nWe will use PostgreSQL.") {
		t.Errorf("expected body content, got:\n%s", output)
	}
}

func TestRenderDoc_NilBody(t *testing.T) {
	doc := &Doc{
		ID:        10,
		DisplayID: "CE-1",
		DocType:   "adr",
		Title:     "Code Exploration",
		Body:      nil,
		Status:    "proposed",
	}

	output := RenderDoc(doc)

	if !strings.Contains(output, "# CE-1: Code Exploration") {
		t.Errorf("expected header, got:\n%s", output)
	}
	if !strings.Contains(output, "**Status:** proposed") {
		t.Errorf("expected status, got:\n%s", output)
	}

	// No extra content after the metadata block
	lines := strings.Split(output, "\n")
	if len(lines) < 10 {
		t.Fatal("expected at least 10 lines")
	}
	if lines[9] != "" && !strings.Contains(lines[9], "**") {
		t.Errorf("expected only header + metadata lines when body is nil, got:\n%s", output)
	}
}
