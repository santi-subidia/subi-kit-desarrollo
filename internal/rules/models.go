package rules

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Metadata contiene la cabecera YAML de una regla.
type Metadata struct {
	Name        string   `yaml:"name"`
	Title       string   `yaml:"title"`
	Category    string   `yaml:"category"`
	AlwaysOn    bool     `yaml:"always_on"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
}

// Rule representa una regla completa del dev-kit.
type Rule struct {
	Metadata   Metadata
	Body       string // Contenido Markdown sin el frontmatter
	RawContent string // Contenido original completo con frontmatter
	RelPath    string // Ruta relativa en el catálogo (ej: core/git-conventions.md)
}

// ParseRule parsea el contenido de un archivo Markdown con frontmatter YAML.
func ParseRule(raw []byte, relPath string) (*Rule, error) {
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
		meta.Name = strings.TrimSuffix(strings.ReplaceAll(relPath, "/", "-"), ".md")
	}

	return &Rule{
		Metadata:   meta,
		Body:       body,
		RawContent: rawStr,
		RelPath:    relPath,
	}, nil
}
