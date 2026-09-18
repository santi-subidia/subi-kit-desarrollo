package targets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/santi-subidia/dev-kit-desarrollo/internal/filemerge"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
)

func TestGenerateAgentsMD_NewFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-agnostic-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAgnosticTarget()
	sampleRules := []*rules.Rule{
		{
			Metadata: rules.Metadata{
				Name:     "clean-code",
				Title:    "Clean Code Guidelines",
				Category: "core",
			},
			Body: "Do not write complex unreadable code.",
		},
	}

	dest, err := target.GenerateAgentsMD(tempDir, sampleRules)
	if err != nil {
		t.Fatalf("GenerateAgentsMD failed: %v", err)
	}

	contentBytes, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}

	content := string(contentBytes)
	if !strings.Contains(content, filemerge.OpenMarker("managed-rules")) {
		t.Errorf("expected open marker for managed-rules")
	}
	if !strings.Contains(content, filemerge.CloseMarker("managed-rules")) {
		t.Errorf("expected close marker for managed-rules")
	}
	if !strings.Contains(content, "Do not write complex unreadable code.") {
		t.Errorf("expected rule content present")
	}
}

func TestGenerateAgentsMD_PreservesUserContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-agnostic-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	userContent := `# Custom Repo Instructions

| Skill | Trigger |
| --- | --- |
| custom-skill | when triggering |

Keep this table safe!`

	if err := os.WriteFile(agentsPath, []byte(userContent), 0644); err != nil {
		t.Fatalf("failed to seed AGENTS.md: %v", err)
	}

	target := NewAgnosticTarget()
	sampleRules := []*rules.Rule{
		{
			Metadata: rules.Metadata{
				Name:     "security",
				Title:    "Security Protocol",
				Category: "safety",
			},
			Body: "Sanitize all user inputs.",
		},
	}

	_, err = target.GenerateAgentsMD(tempDir, sampleRules)
	if err != nil {
		t.Fatalf("GenerateAgentsMD failed: %v", err)
	}

	contentBytes, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read updated AGENTS.md: %v", err)
	}

	content := string(contentBytes)
	// User content preserved
	if !strings.Contains(content, "# Custom Repo Instructions") {
		t.Errorf("expected user header preserved")
	}
	if !strings.Contains(content, "Keep this table safe!") {
		t.Errorf("expected user table preserved")
	}
	// Managed section injected
	if !strings.Contains(content, "Sanitize all user inputs.") {
		t.Errorf("expected managed rules injected")
	}
	if !strings.Contains(content, filemerge.OpenMarker("managed-rules")) {
		t.Errorf("expected managed-rules markers")
	}
}
