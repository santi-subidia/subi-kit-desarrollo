package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	subikit "github.com/santi-subidia/dev-kit-desarrollo"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/agents"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/detector"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/mcp"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/skills"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/targets"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/tui"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/ui"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/updater"
)

var version = "0.4.0"

func main() {
	// Limpieza silenciosa de binarios .old generados en actualizaciones de Windows
	updater.CleanupOld()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	subcommand := os.Args[1]

	rulesFS, err := subikit.GetRulesFS()
	if err != nil {
		ui.Error(fmt.Sprintf("Error al acceder al catálogo de reglas embebido: %v", err))
		os.Exit(1)
	}
	rulesMgr, err := rules.NewManager(rulesFS)
	if err != nil {
		ui.Error(fmt.Sprintf("Error al inicializar gestor de reglas: %v", err))
		os.Exit(1)
	}

	skillsFS, err := subikit.GetSkillsFS()
	if err != nil {
		ui.Error(fmt.Sprintf("Error al acceder al catálogo de skills embebido: %v", err))
		os.Exit(1)
	}
	skillsMgr, err := skills.NewManager(skillsFS)
	if err != nil {
		ui.Error(fmt.Sprintf("Error al inicializar gestor de skills: %v", err))
		os.Exit(1)
	}

	agentsFS, err := subikit.GetAgentsFS()
	if err != nil {
		ui.Error(fmt.Sprintf("Error al acceder al catálogo de agentes embebido: %v", err))
		os.Exit(1)
	}
	agentsMgr, err := agents.NewManager(agentsFS)
	if err != nil {
		ui.Error(fmt.Sprintf("Error al inicializar gestor de agentes: %v", err))
		os.Exit(1)
	}

	mcpMgr := mcp.NewManager()

	switch subcommand {
	case "init":
		handleInit(rulesMgr, skillsMgr, agentsMgr, os.Args[2:])
	case "sync":
		handleSync(rulesMgr, skillsMgr, agentsMgr, os.Args[2:])
	case "list":
		handleList(rulesMgr, skillsMgr, mcpMgr, agentsMgr, os.Args[2:])
	case "doctor":
		handleDoctor(rulesMgr, skillsMgr, mcpMgr, agentsMgr, os.Args[2:])
	case "sdd":
		handleSDD(skillsMgr, os.Args[2:])
	case "mcp":
		handleMCP(mcpMgr, os.Args[2:])
	case "agent":
		handleAgent(agentsMgr, os.Args[2:])
	case "skill":
		handleSkill(skillsMgr, os.Args[2:])
	case "rule":
		handleRule(rulesMgr, os.Args[2:])
	case "tui", "ui", "dashboard":
		handleTUI(rulesMgr, skillsMgr, mcpMgr, agentsMgr, os.Args[2:])
	case "update", "upgrade", "self-update":
		handleUpdate(os.Args[2:])
	case "version", "-v", "--version":
		fmt.Printf("SubiKit CLI v%s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		ui.Error(fmt.Sprintf("Comando desconocido: '%s'", subcommand))
		printUsage()
		os.Exit(1)
	}
}

func handleInit(rulesMgr *rules.Manager, skillsMgr *skills.Manager, agentsMgr *agents.Manager, args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	isGlobal := fs.Bool("global", false, "Instala las directrices en la configuración global de la máquina")
	force := fs.Bool("force", false, "Sobrescribe reglas existentes sin preguntar")
	includeAll := fs.Bool("all", false, "Instala todas las reglas del catálogo sin importar el stack")
	targetPath := fs.String("path", ".", "Ruta del proyecto destino")
	fs.Parse(args)

	ui.PrintBanner()

	antigravityTarget := targets.NewAntigravityTarget()
	agnosticTarget := targets.NewAgnosticTarget()

	allSkills := skillsMgr.GetAll()
	allAgents := agentsMgr.GetAll()

	if *isGlobal {
		ui.Section("Instalación Global")
		coreRules := rulesMgr.GetCoreRules()
		written, err := antigravityTarget.InstallGlobal(coreRules, allSkills, allAgents, *force)
		if err != nil {
			ui.Error(fmt.Sprintf("Fallo en instalación global: %v", err))
			os.Exit(1)
		}
		ui.Success(fmt.Sprintf("Instaladas %d directrices globales en configuración de usuario:", len(written)))
		for _, f := range written {
			ui.Bullet("Elemento Global", filepath.Base(filepath.Dir(f))+"/"+filepath.Base(f))
		}
		return
	}

	// Modo Local (por proyecto)
	absPath, err := filepath.Abs(*targetPath)
	if err != nil {
		absPath = *targetPath
	}

	ui.Section("Detección de Stack")
	stack, err := detector.Detect(absPath)
	if err != nil {
		ui.Error(fmt.Sprintf("Error al detectar stack en %s: %v", absPath, err))
		os.Exit(1)
	}

	ui.Bullet("Proyecto", stack.ProjectName)
	ui.Bullet("Ruta", stack.RootPath)
	if len(stack.Technologies) > 0 {
		ui.Bullet("Tecnologías Detectadas", strings.Join(stack.Technologies, ", "))
	} else {
		ui.Bullet("Tecnologías", "Genérico / Sin stack específico detectado")
	}

	var selectedRules []*rules.Rule
	if *includeAll {
		selectedRules = rulesMgr.GetAll()
	} else {
		selectedRules = rulesMgr.MatchRules(stack.Tags, true)
	}

	ui.Section("Inyección de Reglas, Skills & Agentes")
	written, err := antigravityTarget.InstallProject(absPath, selectedRules, allSkills, allAgents, *force)
	if err != nil {
		ui.Error(fmt.Sprintf("Error al escribir directrices de Antigravity: %v", err))
		os.Exit(1)
	}

	agentsFile, err := agnosticTarget.GenerateAgentsMD(absPath, selectedRules)
	if err != nil {
		ui.Warn(fmt.Sprintf("No se pudo generar AGENTS.md: %v", err))
	} else {
		written = append(written, agentsFile)
	}

	ui.Success(fmt.Sprintf("Configuradas %d reglas, %d skills y %d agentes:", len(selectedRules), len(allSkills), len(allAgents)))
	for _, r := range selectedRules {
		ui.Bullet("Regla ["+r.Metadata.Category+"]", fmt.Sprintf("%s (%s)", r.Metadata.Title, r.Metadata.Name))
	}
	for _, s := range allSkills {
		ui.Bullet("Skill", fmt.Sprintf("%s (%s)", s.Metadata.Name, s.Metadata.Description))
	}
	for _, a := range allAgents {
		tipo := "Subagente"
		if a.Metadata.Type == "orchestrator" {
			tipo = "Orquestador"
		}
		ui.Bullet("Agente ["+tipo+"]", fmt.Sprintf("%s (%s)", a.Metadata.Title, a.Metadata.Name))
	}

	fmt.Println()
	ui.Info(fmt.Sprintf("Directrices listas en: %s/.agents/ y %s/AGENTS.md", absPath, absPath))
}

func handleSync(rulesMgr *rules.Manager, skillsMgr *skills.Manager, agentsMgr *agents.Manager, args []string) {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	targetPath := fs.String("path", ".", "Ruta del proyecto a sincronizar")
	fs.Parse(args)

	handleInit(rulesMgr, skillsMgr, agentsMgr, []string{"-path", *targetPath, "-force"})
}

func handleList(rulesMgr *rules.Manager, skillsMgr *skills.Manager, mcpMgr *mcp.Manager, agentsMgr *agents.Manager, args []string) {
	ui.PrintBanner()

	// 1. Agentes
	ui.Section("Catálogo de Agentes & Subagentes")
	allAgents := agentsMgr.GetAll()
	for _, a := range allAgents {
		tipo := "[Subagente]"
		if a.Metadata.Type == "orchestrator" {
			tipo = "[Orquestador]"
		}
		fmt.Printf("\n  • %-20s %s %s\n", a.Metadata.Name, tipo, a.Metadata.Title)
		if a.Metadata.Description != "" {
			fmt.Printf("    └─ %s\n", a.Metadata.Description)
		}
	}

	// 2. Reglas
	ui.Section("Catálogo de Rules & Conventions")
	cats := rulesMgr.GetCategories()
	allRules := rulesMgr.GetAll()

	for _, cat := range cats {
		fmt.Printf("\n[%s]\n", strings.ToUpper(cat))
		for _, r := range allRules {
			if r.Metadata.Category == cat {
				alwaysOnStr := ""
				if r.Metadata.AlwaysOn {
					alwaysOnStr = " (Siempre Activa)"
				}
				fmt.Printf("  • %-18s %s%s\n", r.Metadata.Name, r.Metadata.Title, alwaysOnStr)
				if r.Metadata.Description != "" {
					fmt.Printf("    └─ %s\n", r.Metadata.Description)
				}
			}
		}
	}

	// 3. Skills
	ui.Section("Catálogo de Skills & Workflows")
	allSkills := skillsMgr.GetAll()
	for _, s := range allSkills {
		fmt.Printf("\n  • %-18s\n", s.Metadata.Name)
		if s.Metadata.Description != "" {
			fmt.Printf("    └─ %s\n", s.Metadata.Description)
		}
	}

	// 4. MCPs
	ui.Section("Catálogo de Servidores MCP")
	mcpCatalog := mcpMgr.GetCatalog()
	cfg, _ := mcpMgr.ReadConfig()
	for _, m := range mcpCatalog {
		installedTag := "[No Instalado]"
		if cfg != nil {
			if _, ok := cfg.MCPServers[m.ID]; ok {
				installedTag = "[Instalado]"
			}
		}
		fmt.Printf("\n  • %-12s %s %s\n", m.ID, installedTag, m.Name)
		fmt.Printf("    └─ %s (Tipo: %s)\n", m.Description, m.Type)
	}
	fmt.Println()
}

func handleDoctor(rulesMgr *rules.Manager, skillsMgr *skills.Manager, mcpMgr *mcp.Manager, agentsMgr *agents.Manager, args []string) {
	fs := flag.NewFlagSet("doctor", flag.ExitOnError)
	targetPath := fs.String("path", ".", "Ruta del proyecto a diagnosticar")
	fs.Parse(args)

	ui.PrintBanner()
	ui.Section("Diagnóstico del Entorno")

	absPath, err := filepath.Abs(*targetPath)
	if err != nil {
		absPath = *targetPath
	}

	antigravityTarget := targets.NewAntigravityTarget()

	// 1. Diagnóstico de Proyecto
	stack, _ := detector.Detect(absPath)
	ui.Bullet("Directorio Analizado", absPath)
	ui.Bullet("Stack Detectado", strings.Join(stack.Technologies, ", "))

	localRulesPath := antigravityTarget.GetProjectRulesPath(absPath)
	if entries, err := os.ReadDir(localRulesPath); err == nil && len(entries) > 0 {
		var active []string
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				active = append(active, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
		ui.Success(fmt.Sprintf("Reglas Locales Activas (.agents/rules/): %s", strings.Join(active, ", ")))
	} else {
		ui.Warn("No se encontraron reglas locales en .agents/rules/. Ejecuta 'subikit init' para crearlas.")
	}

	localSkillsPath := antigravityTarget.GetProjectSkillsPath(absPath)
	if entries, err := os.ReadDir(localSkillsPath); err == nil && len(entries) > 0 {
		var activeSkills []string
		for _, e := range entries {
			if e.IsDir() {
				activeSkills = append(activeSkills, e.Name())
			}
		}
		ui.Success(fmt.Sprintf("Skills Locales Activos (.agents/skills/): %s", strings.Join(activeSkills, ", ")))
	}

	localAgentsPath := antigravityTarget.GetProjectAgentsPath(absPath)
	if entries, err := os.ReadDir(localAgentsPath); err == nil && len(entries) > 0 {
		var activeAgents []string
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				activeAgents = append(activeAgents, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
		ui.Success(fmt.Sprintf("Agentes Locales Activos (.agents/agents/): %s", strings.Join(activeAgents, ", ")))
	}

	// 2. Diagnóstico Global
	globalRulesPath, err := antigravityTarget.GetGlobalRulesPath()
	if err == nil {
		if entries, err := os.ReadDir(globalRulesPath); err == nil && len(entries) > 0 {
			var active []string
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".md") {
					active = append(active, strings.TrimSuffix(e.Name(), ".md"))
				}
			}
			ui.Success(fmt.Sprintf("Reglas Globales Activas (%s): %s", globalRulesPath, strings.Join(active, ", ")))
		}
	}

	// 3. Diagnóstico de Servidores MCP
	ui.Section("Diagnóstico de Servidores MCP")
	mcpStatuses := mcpMgr.Doctor()
	for _, status := range mcpStatuses {
		if status.Installed {
			if status.FoundInPath {
				ui.Success(fmt.Sprintf("%s (%s): %s", status.Name, status.ID, status.Details))
			} else {
				ui.Warn(fmt.Sprintf("%s (%s): %s (Ejecutable: %s)", status.Name, status.ID, status.Details, status.Executable))
			}
		} else {
			ui.Info(fmt.Sprintf("%s (%s): %s", status.Name, status.ID, status.Details))
		}
	}

	// 4. Diagnóstico de Versión y Actualizaciones
	ui.Section("Diagnóstico de Versión y Actualizaciones")
	ui.Bullet("Versión Instalada", "v"+version)
	res, err := updater.CheckLatest(version)
	if err != nil {
		ui.Warn(fmt.Sprintf("No se pudo comprobar actualizaciones en GitHub: %v", err))
	} else if res.UpdateAvail {
		ui.Warn(fmt.Sprintf("¡Actualización disponible en GitHub! %s -> %s (Ejecuta 'subikit update')", version, res.LatestVersion))
	} else {
		ui.Success(fmt.Sprintf("SubiKit está actualizado a la última versión pública (%s)", res.LatestVersion))
	}
	fmt.Println()
}

func handleAgent(agentsMgr *agents.Manager, args []string) {
	if len(args) < 1 {
		fmt.Println("Uso: subikit agent <list|show|set-model> [argumentos]")
		fmt.Println()
		fmt.Println("Comandos de Agentes:")
		fmt.Println("  list            Muestra los roles y subagentes disponibles")
		fmt.Println("  show <nombre>   Muestra las directrices completas del rol")
		fmt.Println("  set-model <nombre> <modelo> Modifica el modelo de un agente local (inherit, flash_lite, flash, pro)")
		return
	}

	action := args[0]
	switch action {
	case "list":
		ui.PrintBanner()
		ui.Section("Catálogo de Agentes y Subagentes")
		allAgents := agentsMgr.GetAll()
		for _, a := range allAgents {
			tipo := "[Subagente]"
			if a.Metadata.Type == "orchestrator" {
				tipo = "[Orquestador]"
			}
			ui.Bullet(a.Metadata.Name, fmt.Sprintf("%s %s - %s", tipo, a.Metadata.Title, a.Metadata.Description))
		}
		fmt.Println()

	case "show":
		if len(args) < 2 {
			ui.Error("Debes especificar el nombre del agente. Ej: subikit agent show architect")
			return
		}
		name := args[1]
		agent, ok := agentsMgr.GetAgent(name)
		if !ok {
			ui.Error(fmt.Sprintf("Agente '%s' no encontrado en el catálogo.", name))
			return
		}
		ui.PrintBanner()
		fmt.Printf("=== ROL: %s (%s) ===\n\n", strings.ToUpper(agent.Metadata.Title), agent.Metadata.Name)
		fmt.Println(agent.Body)
		fmt.Println()
	case "set-model":
		if len(args) < 3 {
			ui.Error("Uso: subikit agent set-model <nombre> <modelo> [--global]")
			return
		}
		name := args[1]
		modelName := args[2]

		validModels := map[string]bool{"inherit": true, "flash_lite": true, "flash": true, "pro": true}
		if !validModels[modelName] {
			ui.Error(fmt.Sprintf("Modelo '%s' no es válido. Opciones: inherit, flash_lite, flash, pro", modelName))
			return
		}

		fs := flag.NewFlagSet("set-model", flag.ExitOnError)
		isGlobal := fs.Bool("global", false, "Modifica el agente en la configuración global")
		fs.Parse(args[3:])

		var agentPath string
		var err error

		t := targets.NewAntigravityTarget()
		if *isGlobal {
			globalPath, errPath := t.GetGlobalRulesPath() // Returns ~/.gemini/config/rules. Better build it
			if errPath != nil {
				ui.Error(fmt.Sprintf("Error obteniendo ruta global: %v", errPath))
				return
			}
			// GetGlobalRulesPath returns config/rules, we need config/agents
			agentPath = filepath.Join(filepath.Dir(globalPath), "agents", name+".md")
		} else {
			absPath, _ := filepath.Abs(".")
			agentPath = filepath.Join(t.GetProjectAgentsPath(absPath), name+".md")
		}

		data, err := os.ReadFile(agentPath)
		if err != nil {
			ui.Error(fmt.Sprintf("No se pudo leer el agente local en %s (¿ejecutaste 'subikit init'?): %v", agentPath, err))
			return
		}

		agent, err := agents.ParseAgent(data, agentPath)
		if err != nil {
			ui.Error(fmt.Sprintf("Error al parsear el agente %s: %v", agentPath, err))
			return
		}

		if err := agent.UpdateModel(modelName); err != nil {
			ui.Error(fmt.Sprintf("Error actualizando modelo: %v", err))
			return
		}

		if err := os.WriteFile(agentPath, []byte(agent.RawContent), 0644); err != nil {
			ui.Error(fmt.Sprintf("Error guardando cambios en %s: %v", agentPath, err))
			return
		}

		ui.Success(fmt.Sprintf("Modelo del agente '%s' actualizado a '%s' en %s", name, modelName, agentPath))

	default:
		ui.Error(fmt.Sprintf("Comando de agente no reconocido: '%s'", action))
	}
}

func handleSkill(skillsMgr *skills.Manager, args []string) {
	if len(args) < 1 {
		fmt.Println("Uso: subikit skill <list|show> [argumentos]")
		fmt.Println()
		fmt.Println("Comandos de Skills:")
		fmt.Println("  list            Muestra todas las skills de ingeniería y frontend disponibles")
		fmt.Println("  show <nombre>   Muestra la guía y directrices completas de la skill")
		return
	}

	action := args[0]
	switch action {
	case "list":
		ui.PrintBanner()
		ui.Section("Catálogo de Skills Embebidas")
		allSkills := skillsMgr.GetAll()
		for _, s := range allSkills {
			ui.Bullet(s.Metadata.Name, s.Metadata.Description)
		}
		fmt.Println()

	case "show":
		if len(args) < 2 {
			ui.Error("Debes especificar el nombre de la skill. Ej: subikit skill show ui-craftsmanship")
			return
		}
		name := args[1]
		skill, ok := skillsMgr.GetSkill(name)
		if !ok {
			ui.Error(fmt.Sprintf("Skill '%s' no encontrada en el catálogo.", name))
			return
		}
		ui.PrintBanner()
		fmt.Printf("=== SKILL: %s ===\n\n", strings.ToUpper(skill.Metadata.Name))
		fmt.Println(skill.Body)
		fmt.Println()

	default:
		ui.Error(fmt.Sprintf("Comando de skill no reconocido: '%s'", action))
	}
}

func handleRule(rulesMgr *rules.Manager, args []string) {
	if len(args) < 1 {
		fmt.Println("Uso: subikit rule <list|show> [argumentos]")
		fmt.Println()
		fmt.Println("Comandos de Reglas:")
		fmt.Println("  list            Muestra todas las reglas de codificación disponibles")
		fmt.Println("  show <nombre>   Muestra el contenido completo de la regla")
		return
	}

	action := args[0]
	switch action {
	case "list":
		ui.PrintBanner()
		ui.Section("Catálogo de Reglas y Convenciones")
		allRules := rulesMgr.GetAll()
		for _, r := range allRules {
			ui.Bullet(r.Metadata.Name, fmt.Sprintf("[%s] %s - %s", r.Metadata.Category, r.Metadata.Title, r.Metadata.Description))
		}
		fmt.Println()

	case "show":
		if len(args) < 2 {
			ui.Error("Debes especificar el nombre de la regla. Ej: subikit rule show tailwind-css")
			return
		}
		name := args[1]
		rule, ok := rulesMgr.GetRule(name)
		if !ok {
			ui.Error(fmt.Sprintf("Regla '%s' no encontrada en el catálogo.", name))
			return
		}
		ui.PrintBanner()
		fmt.Printf("=== REGLA: %s (%s) [%s] ===\n\n", strings.ToUpper(rule.Metadata.Title), rule.Metadata.Name, rule.Metadata.Category)
		fmt.Println(rule.Body)
		fmt.Println()

	default:
		ui.Error(fmt.Sprintf("Comando de regla no reconocido: '%s'", action))
	}
}

func handleMCP(mcpMgr *mcp.Manager, args []string) {
	if len(args) < 1 {
		fmt.Println("Uso: subikit mcp <list|install|doctor> [opciones]")
		fmt.Println()
		fmt.Println("Comandos MCP:")
		fmt.Println("  list                    Muestra los servidores MCP soportados y su estado")
		fmt.Println("  install <id | --all>    Configura servidores MCP en ~/.gemini/config/mcp_config.json")
		fmt.Println("                          Opciones: --token <token>, --cmd <comando>")
		fmt.Println("  doctor                  Verifica la salud y binarios de los MCPs configurados")
		return
	}

	action := args[0]
	switch action {
	case "list":
		ui.PrintBanner()
		ui.Section("Servidores MCP Soportados")
		catalog := mcpMgr.GetCatalog()
		cfg, _ := mcpMgr.ReadConfig()
		for _, m := range catalog {
			installed := "No instalado"
			if cfg != nil {
				if _, ok := cfg.MCPServers[m.ID]; ok {
					installed = "Instalado"
				}
			}
			ui.Bullet(m.ID, fmt.Sprintf("%s [%s] - %s", m.Name, installed, m.Description))
		}
		fmt.Println()

	case "install":
		if len(args) < 2 {
			ui.Error("Debes especificar el ID del servidor o '--all'. Ej: subikit mcp install context7 --token <tu-token>")
			return
		}

		targetID := args[1]
		fs := flag.NewFlagSet("mcp install", flag.ExitOnError)
		token := fs.String("token", "", "Token de autenticación / API Key (ej. para Context7)")
		customCmd := fs.String("cmd", "", "Ruta o comando personalizado para servidores stdio")
		installAll := fs.Bool("all", false, "Instala todos los MCPs del catálogo")
		fs.Parse(args[2:])

		ui.PrintBanner()

		if targetID == "--all" || *installAll {
			ui.Section("Instalación de Todos los MCPs")
			installed, err := mcpMgr.InstallAll(*token)
			if err != nil {
				ui.Error(fmt.Sprintf("Error durante la instalación de MCPs: %v", err))
				return
			}
			ui.Success(fmt.Sprintf("Configurados %d servidores MCP en mcp_config.json:", len(installed)))
			for _, id := range installed {
				ui.Bullet("MCP", id)
			}
		} else {
			ui.Section(fmt.Sprintf("Instalación de MCP: %s", targetID))
			if err := mcpMgr.Install(targetID, *token, *customCmd); err != nil {
				ui.Error(fmt.Sprintf("Error al instalar MCP '%s': %v", targetID, err))
				return
			}
			ui.Success(fmt.Sprintf("Servidor MCP '%s' configurado exitosamente en mcp_config.json", targetID))
		}

		cfgPath, _ := mcpMgr.GetConfigPath()
		ui.Info(fmt.Sprintf("Archivo actualizado: %s (con backup automático)", cfgPath))
		fmt.Println()

	case "doctor":
		ui.PrintBanner()
		ui.Section("Diagnóstico de Servidores MCP")
		statuses := mcpMgr.Doctor()
		for _, st := range statuses {
			if st.Installed {
				if st.FoundInPath {
					ui.Success(fmt.Sprintf("%s: %s", st.Name, st.Details))
				} else {
					ui.Warn(fmt.Sprintf("%s: %s", st.Name, st.Details))
				}
			} else {
				ui.Info(fmt.Sprintf("%s: %s", st.Name, st.Details))
			}
		}
		fmt.Println()

	default:
		ui.Error(fmt.Sprintf("Comando MCP no reconocido: '%s'", action))
	}
}

func handleSDD(skillsMgr *skills.Manager, args []string) {
	if len(args) < 1 {
		fmt.Println("Uso: subikit sdd <new|status> [argumentos]")
		fmt.Println()
		fmt.Println("Comandos SDD:")
		fmt.Println("  new <nombre-feature>  Crea la carpeta .specs/<nombre>/ con las plantillas de las 7 fases")
		fmt.Println("  status                Muestra el estado de las features activas en .specs/")
		return
	}

	action := args[0]
	switch action {
	case "new":
		if len(args) < 2 {
			ui.Error("Debes especificar el nombre de la feature. Ej: subikit sdd new agregar-autenticacion")
			return
		}
		featureName := args[1]
		fs := flag.NewFlagSet("sdd new", flag.ExitOnError)
		targetPath := fs.String("path", ".", "Ruta del proyecto")
		author := fs.String("author", "", "Autor / Responsable de la feature")
		fs.Parse(args[2:])

		absPath, _ := filepath.Abs(*targetPath)
		createdDir, err := skillsMgr.CreateSDDFeature(absPath, featureName, *author)
		if err != nil {
			ui.Error(fmt.Sprintf("Error al crear feature SDD: %v", err))
			return
		}

		ui.PrintBanner()
		ui.Success(fmt.Sprintf("Feature SDD inicializada con éxito en: %s", createdDir))
		ui.Bullet("Fase 1", filepath.Join(createdDir, "spec.md"))
		ui.Bullet("Fase 3", filepath.Join(createdDir, "tech-plan.md"))
		ui.Bullet("Fase 4", filepath.Join(createdDir, "tasks.md"))
		ui.Bullet("Fase 6", filepath.Join(createdDir, "verify.md"))
		ui.Bullet("Fase 7", filepath.Join(createdDir, "archive.md"))
		fmt.Println()
		ui.Info("Comienza completando 'spec.md' y solicita al agente de IA iniciar la Fase 2 (Clarificación).")

	case "status":
		fs := flag.NewFlagSet("sdd status", flag.ExitOnError)
		targetPath := fs.String("path", ".", "Ruta del proyecto")
		fs.Parse(args[1:])

		absPath, _ := filepath.Abs(*targetPath)
		specsDir := filepath.Join(absPath, ".specs")

		entries, err := os.ReadDir(specsDir)
		if err != nil || len(entries) == 0 {
			ui.Info("No hay features activas en .specs/. Usa 'subikit sdd new <nombre>' para crear una.")
			return
		}

		ui.PrintBanner()
		ui.Section("Features SDD Activas (.specs/)")
		for _, e := range entries {
			if e.IsDir() {
				featPath := filepath.Join(specsDir, e.Name())
				status := "Iniciada"
				if fileExists(filepath.Join(featPath, "archive.md")) {
					data, _ := os.ReadFile(filepath.Join(featPath, "archive.md"))
					if strings.Contains(string(data), "ARCHIVADO_Y_MERGEADO") {
						status = "Archivada / Completada"
					}
				}
				ui.Bullet(e.Name(), fmt.Sprintf("Estado: %s (Ruta: .specs/%s)", status, e.Name()))
			}
		}
		fmt.Println()

	default:
		ui.Error(fmt.Sprintf("Subcomando SDD no reconocido: %s", action))
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func handleUpdate(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	checkOnly := fs.Bool("check", false, "Solo comprueba si existe una versión más reciente sin instalar")
	force := fs.Bool("force", false, "Fuerza la reinstalación de la última versión")
	autoYes := fs.Bool("y", false, "Acepta la actualización sin pedir confirmación interactiva")
	fs.BoolVar(autoYes, "yes", false, "Acepta la actualización sin pedir confirmación interactiva")
	fs.Parse(args)

	ui.PrintBanner()
	ui.Section("Actualizador de SubiKit")
	ui.Info("Consultando últimos releases en GitHub...")

	res, err := updater.CheckLatest(version)
	if err != nil {
		ui.Error(fmt.Sprintf("Error al consultar releases en GitHub: %v", err))
		os.Exit(1)
	}

	if res.Release == nil {
		ui.Warn("No se encontraron releases públicos disponibles en GitHub.")
		return
	}

	ui.Bullet("Versión Local", "v"+version)
	ui.Bullet("Último Release en GitHub", res.LatestVersion)

	if *checkOnly {
		if res.UpdateAvail {
			ui.Warn(fmt.Sprintf("¡Nueva versión disponible (%s)! Ejecuta 'subikit update' para instalarla.", res.LatestVersion))
		} else {
			ui.Success("Tu versión de SubiKit está al día.")
		}
		return
	}

	if !res.UpdateAvail && !*force {
		ui.Success(fmt.Sprintf("SubiKit ya está actualizado a la última versión disponible (%s).", "v"+version))
		ui.Info("Usa 'subikit update --force' si deseas forzar la reinstalación del binario.")
		return
	}

	ui.Section(fmt.Sprintf("Instalando SubiKit %s", res.LatestVersion))
	if res.Release.Name != "" {
		ui.Bullet("Nombre del Release", res.Release.Name)
	}
	if res.Release.Body != "" {
		notes := strings.TrimSpace(res.Release.Body)
		if len(notes) > 300 {
			notes = notes[:300] + "..."
		}
		ui.Bullet("Notas de Versión", notes)
	}

	if !*autoYes {
		fmt.Printf("\n¿Deseas descargar e instalar la versión %s ahora? [S/n]: ", res.LatestVersion)
		var answer string
		fmt.Scanln(&answer)
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer == "n" || answer == "no" {
			ui.Info("Actualización cancelada por el usuario.")
			return
		}
	}

	fmt.Println()
	err = updater.ApplyUpdate(res.Release, func(pct float64, status string) {
		ui.Bullet(fmt.Sprintf("%3.0f%%", pct*100), status)
	})

	if err != nil {
		ui.Error(fmt.Sprintf("Fallo al actualizar SubiKit: %v", err))
		os.Exit(1)
	}

	fmt.Println()
	ui.Success(fmt.Sprintf("¡SubiKit se ha actualizado exitosamente a %s!", res.LatestVersion))
	ui.Info("Ejecuta 'subikit version' o 'subikit tui' para utilizar la nueva versión.")
	fmt.Println()
}

func handleTUI(rulesMgr *rules.Manager, skillsMgr *skills.Manager, mcpMgr *mcp.Manager, agentsMgr *agents.Manager, args []string) {
	if err := tui.RunTUI(rulesMgr, skillsMgr, mcpMgr, agentsMgr, version); err != nil {
		ui.Error(fmt.Sprintf("Error al ejecutar la TUI interactiva: %v", err))
		os.Exit(1)
	}
}

func printUsage() {
	ui.PrintBanner()
	fmt.Println("Uso: subikit <comando> [opciones]")
	fmt.Println()
	fmt.Println("Comandos disponibles:")
	fmt.Println("  tui      Inicia la interfaz de terminal interactiva (Dashboard & Glosario)")
	fmt.Println("  init     Inicializa reglas, skills y agentes en el proyecto")
	fmt.Println("           Opciones: --global, --force, --all, --path <dir>")
	fmt.Println("  update   Comprueba y actualiza SubiKit a la última versión de GitHub Releases")
	fmt.Println("           Opciones: --check, --force, -y")
	fmt.Println("  agent    Gestión de Agentes y Subagentes (Orquestador y Especialistas)")
	fmt.Println("           Subcomandos: subikit agent list, subikit agent show <nombre>")
	fmt.Println("  skill    Gestión de Skills de Ingeniería y Frontend Craftsmanship")
	fmt.Println("           Subcomandos: subikit skill list, subikit skill show <nombre>")
	fmt.Println("  rule     Gestión de Reglas y Convenciones de Código")
	fmt.Println("           Subcomandos: subikit rule list, subikit rule show <nombre>")
	fmt.Println("  sdd      Gestión del flujo Spec-Driven Development (SDD)")
	fmt.Println("           Subcomandos: subikit sdd new <nombre>, subikit sdd status")
	fmt.Println("  mcp      Gestión de servidores MCP (Model Context Protocol)")
	fmt.Println("           Subcomandos: subikit mcp list, subikit mcp install, subikit mcp doctor")
	fmt.Println("  sync     Sincroniza y actualiza directrices con el catálogo embebido")
	fmt.Println("  list     Muestra el catálogo completo de Agentes, Rules, Skills y MCPs")
	fmt.Println("  doctor   Diagnostica el estado completo del entorno")
	fmt.Println("  version  Muestra la versión instalada del CLI")
	fmt.Println("  help     Muestra esta ayuda")
	fmt.Println()
}

