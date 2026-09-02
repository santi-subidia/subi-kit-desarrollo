package agents

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Metadata contiene la cabecera YAML de un Agente o Subagente.
type Metadata struct {
	Name        string   `yaml:"name"`
	Title       string   `yaml:"title"`
	Type        string   `yaml:"type"` // "orchestrator" o "subagent"
	Description string   `yaml:"description"`
	Tools       []string `yaml:"tools"`
	Subagents   []string `yaml:"subagents"`
}

// Agent representa un agente o subagente completo del dev-kit.
type Agent struct {
	Metadata   Metadata
	Body       string // Contenido Markdown sin el frontmatter
	RawContent string // Contenido original completo con frontmatter
	RelPath    string // Ruta relativa (ej: orchestrator.md o subagents/architect.md)
}

// ParseAgent parsea el contenido de un archivo de agente en Markdown con YAML frontmatter.
func ParseAgent(raw []byte, relPath string) (*Agent, error) {
	rawStr := string(raw)
	parts := strings.SplitN(rawStr, "---", 3)

	if len(parts) < 3 {
		return nil, fmt.Errorf("formato inválido en %s: se esperaba bloque YAML delimitado por '---'", relPath)
	}

	yamlHeader := parts[1]
	body := strings.TrimSpace(parts[2])

	var meta Metadata
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(yamlHeader)))
	if err := decoder.Decode(&meta); err != nil {
		return nil, fmt.Errorf("error al decodificar YAML frontmatter en %s: %w", relPath, err)
	}

	if meta.Name == "" {
		meta.Name = strings.TrimSuffix(relPath, ".md")
	}

	return &Agent{
		Metadata:   meta,
		Body:       body,
		RawContent: rawStr,
		RelPath:    relPath,
	}, nil
}
