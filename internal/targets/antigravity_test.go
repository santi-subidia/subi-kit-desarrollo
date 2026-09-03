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
