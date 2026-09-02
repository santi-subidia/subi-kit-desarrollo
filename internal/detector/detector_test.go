package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMonorepoProject(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "devkit-test-monorepo-*")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Crear frontend/
	frontDir := filepath.Join(tempDir, "frontend")
	os.MkdirAll(frontDir, 0755)
	pkgJSON := `{
		"name": "frontend",
		"dependencies": {
			"next": "16.0.0",
			"react": "19.0.0",
			"tailwindcss": "4.0.0"
		},
		"devDependencies": {
			"typescript": "^5.0.0",
			"jest": "^29.0.0"
		}
	}`
	os.WriteFile(filepath.Join(frontDir, "package.json"), []byte(pkgJSON), 0644)
	os.WriteFile(filepath.Join(frontDir, "tsconfig.json"), []byte("{}"), 0644)

	// Crear backend/ con .slnx
	backDir := filepath.Join(tempDir, "backend")
	os.MkdirAll(backDir, 0755)
	os.WriteFile(filepath.Join(backDir, "Barberia.slnx"), []byte("<Solution />"), 0644)

	// Crear docker-compose.yml en la raíz
	os.WriteFile(filepath.Join(tempDir, "docker-compose.yml"), []byte("version: '3'"), 0644)

	stack, err := Detect(tempDir)
	if err != nil {
		t.Fatalf("Detect error: %v", err)
	}

	expectedTechs := []string{".NET / C#", "Docker", "Next.js", "React", "Tailwind CSS", "Testing (Jest/Vitest)", "TypeScript"}
	for _, tech := range expectedTechs {
		found := false
		for _, t := range stack.Technologies {
			if t == tech {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected '%s' to be detected in monorepo, detected: %v", tech, stack.Technologies)
		}
	}
}
