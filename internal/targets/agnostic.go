package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santi-subidia/dev-kit-desarrollo/internal/filemerge"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
)

// AgnosticTarget genera archivos de reglas consolidados en formato Markdown universal (AGENTS.md o .cursorrules).
type AgnosticTarget struct{}

func NewAgnosticTarget() *AgnosticTarget {
	return &AgnosticTarget{}
}

// GenerateAgentsMD genera o inyecta directrices consolidadas en AGENTS.md preservando cualquier contenido existente.
// Utiliza inyección de secciones idempotentes con marcadores HTML y escrituras atómicas en disco.
func (t *AgnosticTarget) GenerateAgentsMD(projectRoot string, selectedRules []*rules.Rule) (string, error) {
	var builder strings.Builder

	// 1. Ubicar orchestrator-role primero si está presente
	for _, r := range selectedRules {
		if r.Metadata.Name == "orchestrator-role" {
			builder.WriteString("## ROL PRINCIPAL: AGENTE ORQUESTADOR & TECH LEAD\n\n")
			builder.WriteString(r.Body)
			builder.WriteString("\n\n---\n\n")
			break
		}
	}

	// 2. Agrupar por categoría
	byCat := make(map[string][]*rules.Rule)
	var catOrder []string
	seenCat := make(map[string]bool)

	for _, r := range selectedRules {
		if r.Metadata.Name == "orchestrator-role" {
			continue
		}
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
	existingBytes, err := os.ReadFile(destPath)

	var finalContent string
	managedSection := strings.TrimSpace(builder.String())

	if err == nil {
		// Archivo existente: inyectar o actualizar bloque gestionado sin borrar contenido previo del proyecto
		finalContent = filemerge.InjectMarkdownSection(string(existingBytes), "managed-rules", managedSection)
	} else if os.IsNotExist(err) {
		// Archivo nuevo: incluir encabezado estándar e inyectar sección con marcadores
		header := "# Directrices y Reglas del Proyecto (Dev-Kit)\n\n> Este archivo consolida las reglas activas de arquitectura, calidad y convenciones del proyecto.\n\n"
		finalContent = filemerge.InjectMarkdownSection(header, "managed-rules", managedSection)
	} else {
		return "", fmt.Errorf("error al leer AGENTS.md: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(destPath, []byte(finalContent), 0644); err != nil {
		return "", fmt.Errorf("error al escribir AGENTS.md atómicamente: %w", err)
	}

	return destPath, nil
}
