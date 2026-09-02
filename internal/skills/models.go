package skills

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Metadata contiene la cabecera YAML de un Skill.
type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Skill representa una habilidad completa del dev-kit.
type Skill struct {
	Metadata   Metadata
	Body       string            // Contenido Markdown sin el frontmatter
	RawContent string            // Contenido original completo con frontmatter
	RelPath    string            // Ruta relativa (ej: sdd-workflow/SKILL.md)
	Templates  map[string]string // Plantillas asociadas al skill (ej: spec.template.md)
}

// ParseSkill parsea el contenido de un archivo SKILL.md.
func ParseSkill(raw []byte, relPath string) (*Skill, error) {
	rawStr := string(raw)
	parts := strings.SplitN(rawStr, "---", 3)

	if len(parts) < 3 {
		return nil, fmt.Errorf("formato inválido en %s: se esperaba cabecera YAML '---'", relPath)
	}

	yamlHeader := parts[1]
	body := strings.TrimSpace(parts[2])

	var meta Metadata
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(yamlHeader)))
	if err := decoder.Decode(&meta); err != nil {
		return nil, fmt.Errorf("error al decodificar YAML frontmatter en %s: %w", relPath, err)
	}

	if meta.Name == "" {
		meta.Name = "unnamed-skill"
	}

	return &Skill{
		Metadata:   meta,
		Body:       body,
		RawContent: rawStr,
		RelPath:    relPath,
		Templates:  make(map[string]string),
	}, nil
}
