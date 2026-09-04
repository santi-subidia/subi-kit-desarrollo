package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/agents"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/detector"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/mcp"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/skills"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/targets"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/updater"
)

// RenderCommandsView renderiza el panel izquierdo (lista) y derecho (detalles) de comandos.
func RenderCommandsView(commands []CommandDoc, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	var listBuilder strings.Builder

	for i, cmd := range commands {
		isSelected := i == selectedIdx
		badge := BadgeCommandStyle.Render(cmd.Category)

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-10s %s\n  %s", cmd.Name, badge, ItemDescStyle.Render(truncate(cmd.Description, leftWidth-14))),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-10s %s\n  %s", cmd.Name, badge, ItemDescStyle.Render(truncate(cmd.Description, leftWidth-14))),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(commands) {
		return listBuilder.String(), "Selecciona un comando para ver los detalles."
	}

	cmd := commands[selectedIdx]
	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("⚡ COMANDO: subikit "+cmd.Name) + "  " + BadgeCommandStyle.Render(cmd.Category) + "\n\n")
	detailBuilder.WriteString(DetailSectionStyle.Render("📖 SINTAXIS") + "\n")
	detailBuilder.WriteString(CodeBlockStyle.Width(detailWidth-4).Render("  "+cmd.Syntax) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("📝 DESCRIPCIÓN") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(cmd.Description) + "\n\n")

	if len(cmd.Flags) > 0 {
		detailBuilder.WriteString(DetailSectionStyle.Render("⚙️ FLAGS Y OPCIONES") + "\n")
		for _, f := range cmd.Flags {
			flagHeader := fmt.Sprintf("  %s %s", KeyBadgeStyle.Render("•"), lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(f.Flag))
			flagDesc := contentStyle.Foreground(ColorDim).Render("    └─ " + f.Description)
			detailBuilder.WriteString(flagHeader + "\n" + flagDesc + "\n")
		}
		detailBuilder.WriteString("\n")
	}

	if len(cmd.Examples) > 0 {
		detailBuilder.WriteString(DetailSectionStyle.Render("💡 EJEMPLOS DE USO") + "\n")
		for _, ex := range cmd.Examples {
			detailBuilder.WriteString(CodeBlockStyle.Width(detailWidth-4).Render("  $ "+ex) + "\n")
		}
		detailBuilder.WriteString("\n")
	}

	if cmd.Details != "" {
		detailBuilder.WriteString(DetailSectionStyle.Render("🔍 DETALLES & FUNCIONAMIENTO") + "\n")
		detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(cmd.Details) + "\n")
	}

	return listBuilder.String(), detailBuilder.String()
}

// RenderAgentsView renderiza los agentes y subagentes.
func RenderAgentsView(allAgents []*agents.Agent, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	var listBuilder strings.Builder

	for i, a := range allAgents {
		isSelected := i == selectedIdx
		badge := BadgeCoreStyle.Render("Subagente")
		if a.Metadata.Type == "orchestrator" {
			badge = BadgeAgentStyle.Render("Orquestador")
		}

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-16s %s\n  %s", a.Metadata.Name, badge, ItemDescStyle.Render(truncate(a.Metadata.Title, leftWidth-14))),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-16s %s\n  %s", a.Metadata.Name, badge, ItemDescStyle.Render(truncate(a.Metadata.Title, leftWidth-14))),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(allAgents) {
		return listBuilder.String(), "Selecciona un agente para ver sus directrices."
	}

	a := allAgents[selectedIdx]
	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	badge := BadgeCoreStyle.Render("Subagente")
	if a.Metadata.Type == "orchestrator" {
		badge = BadgeAgentStyle.Render("Orquestador")
	}

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("🤖 ROL: "+strings.ToUpper(a.Metadata.Title)) + "  " + badge + "\n")
	detailBuilder.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("Identificador: "+a.Metadata.Name) + "\n\n")

	if a.Metadata.Description != "" {
		detailBuilder.WriteString(DetailSectionStyle.Render("📋 PROPÓSITO") + "\n")
		detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(a.Metadata.Description) + "\n\n")
	}

	detailBuilder.WriteString(DetailSectionStyle.Render("📜 DIRECTRICES & SYSTEM PROMPT COMPLETO") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(strings.TrimSpace(a.Body)) + "\n")

	return listBuilder.String(), detailBuilder.String()
}

// RenderRulesView renderiza las reglas y convenciones.
func RenderRulesView(allRules []*rules.Rule, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	var listBuilder strings.Builder

	for i, r := range allRules {
		isSelected := i == selectedIdx
		badgeCat := BadgeCommandStyle.Render(strings.ToUpper(r.Metadata.Category))

		alwaysOnBadge := ""
		if r.Metadata.AlwaysOn {
			alwaysOnBadge = " " + BadgeWarnStyle.Render("Always-On")
		}

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-16s %s%s\n  %s", r.Metadata.Name, badgeCat, alwaysOnBadge, ItemDescStyle.Render(truncate(r.Metadata.Title, leftWidth-14))),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-16s %s%s\n  %s", r.Metadata.Name, badgeCat, alwaysOnBadge, ItemDescStyle.Render(truncate(r.Metadata.Title, leftWidth-14))),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(allRules) {
		return listBuilder.String(), "Selecciona una regla para ver el contenido."
	}

	r := allRules[selectedIdx]
	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	badgeCat := BadgeCommandStyle.Render(strings.ToUpper(r.Metadata.Category))
	alwaysOnBadge := ""
	if r.Metadata.AlwaysOn {
		alwaysOnBadge = " " + BadgeWarnStyle.Render("Always-On")
	}

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("📏 REGLA: "+r.Metadata.Title) + "  " + badgeCat + alwaysOnBadge + "\n")
	detailBuilder.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("Nombre: "+r.Metadata.Name+" | Categoría: "+r.Metadata.Category) + "\n\n")

	if len(r.Metadata.Tags) > 0 {
		detailBuilder.WriteString(DetailSectionStyle.Render("🏷️ TAGS / TECNOLOGÍAS") + "\n")
		detailBuilder.WriteString(contentStyle.Foreground(ColorHighlight).Render(strings.Join(r.Metadata.Tags, ", ")) + "\n\n")
	}

	if r.Metadata.Description != "" {
		detailBuilder.WriteString(DetailSectionStyle.Render("📋 RESUMEN") + "\n")
		detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(r.Metadata.Description) + "\n\n")
	}

	detailBuilder.WriteString(DetailSectionStyle.Render("📖 NORMATIVA & CONTENIDO COMPLETO") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(strings.TrimSpace(r.Body)) + "\n")

	return listBuilder.String(), detailBuilder.String()
}

// RenderSkillsView renderiza las skills y workflows.
func RenderSkillsView(allSkills []*skills.Skill, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	var listBuilder strings.Builder

	for i, s := range allSkills {
		isSelected := i == selectedIdx
		badge := BadgeSuccessStyle.Render("Skill")

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-18s %s\n  %s", s.Metadata.Name, badge, ItemDescStyle.Render(truncate(s.Metadata.Description, leftWidth-14))),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-18s %s\n  %s", s.Metadata.Name, badge, ItemDescStyle.Render(truncate(s.Metadata.Description, leftWidth-14))),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(allSkills) {
		return listBuilder.String(), "Selecciona una skill para ver la guía completa."
	}

	s := allSkills[selectedIdx]
	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("⚡ SKILL: "+strings.ToUpper(s.Metadata.Name)) + "  " + BadgeSuccessStyle.Render("Engineering Workflow") + "\n\n")

	if s.Metadata.Description != "" {
		detailBuilder.WriteString(DetailSectionStyle.Render("📋 DESCRIPCIÓN") + "\n")
		detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(s.Metadata.Description) + "\n\n")
	}

	detailBuilder.WriteString(DetailSectionStyle.Render("📘 GUÍA Y METODOLOGÍA (SKILL.md COMPLETO)") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(strings.TrimSpace(s.Body)) + "\n")

	return listBuilder.String(), detailBuilder.String()
}

// RenderMCPView renderiza los servidores MCP.
func RenderMCPView(mcpMgr *mcp.Manager, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	catalog := mcpMgr.GetCatalog()
	cfg, _ := mcpMgr.ReadConfig()
	var listBuilder strings.Builder

	for i, m := range catalog {
		isSelected := i == selectedIdx
		isInstalled := false
		if cfg != nil {
			if _, ok := cfg.MCPServers[m.ID]; ok {
				isInstalled = true
			}
		}

		var statusBadge string
		if isInstalled {
			statusBadge = BadgeSuccessStyle.Render("Instalado")
		} else {
			statusBadge = BadgeWarnStyle.Render("No Instalado")
		}

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-14s %s\n  %s", m.ID, statusBadge, ItemDescStyle.Render(truncate(m.Description, leftWidth-14))),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-14s %s\n  %s", m.ID, statusBadge, ItemDescStyle.Render(truncate(m.Description, leftWidth-14))),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(catalog) {
		return listBuilder.String(), "Selecciona un servidor MCP para ver la configuración."
	}

	m := catalog[selectedIdx]
	isInstalled := false
	if cfg != nil {
		if _, ok := cfg.MCPServers[m.ID]; ok {
			isInstalled = true
		}
	}

	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	var statusBadge string
	if isInstalled {
		statusBadge = BadgeSuccessStyle.Render("Instalado y Configurado")
	} else {
		statusBadge = BadgeWarnStyle.Render("No Instalado")
	}

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("🔌 SERVIDOR MCP: "+m.Name) + "  " + statusBadge + "\n")
	detailBuilder.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("ID: "+m.ID+" | Tipo de Transporte: "+m.Type) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("📋 DESCRIPCIÓN") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(m.Description) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("🚀 COMANDO DE INSTALACIÓN") + "\n")
	installCmd := fmt.Sprintf("subikit mcp install %s", m.ID)
	if m.RequiresAuth {
		installCmd += " --token <tu-token>"
	}
	detailBuilder.WriteString(CodeBlockStyle.Width(detailWidth-4).Render("  $ "+installCmd) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("🛠️ CONFIGURACIÓN TÉCNICA") + "\n")
	if m.Type == "stdio" {
		detailBuilder.WriteString(fmt.Sprintf("  • Comando ejecutable: %s\n", KeyBadgeStyle.Render(m.DefaultCmd)))
		if len(m.DefaultArgs) > 0 {
			detailBuilder.WriteString(fmt.Sprintf("  • Argumentos por defecto: %s\n", strings.Join(m.DefaultArgs, " ")))
		}
	} else if m.Type == "http" {
		detailBuilder.WriteString(fmt.Sprintf("  • URL del Servidor: %s\n", KeyBadgeStyle.Render(m.DefaultURL)))
	}
	if m.RequiresAuth {
		detailBuilder.WriteString(fmt.Sprintf("  • Variable de Entorno Token / Header: %s\n", m.EnvKey))
	}

	return listBuilder.String(), detailBuilder.String()
}

// RenderDoctorView renderiza el diagnóstico completo del entorno.
func RenderDoctorView(rulesMgr *rules.Manager, skillsMgr *skills.Manager, mcpMgr *mcp.Manager, agentsMgr *agents.Manager, updateRes *updater.UpdateResult, detailWidth int) string {
	var sb strings.Builder

	antigravityTarget := targets.NewAntigravityTarget()
	absPath, _ := filepath.Abs(".")
	stack, _ := detector.Detect(absPath)

	sb.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("🩺 DIAGNÓSTICO EN VIVO DEL ENTORNO") + "\n\n")

	// 1. Proyecto
	sb.WriteString(DetailSectionStyle.Render("📁 1. PROYECTO ACTUAL & DETECCIÓN DE STACK") + "\n")
	sb.WriteString(fmt.Sprintf("  • Directorio analizado: %s\n", absPath))
	sb.WriteString(fmt.Sprintf("  • Nombre del proyecto:  %s\n", stack.ProjectName))
	if len(stack.Technologies) > 0 {
		sb.WriteString(fmt.Sprintf("  • Tecnologías detectadas: %s\n", BadgeSuccessStyle.Render(strings.Join(stack.Technologies, ", "))))
	} else {
		sb.WriteString("  • Tecnologías detectadas: " + BadgeWarnStyle.Render("Genérico / Sin stack específico") + "\n")
	}
	sb.WriteString("\n")

	// 2. Directrices Locales
	sb.WriteString(DetailSectionStyle.Render("📂 2. ESTADO DE DIRECTRICES LOCALES (.agents/)") + "\n")
	localRulesPath := antigravityTarget.GetProjectRulesPath(absPath)
	if entries, err := os.ReadDir(localRulesPath); err == nil && len(entries) > 0 {
		var active []string
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				active = append(active, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
		sb.WriteString(fmt.Sprintf("  %s Reglas locales activas (%d): %s\n", BadgeSuccessStyle.Render("OK"), len(active), strings.Join(active, ", ")))
	} else {
		sb.WriteString(fmt.Sprintf("  %s No hay reglas en .agents/rules/. Ejecuta 'subikit init' para crearlas.\n", BadgeWarnStyle.Render("PENDIENTE")))
	}

	localSkillsPath := antigravityTarget.GetProjectSkillsPath(absPath)
	if entries, err := os.ReadDir(localSkillsPath); err == nil && len(entries) > 0 {
		var active []string
		for _, e := range entries {
			if e.IsDir() {
				active = append(active, e.Name())
			}
		}
		sb.WriteString(fmt.Sprintf("  %s Skills locales activas (%d): %s\n", BadgeSuccessStyle.Render("OK"), len(active), strings.Join(active, ", ")))
	} else {
		sb.WriteString(fmt.Sprintf("  %s No hay skills en .agents/skills/.\n", BadgeWarnStyle.Render("PENDIENTE")))
	}

	localAgentsPath := antigravityTarget.GetProjectAgentsPath(absPath)
	if entries, err := os.ReadDir(localAgentsPath); err == nil && len(entries) > 0 {
		var active []string
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				active = append(active, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
		sb.WriteString(fmt.Sprintf("  %s Agentes locales activos (%d): %s\n", BadgeSuccessStyle.Render("OK"), len(active), strings.Join(active, ", ")))
	} else {
		sb.WriteString(fmt.Sprintf("  %s No hay agentes en .agents/agents/.\n", BadgeWarnStyle.Render("PENDIENTE")))
	}

	agentsMDPath := filepath.Join(absPath, "AGENTS.md")
	geminiMDPath := filepath.Join(absPath, "GEMINI.md")

	hasAgents := false
	if _, err := os.Stat(agentsMDPath); err == nil {
		hasAgents = true
	}
	hasGemini := false
	if _, err := os.Stat(geminiMDPath); err == nil {
		hasGemini = true
	}

	if hasAgents && hasGemini {
		sb.WriteString(fmt.Sprintf("  %s ¡Conflicto de directrices! Coexisten AGENTS.md y GEMINI.md (duplica tokens en Antigravity). Ejecuta 'subikit sync'.\n", BadgeDangerStyle.Render("CONFLICTO")))
	} else if hasAgents {
		if data, err := os.ReadFile(agentsMDPath); err == nil {
			estTokens := len(data) / 4
			sb.WriteString(fmt.Sprintf("  %s AGENTS.md universal presente (~%d tokens estimados).\n", BadgeSuccessStyle.Render("OK"), estTokens))
		} else {
			sb.WriteString(fmt.Sprintf("  %s AGENTS.md universal presente en la raíz del proyecto.\n", BadgeSuccessStyle.Render("OK")))
		}
	} else {
		sb.WriteString(fmt.Sprintf("  %s AGENTS.md no encontrado en la raíz.\n", BadgeWarnStyle.Render("AVISO")))
	}
	sb.WriteString("\n")

	// 3. Directrices Globales
	sb.WriteString(DetailSectionStyle.Render("🌐 3. CONFIGURACIÓN GLOBAL DEL USUARIO") + "\n")
	globalRulesPath, err := antigravityTarget.GetGlobalRulesPath()
	if err == nil {
		if entries, err := os.ReadDir(globalRulesPath); err == nil && len(entries) > 0 {
			var active []string
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".md") {
					active = append(active, strings.TrimSuffix(e.Name(), ".md"))
				}
			}
			sb.WriteString(fmt.Sprintf("  %s Reglas globales activas (%s): %s\n", BadgeSuccessStyle.Render("OK"), globalRulesPath, strings.Join(active, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("  %s Sin reglas globales. Puedes ejecutar 'subikit init --global'.\n", BadgeWarnStyle.Render("INFO")))
		}
	}
	sb.WriteString("\n")

	// 4. Servidores MCP
	sb.WriteString(DetailSectionStyle.Render("🔌 4. SERVIDORES MCP (MODEL CONTEXT PROTOCOL)") + "\n")
	statuses := mcpMgr.Doctor()
	for _, st := range statuses {
		if st.Installed {
			if st.FoundInPath {
				sb.WriteString(fmt.Sprintf("  %s %s: %s\n", BadgeSuccessStyle.Render("INSTALADO"), st.Name, st.Details))
			} else {
				sb.WriteString(fmt.Sprintf("  %s %s: %s (Ejecutable: %s)\n", BadgeWarnStyle.Render("ADVERTENCIA"), st.Name, st.Details, st.Executable))
			}
		} else {
			sb.WriteString(fmt.Sprintf("  %s %s: %s\n", BadgeCommandStyle.Render("DISPONIBLE"), st.Name, st.Details))
		}
	}

	// 5. Versión y Actualizaciones
	sb.WriteString("\n")
	sb.WriteString(DetailSectionStyle.Render("⚡ 5. ESTADO DE VERSIÓN Y ACTUALIZACIONES") + "\n")
	if updateRes != nil {
		if updateRes.UpdateAvail {
			sb.WriteString(fmt.Sprintf("  %s ¡Nueva versión disponible! Instalada: v%s  ->  Última en GitHub: %s\n", BadgeWarnStyle.Render("UPDATE"), updateRes.CurrentVersion, updateRes.LatestVersion))
			sb.WriteString("  " + CodeBlockStyle.Width(detailWidth-6).Render("Presiona 'u' en la TUI o ejecuta 'subikit update' en tu terminal.") + "\n")
		} else {
			sb.WriteString(fmt.Sprintf("  %s SubiKit está al día con la última versión pública (v%s).\n", BadgeSuccessStyle.Render("AL DÍA"), updateRes.CurrentVersion))
		}
	} else {
		sb.WriteString(fmt.Sprintf("  %s Consultando releases en GitHub en segundo plano...\n", BadgeCommandStyle.Render("BUSCANDO")))
	}

	return sb.String()
}

// GetSDDFeatures busca y evalúa las features existentes en .specs/
func GetSDDFeatures() []SDDFeatureStatus {
	var feats []SDDFeatureStatus
	specsDir := filepath.Join(".", ".specs")

	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return feats
	}

	for _, e := range entries {
		if e.IsDir() {
			featPath := filepath.Join(specsDir, e.Name())
			status := "En progreso"
			isDone := false

			hasSpec := fileExists(filepath.Join(featPath, "spec.md"))
			hasPlan := fileExists(filepath.Join(featPath, "tech-plan.md"))
			hasTasks := fileExists(filepath.Join(featPath, "tasks.md"))
			hasVerify := fileExists(filepath.Join(featPath, "verify.md"))

			if fileExists(filepath.Join(featPath, "archive.md")) {
				data, _ := os.ReadFile(filepath.Join(featPath, "archive.md"))
				if strings.Contains(string(data), "ARCHIVADO_Y_MERGEADO") {
					status = "Archivada / Completada"
					isDone = true
				}
			}

			feats = append(feats, SDDFeatureStatus{
				Name:      e.Name(),
				Path:      featPath,
				Status:    status,
				HasSpec:   hasSpec,
				HasPlan:   hasPlan,
				HasTasks:  hasTasks,
				HasVerify: hasVerify,
				IsDone:    isDone,
			})
		}
	}
	return feats
}

// RenderSDDView renderiza la sección de SDD.
func RenderSDDView(feats []SDDFeatureStatus, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	if len(feats) == 0 {
		left := ListItemStyle.Width(leftWidth - 4).Render("No se encontraron features en .specs/.\n\nUsa 'subikit sdd new <nombre>' para inicializar una.")
		right := DetailTitleStyle.Width(detailWidth-4).Render("📐 SPEC-DRIVEN DEVELOPMENT (SDD)") + "\n\n" +
			DocTextStyle.Width(detailWidth-4).Render("El directorio .specs/ no contiene features activas aún.\n\nPara crear tu primera feature estructurada en 7 fases:\n") +
			CodeBlockStyle.Width(detailWidth-4).Render("  $ subikit sdd new mi-nueva-feature") + "\n\n" +
			DetailSectionStyle.Render("LAS 7 FASES DE SDD") + "\n" +
			DocTextStyle.Width(detailWidth-4).Render("1. Fase 1: Especificación funcional (spec.md)\n2. Fase 2: Clarificación y resolución de dudas\n3. Fase 3: Plan Técnico de arquitectura (tech-plan.md)\n4. Fase 4: Desglose de tareas atómicas (tasks.md)\n5. Fase 5: Implementación guiada\n6. Fase 6: Verificación y pruebas (verify.md)\n7. Fase 7: Archivo y merge (archive.md)")
		return left, right
	}

	var listBuilder strings.Builder
	for i, f := range feats {
		isSelected := i == selectedIdx
		badge := BadgeWarnStyle.Render("En Progreso")
		if f.IsDone {
			badge = BadgeSuccessStyle.Render("Completada")
		}

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-16s %s\n  %s", f.Name, badge, ItemDescStyle.Render(f.Path)),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-16s %s\n  %s", f.Name, badge, ItemDescStyle.Render(f.Path)),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(feats) {
		return listBuilder.String(), "Selecciona una feature para ver su estado."
	}

	f := feats[selectedIdx]
	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	badge := BadgeWarnStyle.Render("En Progreso")
	if f.IsDone {
		badge = BadgeSuccessStyle.Render("Completada")
	}

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("📐 FEATURE SDD: "+f.Name) + "  " + badge + "\n")
	detailBuilder.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("Ruta: "+f.Path) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("📊 CHECKLIST DE LAS 7 FASES") + "\n")

	check := func(ok bool, name string) string {
		if ok {
			return BadgeSuccessStyle.Render("✓") + " " + name + "\n"
		}
		return BadgeWarnStyle.Render("○") + " " + lipgloss.NewStyle().Foreground(ColorDim).Render(name) + "\n"
	}

	detailBuilder.WriteString("  " + check(f.HasSpec, "Fase 1: spec.md (Especificación)"))
	detailBuilder.WriteString("  " + check(f.HasPlan, "Fase 3: tech-plan.md (Plan Técnico)"))
	detailBuilder.WriteString("  " + check(f.HasTasks, "Fase 4: tasks.md (Tareas Atómicas)"))
	detailBuilder.WriteString("  " + check(f.HasVerify, "Fase 6: verify.md (Plan de Verificación)"))
	detailBuilder.WriteString("  " + check(f.IsDone, "Fase 7: archive.md (Feature Archivada)"))
	detailBuilder.WriteString("\n")

	specFile := filepath.Join(f.Path, "spec.md")
	if data, err := os.ReadFile(specFile); err == nil {
		detailBuilder.WriteString(DetailSectionStyle.Render("📄 VISTA PREVIA COMPLETA: spec.md") + "\n")
		detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(string(data)) + "\n")
	}

	return listBuilder.String(), detailBuilder.String()
}

// RenderGlossaryView renderiza el glosario explicativo.
func RenderGlossaryView(terms []GlossaryTerm, selectedIdx int, leftWidth int, detailWidth int) (string, string) {
	var listBuilder strings.Builder

	for i, t := range terms {
		isSelected := i == selectedIdx
		badge := BadgeCommandStyle.Render(t.Category)

		var itemText string
		if isSelected {
			itemText = SelectedListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("▶ %-20s %s\n  %s", t.Term, badge, ItemDescStyle.Render(truncate(t.Summary, leftWidth-14))),
			)
		} else {
			itemText = ListItemStyle.Width(leftWidth - 4).Render(
				fmt.Sprintf("  %-20s %s\n  %s", t.Term, badge, ItemDescStyle.Render(truncate(t.Summary, leftWidth-14))),
			)
		}
		listBuilder.WriteString(itemText + "\n\n")
	}

	if selectedIdx < 0 || selectedIdx >= len(terms) {
		return listBuilder.String(), "Selecciona un término del glosario."
	}

	t := terms[selectedIdx]
	var detailBuilder strings.Builder
	contentStyle := lipgloss.NewStyle().Width(detailWidth - 4)

	detailBuilder.WriteString(DetailTitleStyle.Width(detailWidth-4).Render("📚 CONCEPTO: "+t.Term) + "  " + BadgeCommandStyle.Render(t.Category) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("💡 DEFINICIÓN RÁPIDA") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(t.Summary) + "\n\n")

	detailBuilder.WriteString(DetailSectionStyle.Render("📖 EXPLICACIÓN DIDÁCTICA Y APLICACIÓN COMPLETA") + "\n")
	detailBuilder.WriteString(contentStyle.Foreground(ColorText).Render(t.Explanation) + "\n")

	return listBuilder.String(), detailBuilder.String()
}

// RenderHelpModal renderiza el modal flotante de ayuda.
func RenderHelpModal() string {
	var sb strings.Builder
	sb.WriteString(DetailTitleStyle.Render("⌨️ ATAJOS DE TECLADO Y NAVEGACIÓN NEOVIM") + "\n\n")

	sb.WriteString(DetailSectionStyle.Render("NAVEGACIÓN ENTRE PANELES Y PESTAÑAS") + "\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("l / Right / Enter") + " : Enfocar panel de lectura (modo texto)\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("h / Left / Esc") + "   : Volver al panel de lista izquierda\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("Tab / Shift+Tab") + "   : Cambiar a pestaña siguiente / anterior\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("1 - 8") + "             : Salto directo a pestañas (1:Comandos ... 8:Glosario)\n\n")

	sb.WriteString(DetailSectionStyle.Render("EN MODO LECTURA (PANEL DERECHO - ESTILO NEOVIM)") + "\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("j / ↓") + "             : Bajar línea por línea (Line-by-line scroll)\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("k / ↑") + "             : Subir línea por línea\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("d / Ctrl+d") + "        : Media página hacia abajo (Half-page down)\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("u / Ctrl+u") + "        : Media página hacia arriba (Half-page up)\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("g / Home") + "          : Ir al principio del documento\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("G / End") + "           : Ir al final del documento\n")
	sb.WriteString("  " + KeyBadgeVioletStyle.Render("h / Esc") + "           : Salir del modo lectura y volver a la lista\n\n")

	sb.WriteString(DetailSectionStyle.Render("EN MODO LISTA (PANEL IZQUIERDO)") + "\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("j / ↓") + "             : Bajar en la lista\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("k / ↑") + "             : Subir en la lista\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("g / G") + "             : Primer / último elemento\n\n")

	sb.WriteString("  " + KeyBadgeStyle.Render("?") + "                 : Abrir / Cerrar esta ventana de ayuda\n")
	sb.WriteString("  " + KeyBadgeStyle.Render("q / Ctrl+C") + "        : Salir de la TUI\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorDim).Render("Presiona '?' o 'Esc' para volver."))
	return ModalStyle.Render(sb.String())
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
