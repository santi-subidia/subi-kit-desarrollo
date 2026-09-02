package skills

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager administra la carga, consulta y generación de Skills del dev-kit.
type Manager struct {
	skills []*Skill
	byName map[string]*Skill
	rawFS  fs.FS
}

// NewManager carga los skills desde un filesystem embebido o local.
func NewManager(catalogFS fs.FS) (*Manager, error) {
	m := &Manager{
		skills: make([]*Skill, 0),
		byName: make(map[string]*Skill),
		rawFS:  catalogFS,
	}

	entries, err := fs.ReadDir(catalogFS, ".")
	if err != nil {
		return nil, fmt.Errorf("error al leer directorio de skills: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillDir := entry.Name()
		skillFilePath := filepath.ToSlash(filepath.Join(skillDir, "SKILL.md"))

		data, err := fs.ReadFile(catalogFS, skillFilePath)
		if err != nil {
			continue // No tiene SKILL.md
		}

		skill, err := ParseSkill(data, skillFilePath)
		if err != nil {
			continue
		}

		// Leer plantillas si existen en templates/
		templatesDir := filepath.ToSlash(filepath.Join(skillDir, "templates"))
		if tmplEntries, err := fs.ReadDir(catalogFS, templatesDir); err == nil {
			for _, tEntry := range tmplEntries {
				if !tEntry.IsDir() && strings.HasSuffix(tEntry.Name(), ".md") {
					tPath := filepath.ToSlash(filepath.Join(templatesDir, tEntry.Name()))
					if tData, err := fs.ReadFile(catalogFS, tPath); err == nil {
						skill.Templates[tEntry.Name()] = string(tData)
					}
				}
			}
		}

		m.skills = append(m.skills, skill)
		m.byName[skill.Metadata.Name] = skill
	}

	return m, nil
}

// GetAll retorna todos los skills disponibles.
func (m *Manager) GetAll() []*Skill {
	return m.skills
}

// GetSkill obtiene un skill por su nombre.
func (m *Manager) GetSkill(name string) (*Skill, bool) {
	s, ok := m.byName[name]
	return s, ok
}

// CreateSDDFeature inicializa la carpeta .specs/<feature-name>/ con todas las plantillas SDD.
func (m *Manager) CreateSDDFeature(projectRoot string, featureName string, author string) (string, error) {
	if featureName == "" {
		return "", fmt.Errorf("el nombre de la feature no puede estar vacío")
	}

	// Normalizar nombre de feature (kebab-case)
	normName := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(featureName), " ", "-"))
	specsDir := filepath.Join(projectRoot, ".specs", normName)

	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return "", fmt.Errorf("no se pudo crear la carpeta de la feature %s: %w", specsDir, err)
	}

	sddSkill, ok := m.byName["sdd-workflow"]
	if !ok {
		return "", fmt.Errorf("skill 'sdd-workflow' no encontrado en el catálogo")
	}

	now := time.Now().Format("2006-01-02")
	if author == "" {
		author = os.Getenv("USERNAME")
		if author == "" {
			author = os.Getenv("USER")
		}
		if author == "" {
			author = "Developer"
		}
	}

	replacer := strings.NewReplacer(
		"{{FEATURE_NAME}}", normName,
		"{{DATE}}", now,
		"{{AUTHOR}}", author,
	)

	templateFiles := map[string]string{
		"spec.template.md":      "spec.md",
		"tech-plan.template.md": "tech-plan.md",
		"tasks.template.md":     "tasks.md",
		"verify.template.md":    "verify.md",
		"archive.template.md":   "archive.md",
	}

	for tmplName, targetFile := range templateFiles {
		content, exists := sddSkill.Templates[tmplName]
		if !exists {
			continue
		}

		filledContent := replacer.Replace(content)
		destPath := filepath.Join(specsDir, targetFile)

		// Solo crear si no existe para no sobrescribir trabajo previo
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			if err := os.WriteFile(destPath, []byte(filledContent), 0644); err != nil {
				return "", fmt.Errorf("error al escribir %s: %w", destPath, err)
			}
		}
	}

	return specsDir, nil
}
