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

	if !strings.HasPrefix(output, "# DGR-3: Deploy on Railway + Supabase") {
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
	if len(lines) < 3 {
		t.Fatal("expected at least 3 lines")
	}
	if lines[2] != "" && !strings.Contains(lines[2], "**") {
		t.Errorf("expected only header + metadata lines when body is nil, got:\n%s", output)
	}
}
