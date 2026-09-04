package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santi-subidia/dev-kit-desarrollo/internal/agents"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/skills"
)

const (
	globalGeminiStartTag = "<!-- BEGIN SUBIKIT ORCHESTRATOR ROLE -->"
	globalGeminiEndTag   = "<!-- END SUBIKIT ORCHESTRATOR ROLE -->"
)

// AntigravityTarget maneja la instalación de reglas, skills y agentes para Antigravity y Gemini CLI.
type AntigravityTarget struct{}

func NewAntigravityTarget() *AntigravityTarget {
	return &AntigravityTarget{}
}

// InstallProject instala rules, skills y agentes en el directorio .agents/ del proyecto.
func (t *AntigravityTarget) InstallProject(projectRoot string, selectedRules []*rules.Rule, selectedSkills []*skills.Skill, selectedAgents []*agents.Agent, force bool) ([]string, error) {
	var writtenFiles []string

	// 1. Instalar Rules en .agents/rules/
	rulesDir := filepath.Join(projectRoot, ".agents", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear el directorio de reglas %s: %w", rulesDir, err)
	}

	for _, r := range selectedRules {
		destPath := filepath.Join(rulesDir, r.Metadata.Name+".md")
		if err := os.WriteFile(destPath, []byte(r.RawContent), 0644); err != nil {
			return writtenFiles, fmt.Errorf("error al escribir regla %s: %w", destPath, err)
		}
		writtenFiles = append(writtenFiles, destPath)
	}

	// 2. Instalar Skills en .agents/skills/<skillName>/SKILL.md
	for _, s := range selectedSkills {
		skillDir := filepath.Join(projectRoot, ".agents", "skills", s.Metadata.Name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return writtenFiles, fmt.Errorf("no se pudo crear directorio del skill %s: %w", skillDir, err)
		}

		destPath := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(destPath, []byte(s.RawContent), 0644); err != nil {
			return writtenFiles, fmt.Errorf("error al escribir skill %s: %w", destPath, err)
		}
		writtenFiles = append(writtenFiles, destPath)
	}

	// 3. Instalar Agentes y Subagentes en .agents/agents/
	agentsDir := filepath.Join(projectRoot, ".agents", "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return writtenFiles, fmt.Errorf("no se pudo crear directorio de agentes %s: %w", agentsDir, err)
	}

	for _, a := range selectedAgents {
		destPath := filepath.Join(agentsDir, a.Metadata.Name+".md")
		if err := os.WriteFile(destPath, []byte(a.RawContent), 0644); err != nil {
			return writtenFiles, fmt.Errorf("error al escribir agente %s: %w", destPath, err)
		}
		writtenFiles = append(writtenFiles, destPath)
	}

	// 4. Asegurar que .codegraph/ esté en .gitignore
	if err := t.EnsureGitignore(projectRoot); err == nil {
		writtenFiles = append(writtenFiles, filepath.Join(projectRoot, ".gitignore"))
	}

	return writtenFiles, nil
}

// EnsureGitignore verifica y añade .codegraph/ al archivo .gitignore del proyecto si no está presente.
func (t *AntigravityTarget) EnsureGitignore(projectRoot string) error {
	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	var existingContent string

	data, err := os.ReadFile(gitignorePath)
	if err == nil {
		existingContent = string(data)
		if strings.Contains(existingContent, ".codegraph/") || strings.Contains(existingContent, ".codegraph") {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	var toAppend string
	if existingContent == "" {
		toAppend = "# Dev-Kit y MCPs\n.codegraph/\n"
	} else {
		toAppend = strings.TrimRight(existingContent, "\r\n") + "\n\n# Dev-Kit y MCPs\n.codegraph/\n"
	}

	return os.WriteFile(gitignorePath, []byte(toAppend), 0644)
}

// GenerateGeminiMD genera o sincroniza el archivo GEMINI.md consolidado con directrices y rol de orquestador.
func (t *AntigravityTarget) GenerateGeminiMD(projectRoot string, selectedRules []*rules.Rule) (string, error) {
	var builder strings.Builder

	builder.WriteString("# Directrices del Proyecto para Antigravity & Asistentes de IA\n\n")
	builder.WriteString("> Este archivo define las directrices maestras del proyecto y el rol permanente de Agente Orquestador.\n\n")

	// 1. Ubicar orchestrator-role primero si está presente
	for _, r := range selectedRules {
		if r.Metadata.Name == "orchestrator-role" {
			builder.WriteString(r.Body)
			builder.WriteString("\n\n---\n\n")
			break
		}
	}

	// 2. Agrupar el resto por categoría
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

	destPath := filepath.Join(projectRoot, "GEMINI.md")
	if err := os.WriteFile(destPath, []byte(builder.String()), 0644); err != nil {
		return "", fmt.Errorf("error al escribir GEMINI.md: %w", err)
	}

	return destPath, nil
}

// InstallGlobal instala rules, skills y agentes en la configuración global (~/.gemini/config/).
// Además sincroniza de forma segura el rol de Orquestador en ~/.gemini/GEMINI.md.
func (t *AntigravityTarget) InstallGlobal(selectedRules []*rules.Rule, selectedSkills []*skills.Skill, selectedAgents []*agents.Agent, force bool) ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el directorio home del usuario: %w", err)
	}

	var writtenFiles []string

	// Reglas globales
	rulesDir := filepath.Join(homeDir, ".gemini", "config", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear directorio global de reglas %s: %w", rulesDir, err)
	}

	for _, r := range selectedRules {
		destPath := filepath.Join(rulesDir, r.Metadata.Name+".md")
		if err := os.WriteFile(destPath, []byte(r.RawContent), 0644); err != nil {
			return writtenFiles, fmt.Errorf("error al escribir regla global %s: %w", destPath, err)
		}
		writtenFiles = append(writtenFiles, destPath)
	}

	// Skills globales
	for _, s := range selectedSkills {
		skillDir := filepath.Join(homeDir, ".gemini", "config", "skills", s.Metadata.Name)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return writtenFiles, fmt.Errorf("no se pudo crear directorio global del skill %s: %w", skillDir, err)
		}

		destPath := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(destPath, []byte(s.RawContent), 0644); err != nil {
			return writtenFiles, fmt.Errorf("error al escribir skill global %s: %w", destPath, err)
		}
		writtenFiles = append(writtenFiles, destPath)
	}

	// Agentes globales
	agentsDir := filepath.Join(homeDir, ".gemini", "config", "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return writtenFiles, fmt.Errorf("no se pudo crear directorio global de agentes %s: %w", agentsDir, err)
	}

	for _, a := range selectedAgents {
		destPath := filepath.Join(agentsDir, a.Metadata.Name+".md")
		if err := os.WriteFile(destPath, []byte(a.RawContent), 0644); err != nil {
			return writtenFiles, fmt.Errorf("error al escribir agente global %s: %w", destPath, err)
		}
		writtenFiles = append(writtenFiles, destPath)
	}

	// Sincronizar GEMINI.md global para rol de Orquestador
	var orchRule *rules.Rule
	for _, r := range selectedRules {
		if r.Metadata.Name == "orchestrator-role" {
			orchRule = r
			break
		}
	}
	if orchRule != nil {
		globalGemini, err := t.syncGlobalGeminiMD(homeDir, orchRule)
		if err == nil && globalGemini != "" {
			writtenFiles = append(writtenFiles, globalGemini)
		}
	}

	return writtenFiles, nil
}

func (t *AntigravityTarget) syncGlobalGeminiMD(homeDir string, orchestratorRule *rules.Rule) (string, error) {
	if orchestratorRule == nil {
		return "", nil
	}

	geminiPath := filepath.Join(homeDir, ".gemini", "GEMINI.md")
	block := fmt.Sprintf("%s\n## %s\n\n%s\n%s\n",
		globalGeminiStartTag,
		orchestratorRule.Metadata.Title,
		orchestratorRule.Body,
		globalGeminiEndTag,
	)

	existingBytes, err := os.ReadFile(geminiPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(geminiPath), 0755); err != nil {
				return "", err
			}
			if err := os.WriteFile(geminiPath, []byte(block), 0644); err != nil {
				return "", err
			}
			return geminiPath, nil
		}
		return "", err
	}

	content := string(existingBytes)
	startIdx := strings.Index(content, globalGeminiStartTag)
	endIdx := strings.Index(content, globalGeminiEndTag)

	var newContent string
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Reemplazar bloque existente manteniendo el resto intacto
		endPos := endIdx + len(globalGeminiEndTag)
		newContent = content[:startIdx] + strings.TrimSpace(block) + content[endPos:]
	} else {
		// Anexar al final sin tocar los bloques previos (CodeGraph, Context7, Engram)
		newContent = strings.TrimSpace(content) + "\n\n" + block
	}

	if err := os.WriteFile(geminiPath, []byte(newContent), 0644); err != nil {
		return "", err
	}

	return geminiPath, nil
}

// GetProjectRulesPath retorna la ruta de reglas locales.
func (t *AntigravityTarget) GetProjectRulesPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".agents", "rules")
}

// GetProjectSkillsPath retorna la ruta de skills locales.
func (t *AntigravityTarget) GetProjectSkillsPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".agents", "skills")
}

// GetProjectAgentsPath retorna la ruta de agentes locales.
func (t *AntigravityTarget) GetProjectAgentsPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".agents", "agents")
}

// GetGlobalRulesPath retorna la ruta de configuración global de Antigravity.
func (t *AntigravityTarget) GetGlobalRulesPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".gemini", "config", "rules"), nil
}
