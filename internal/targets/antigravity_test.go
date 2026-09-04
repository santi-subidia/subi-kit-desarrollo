package targets

import (
	"os"
	"strings"
	"testing"

	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
)

func TestGenerateGeminiMD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-antigravity-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAntigravityTarget()
	sampleRules := []*rules.Rule{
		{
			Metadata: rules.Metadata{
				Name:     "orchestrator-role",
				Title:    "Rol de Agente Orquestador & Protocolo de Delegación a Subagentes",
				Category: "core",
			},
			Body: "Eres el Agente Orquestador.",
		},
		{
			Metadata: rules.Metadata{
				Name:     "clean-code",
				Title:    "Clean Code y Simplicidad",
				Category: "core",
			},
			Body: "Reglas de clean code.",
		},
	}

	geminiPath, err := target.GenerateGeminiMD(tempDir, sampleRules)
	if err != nil {
		t.Fatalf("GenerateGeminiMD error: %v", err)
	}

	content, err := os.ReadFile(geminiPath)
	if err != nil {
		t.Fatalf("failed to read generated GEMINI.md: %v", err)
	}

	strContent := string(content)
	if !strings.Contains(strContent, "Eres el Agente Orquestador.") {
		t.Errorf("expected GEMINI.md to contain orchestrator body, got: %s", strContent)
	}
	if !strings.Contains(strContent, "Módulo: CORE") {
		t.Errorf("expected GEMINI.md to contain category headers, got: %s", strContent)
	}
}

func TestSyncGlobalGeminiMD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-global-gemini-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAntigravityTarget()
	orchRule := &rules.Rule{
		Metadata: rules.Metadata{
			Name:  "orchestrator-role",
			Title: "Rol de Agente Orquestador",
		},
		Body: "Orchestrator instructions.",
	}

	// 1. Initial creation
	geminiPath, err := target.syncGlobalGeminiMD(tempDir, orchRule)
	if err != nil {
		t.Fatalf("syncGlobalGeminiMD failed on initial write: %v", err)
	}

	data, err := os.ReadFile(geminiPath)
	if err != nil {
		t.Fatalf("failed to read global GEMINI.md: %v", err)
	}
	if !strings.Contains(string(data), globalGeminiStartTag) {
		t.Errorf("expected start tag in %s", string(data))
	}
	if !strings.Contains(string(data), "Orchestrator instructions.") {
		t.Errorf("expected body in %s", string(data))
	}

	// 2. Idempotent update
	orchRule.Body = "Updated orchestrator instructions."
	geminiPath, err = target.syncGlobalGeminiMD(tempDir, orchRule)
	if err != nil {
		t.Fatalf("syncGlobalGeminiMD failed on update: %v", err)
	}

	data, err = os.ReadFile(geminiPath)
	if err != nil {
		t.Fatalf("failed to read updated GEMINI.md: %v", err)
	}
	if !strings.Contains(string(data), "Updated orchestrator instructions.") {
		t.Errorf("expected updated body, got: %s", string(data))
	}
	if strings.Count(string(data), globalGeminiStartTag) != 1 {
		t.Errorf("expected exactly 1 tag block, got count %d", strings.Count(string(data), globalGeminiStartTag))
	}
}

func TestEnsureGitignore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-gitignore-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAntigravityTarget()

	// 1. Create on non-existent .gitignore
	if err := target.EnsureGitignore(tempDir); err != nil {
		t.Fatalf("EnsureGitignore failed: %v", err)
	}

	content, err := os.ReadFile(tempDir + "/.gitignore")
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}
	if !strings.Contains(string(content), ".codegraph/") {
		t.Errorf("expected .codegraph/ in .gitignore, got: %s", string(content))
	}

	// 2. Idempotent check on existing .gitignore
	if err := target.EnsureGitignore(tempDir); err != nil {
		t.Fatalf("EnsureGitignore failed on second run: %v", err)
	}

	content2, _ := os.ReadFile(tempDir + "/.gitignore")
	if strings.Count(string(content2), ".codegraph/") != 1 {
		t.Errorf("expected exactly 1 .codegraph/ occurrence, got: %s", string(content2))
	}
}

func TestCleanupLegacyGeminiMD(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-cleanup-gemini-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAntigravityTarget()

	// 1. When GEMINI.md doesn't exist, should return false, nil
	cleaned, err := target.CleanupLegacyGeminiMD(tempDir)
	if err != nil {
		t.Fatalf("unexpected error when file does not exist: %v", err)
	}
	if cleaned {
		t.Errorf("expected cleaned to be false, got true")
	}

	// 2. When GEMINI.md exists, should remove it and return true, nil
	geminiPath := tempDir + "/GEMINI.md"
	if err := os.WriteFile(geminiPath, []byte("# Legacy rules"), 0644); err != nil {
		t.Fatalf("failed to write test GEMINI.md: %v", err)
	}

	cleaned, err = target.CleanupLegacyGeminiMD(tempDir)
	if err != nil {
		t.Fatalf("unexpected error when removing file: %v", err)
	}
	if !cleaned {
		t.Errorf("expected cleaned to be true, got false")
	}

	// Verify file is gone
	if _, err := os.Stat(geminiPath); !os.IsNotExist(err) {
		t.Errorf("expected GEMINI.md to be deleted, but still exists")
	}

	// 3. Second run should return false, nil (idempotent)
	cleaned, err = target.CleanupLegacyGeminiMD(tempDir)
	if err != nil || cleaned {
		t.Errorf("expected cleaned to be false on second run, got cleaned=%v, err=%v", cleaned, err)
	}
}

func TestCleanupLegacySubagents(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-cleanup-subagents-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAntigravityTarget()

	// 1. When .agents/subagents does not exist
	cleaned, err := target.CleanupLegacySubagents(tempDir)
	if err != nil || cleaned {
		t.Fatalf("expected cleaned=false, got %v, err=%v", cleaned, err)
	}

	// 2. When .agents/subagents exists
	legacyDir := tempDir + "/.agents/subagents"
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatalf("failed to create legacy dir: %v", err)
	}
	_ = os.WriteFile(legacyDir+"/architect.md", []byte("legacy"), 0644)

	cleaned, err = target.CleanupLegacySubagents(tempDir)
	if err != nil || !cleaned {
		t.Fatalf("expected cleaned=true, got %v, err=%v", cleaned, err)
	}

	if _, err := os.Stat(legacyDir); !os.IsNotExist(err) {
		t.Errorf("expected legacy dir to be removed")
	}
}

func TestRemoveLocalRedundantAgents(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "subikit-remove-redundant-agents-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	target := NewAntigravityTarget()

	// 1. When directory does not exist
	removed, err := target.RemoveLocalRedundantAgents(tempDir)
	if err != nil || removed != 0 {
		t.Fatalf("expected removed=0, got %d, err=%v", removed, err)
	}

	// 2. When agents exist
	agentsDir := tempDir + "/.agents/agents"
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("failed to create agents dir: %v", err)
	}
	_ = os.WriteFile(agentsDir+"/architect.md", []byte("arch"), 0644)
	_ = os.WriteFile(agentsDir+"/fullstack.md", []byte("fullstack"), 0644)

	removed, err = target.RemoveLocalRedundantAgents(tempDir)
	if err != nil || removed != 2 {
		t.Fatalf("expected removed=2, got %d, err=%v", removed, err)
	}

	// Check that directory is gone or empty
	entries, _ := os.ReadDir(agentsDir)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries in agents dir, got %d", len(entries))
	}
}



