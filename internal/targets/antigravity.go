package targets

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/santi-subidia/dev-kit-desarrollo/internal/agents"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/skills"
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

	return writtenFiles, nil
}

// InstallGlobal instala rules, skills y agentes en la configuración global (~/.gemini/config/).
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

	return writtenFiles, nil
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
