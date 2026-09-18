package filemerge

import (
	"strings"
	"testing"
)

func TestInjectMarkdownSection_Append(t *testing.T) {
	initial := "# My Project\n\nSome introductory text.\n"
	content := "### Rule 1\nFollow clean code."

	res := InjectMarkdownSection(initial, "rules", content)

	if !strings.Contains(res, "# My Project") {
		t.Errorf("expected initial content preserved")
	}
	if !strings.Contains(res, OpenMarker("rules")) {
		t.Errorf("expected open marker")
	}
	if !strings.Contains(res, CloseMarker("rules")) {
		t.Errorf("expected close marker")
	}
	if !strings.Contains(res, "Follow clean code.") {
		t.Errorf("expected rule content injected")
	}
}

func TestInjectMarkdownSection_Replace(t *testing.T) {
	initial := "# Header\n\n<!-- subikit:rules -->\nOld content\n<!-- /subikit:rules -->\n\n## Footer\n"
	newContent := "Brand new content"

	res := InjectMarkdownSection(initial, "rules", newContent)

	if strings.Contains(res, "Old content") {
		t.Errorf("expected old content replaced, got %s", res)
	}
	if !strings.Contains(res, "Brand new content") {
		t.Errorf("expected new content present, got %s", res)
	}
	if !strings.Contains(res, "# Header") || !strings.Contains(res, "## Footer") {
		t.Errorf("expected Header and Footer preserved, got %s", res)
	}
}

func TestInjectMarkdownSection_Remove(t *testing.T) {
	initial := "# Header\n\n<!-- subikit:rules -->\nOld content to remove\n<!-- /subikit:rules -->\n\n## Footer\n"

	res := InjectMarkdownSection(initial, "rules", "")

	if strings.Contains(res, "Old content to remove") {
		t.Errorf("expected content removed")
	}
	if strings.Contains(res, OpenMarker("rules")) || strings.Contains(res, CloseMarker("rules")) {
		t.Errorf("expected markers removed")
	}
	if !strings.Contains(res, "# Header") || !strings.Contains(res, "## Footer") {
		t.Errorf("expected Header and Footer preserved")
	}
}

func TestInjectMarkdownSection_DuplicateBlocks(t *testing.T) {
	initial := `# Header
<!-- subikit:rules -->
Block 1
<!-- /subikit:rules -->

Middle text

<!-- subikit:rules -->
Block 2 (duplicate)
<!-- /subikit:rules -->
# Tail`

	res := InjectMarkdownSection(initial, "rules", "Single active block")

	// Count occurrences of open marker
	count := strings.Count(res, OpenMarker("rules"))
	if count != 1 {
		t.Errorf("expected exactly 1 open marker, got %d. Content:\n%s", count, res)
	}
	if strings.Contains(res, "Block 1") || strings.Contains(res, "Block 2") {
		t.Errorf("expected old blocks purged")
	}
	if !strings.Contains(res, "Single active block") {
		t.Errorf("expected new block active")
	}
	if !strings.Contains(res, "Middle text") || !strings.Contains(res, "# Tail") {
		t.Errorf("expected Middle text and Tail preserved")
	}
}

func TestInjectMarkdownSection_OrphanMarkers(t *testing.T) {
	// Orphan close marker followed by valid pair
	corrupted := `<!-- /subikit:rules -->
# Title
<!-- subikit:rules -->
Valid
<!-- /subikit:rules -->`

	res := InjectMarkdownSection(corrupted, "rules", "Fresh Content")

	if strings.Count(res, CloseMarker("rules")) != 1 {
		t.Errorf("expected orphan close marker stripped, got:\n%s", res)
	}
	if !strings.Contains(res, "Fresh Content") {
		t.Errorf("expected fresh content injected")
	}
}

func TestInjectSectionWithMarkers_Custom(t *testing.T) {
	open := "<!-- BEGIN CUSTOM -->"
	close := "<!-- END CUSTOM -->"
	existing := "Pre-text\n" + open + "\nOld\n" + close + "\nPost-text"

	res := InjectSectionWithMarkers(existing, open, close, "Updated Custom")

	if !strings.Contains(res, "Pre-text") || !strings.Contains(res, "Post-text") {
		t.Errorf("surrounding text must be preserved")
	}
	if !strings.Contains(res, "Updated Custom") {
		t.Errorf("new content must be present")
	}
	if strings.Contains(res, "Old") {
		t.Errorf("old content must be replaced")
	}

	extracted := ExtractSectionWithMarkers(res, open, close)
	if extracted != "Updated Custom" {
		t.Errorf("expected extracted %q, got %q", "Updated Custom", extracted)
	}
}
