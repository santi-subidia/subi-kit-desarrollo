package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
)

// AgnosticTarget genera archivos de reglas consolidados en formato Markdown universal (AGENTS.md o .cursorrules).
type AgnosticTarget struct{}

func NewAgnosticTarget() *AgnosticTarget {
	return &AgnosticTarget{}
}

// GenerateAgentsMD genera un archivo AGENTS.md consolidado con todas las reglas seleccionadas.
func (t *AgnosticTarget) GenerateAgentsMD(projectRoot string, selectedRules []*rules.Rule) (string, error) {
	var builder strings.Builder

	builder.WriteString("# Directrices y Reglas del Proyecto (Dev-Kit)\n\n")
	builder.WriteString("> Este archivo consolida las reglas activas de arquitectura, calidad y convenciones del proyecto.\n\n")

	// Agrupar por categoría
	byCat := make(map[string][]*rules.Rule)
	var catOrder []string
	seenCat := make(map[string]bool)

	for _, r := range selectedRules {
		cat := r.Metadata.Category
		if cat == "" {
			cat = "general"
		}
		if !seenCat[cat] {
			seenCat[cat] = true
			catOrder = append(catOrder, cat)
		}
		byCat[cat] = append(byCat[cat], r)
	}

	for _, cat := range catOrder {
		builder.WriteString(fmt.Sprintf("## Módulo: %s\n\n", strings.ToUpper(cat)))
		for _, r := range byCat[cat] {
			builder.WriteString(fmt.Sprintf("### %s\n", r.Metadata.Title))
			if r.Metadata.Description != "" {
				builder.WriteString(fmt.Sprintf("_%s_\n\n", r.Metadata.Description))
			}
			builder.WriteString(r.Body)
			builder.WriteString("\n\n---\n\n")
		}
	}

	destPath := filepath.Join(projectRoot, "AGENTS.md")
	if err := os.WriteFile(destPath, []byte(builder.String()), 0644); err != nil {
		return "", fmt.Errorf("error al escribir AGENTS.md: %w", err)
	}

	return destPath, nil
}
