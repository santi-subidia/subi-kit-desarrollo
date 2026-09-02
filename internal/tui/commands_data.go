package tui

// GetCommandsData retorna el catálogo exhaustivo de comandos documentados para la TUI.
func GetCommandsData() []CommandDoc {
	return []CommandDoc{
		{
			Name:        "init",
			Category:    "Setup & Directrices",
			Syntax:      "subikit init [--global] [--force] [--all] [--path <dir>]",
			Description: "Analiza el stack tecnológico del proyecto destino e inyecta inteligentemente las reglas, skills y agentes necesarios.",
			Flags: []CommandFlag{
				{Flag: "--global", Description: "Instala las directrices en la configuración global del usuario (~/.gemini/antigravity/rules/)."},
				{Flag: "--force", Description: "Sobrescribe archivos de reglas existentes sin pedir confirmación interactiva."},
				{Flag: "--all", Description: "Instala todas las reglas del catálogo sin importar los tags del stack detectado."},
				{Flag: "--path <dir>", Description: "Ruta del proyecto destino a inicializar (por defecto: directorio actual .)."},
			},
			Examples: []string{
				"subikit init",
				"subikit init --global",
				"subikit init --path ../mi-proyecto-web --force",
				"subikit init --all",
			},
			Details: "El comando analiza package.json, go.mod, pubspec.yaml, requirements.txt, etc., detectando frameworks (React, Next.js, FastAPI, Flutter, Go, etc.) y genera las directrices en .agents/rules/, .agents/skills/, .agents/agents/ y el archivo agnóstico AGENTS.md.",
		},
		{
			Name:        "sync",
			Category:    "Setup & Directrices",
			Syntax:      "subikit sync [--path <dir>]",
			Description: "Sincroniza y actualiza las reglas, skills y agentes locales con las últimas versiones del catálogo embebido.",
			Flags: []CommandFlag{
				{Flag: "--path <dir>", Description: "Ruta del proyecto a sincronizar (por defecto: .)."},
			},
			Examples: []string{
				"subikit sync",
				"subikit sync --path ./backend",
			},
			Details: "Equivale a ejecutar 'subikit init --force'. Es ideal tras actualizar el binario de SubiKit para recibir las últimas directrices y mejoras en los prompts de los agentes.",
		},
		{
			Name:        "update",
			Category:    "Distribución & Mantenimiento",
			Syntax:      "subikit update [--check] [--force] [-y]",
			Description: "Comprueba la disponibilidad de nuevas versiones en GitHub Releases y actualiza el binario de SubiKit automáticamente en caliente.",
			Flags: []CommandFlag{
				{Flag: "--check", Description: "Solo verifica si hay una nueva versión disponible en GitHub sin descargar ni instalar."},
				{Flag: "--force", Description: "Fuerza la reinstalación de la última versión pública aunque sea igual a la actual."},
				{Flag: "-y, --yes", Description: "Acepta la instalación de la actualización sin solicitar confirmación interactiva."},
			},
			Examples: []string{
				"subikit update",
				"subikit update --check",
				"subikit update --force",
				"subikit update -y",
			},
			Details: "Consulta los releases públicos en GitHub, detecta tu sistema operativo (Windows/Linux/macOS) y arquitectura (amd64/arm64), descarga el asset correspondiente y realiza un reemplazo atómico seguro del ejecutable en curso.",
		},
		{
			Name:        "doctor",
			Category:    "Diagnóstico",
			Syntax:      "subikit doctor [--path <dir>]",
			Description: "Realiza un diagnóstico integral del entorno: stack detectado, reglas locales/globales activas y estado de salud de MCPs.",
			Flags: []CommandFlag{
				{Flag: "--path <dir>", Description: "Ruta del proyecto a diagnosticar (por defecto: .)."},
			},
			Examples: []string{
				"subikit doctor",
				"subikit doctor --path ../otra-app",
			},
			Details: "Verifica si los archivos .agents/ y AGENTS.md existen y están poblados, y comprueba si los servidores MCP configurados (como codegraph o context7) tienen sus binarios disponibles en el PATH del sistema.",
		},
		{
			Name:        "agent",
			Category:    "Catálogo & Roles",
			Syntax:      "subikit agent <list | show <nombre> | set-model <nombre> <modelo>>",
			Description: "Gestiona y permite inspeccionar los roles de IA disponibles en el catálogo (Orquestador y Subagentes especializados).",
			Flags: []CommandFlag{
				{Flag: "list", Description: "Muestra todos los agentes disponibles con su tipo y descripción."},
				{Flag: "show <nombre>", Description: "Imprime las directrices y system prompt completo del agente."},
				{Flag: "set-model <nombre> <modelo>", Description: "Modifica el modelo (inherit, flash_lite, flash, pro) del agente local."},
			},
			Examples: []string{
				"subikit agent list",
				"subikit agent show architect",
				"subikit agent set-model architect pro",
				"subikit agent set-model orchestrator flash --global",
			},
			Details: "Permite estudiar los roles definidos para delegar tareas a agentes especializados como el Arquitecto, QA/Testing, Diseñador Frontend, Especialista en Base de Datos o el Orquestador central.",
		},
		{
			Name:        "skill",
			Category:    "Catálogo & Workflows",
			Syntax:      "subikit skill <list | show <nombre>>",
			Description: "Gestiona y muestra las guías de ingeniería, buenas prácticas y habilidades de frontend craftsmanship.",
			Flags: []CommandFlag{
				{Flag: "list", Description: "Lista todas las skills embebidas en SubiKit."},
				{Flag: "show <nombre>", Description: "Muestra el contenido completo (SKILL.md) de la habilidad indicada."},
			},
			Examples: []string{
				"subikit skill list",
				"subikit skill show ui-craftsmanship",
				"subikit skill show sdd-workflow",
				"subikit skill show gemini-cli-rules",
			},
			Details: "Las skills proporcionan estándares de alto nivel para diseño visual refinado, arquitectura limpia, directrices para modelos de lenguaje y metodologías de desarrollo estructuradas.",
		},
		{
			Name:        "rule",
			Category:    "Catálogo & Convenciones",
			Syntax:      "subikit rule <list | show <nombre>>",
			Description: "Permite explorar y consultar las reglas de codificación y convenciones por lenguaje y framework.",
			Flags: []CommandFlag{
				{Flag: "list", Description: "Lista todas las reglas agrupadas por categoría (Core, Frontend, Backend, etc.)."},
				{Flag: "show <nombre>", Description: "Imprime el cuerpo completo de la regla especificada."},
			},
			Examples: []string{
				"subikit rule list",
				"subikit rule show tailwind-css",
				"subikit rule show go-best-practices",
				"subikit rule show react-standards",
			},
			Details: "Cada regla incluye tags para autoselección por stack, indicador de 'Always-On' y pautas de código concisas y accionables para los modelos de IA.",
		},
		{
			Name:        "sdd",
			Category:    "Flujo de Desarrollo",
			Syntax:      "subikit sdd <new <nombre> | status> [--path <dir>] [--author <autor>]",
			Description: "Administra el flujo Spec-Driven Development (SDD) creando la estructura de carpetas de las 7 fases para features.",
			Flags: []CommandFlag{
				{Flag: "new <nombre>", Description: "Crea el directorio .specs/<nombre>/ con las plantillas de las 7 fases de SDD."},
				{Flag: "status", Description: "Muestra el estado de avance de las features activas en .specs/."},
				{Flag: "--path <dir>", Description: "Ruta del proyecto destino."},
				{Flag: "--author <nombre>", Description: "Nombre del autor o responsable de la feature."},
			},
			Examples: []string{
				"subikit sdd new agregar-autenticacion-jwt",
				"subikit sdd new carrito-compras --author 'Santiago'",
				"subikit sdd status",
			},
			Details: "SDD estructura el desarrollo asistido por IA en 7 fases: 1. Especificación (spec.md) -> 2. Clarificación -> 3. Plan Técnico (tech-plan.md) -> 4. Tareas Atómicas (tasks.md) -> 5. Implementación -> 6. Verificación (verify.md) -> 7. Archivo (archive.md).",
		},
		{
			Name:        "mcp",
			Category:    "Servidores MCP",
			Syntax:      "subikit mcp <list | install <id|--all> | doctor> [--token <key>] [--cmd <path>]",
			Description: "Configura, verifica y administra los servidores MCP (Model Context Protocol) en ~/.gemini/config/mcp_config.json.",
			Flags: []CommandFlag{
				{Flag: "list", Description: "Muestra el catálogo de MCPs soportados y si están configurados."},
				{Flag: "install <id>", Description: "Instala y configura un servidor MCP específico en la configuración global."},
				{Flag: "install --all", Description: "Instala todos los servidores MCP disponibles en el catálogo."},
				{Flag: "doctor", Description: "Verifica que los binarios y dependencias de los MCPs configurados funcionen."},
				{Flag: "--token <key>", Description: "Token de autenticación / API Key (ej. para Context7 o servicios cloud)."},
				{Flag: "--cmd <path>", Description: "Comando o ruta personalizada para el ejecutable del servidor stdio."},
			},
			Examples: []string{
				"subikit mcp list",
				"subikit mcp install context7 --token mi_token_secreto",
				"subikit mcp install codegraph",
				"subikit mcp doctor",
			},
			Details: "Soporta servidores esenciales como Context7 (documentación oficial viva), CodeGraph (navegación AST de código), Engram (memoria persistente), Browser, etc. Incluye backup automático antes de modificar el archivo de configuración.",
		},
		{
			Name:        "list",
			Category:    "Catálogo & Inspección",
			Syntax:      "subikit list",
			Description: "Muestra un resumen completo de todo el catálogo embebido en SubiKit: Agentes, Reglas, Skills y MCPs.",
			Flags:       []CommandFlag{},
			Examples: []string{
				"subikit list",
			},
			Details: "Ideal para tener un panorama global de todas las capacidades disponibles en la versión actual de SubiKit desde la línea de comandos.",
		},
		{
			Name:        "tui",
			Category:    "Interfaz Interactiva",
			Syntax:      "subikit tui",
			Description: "Inicia la interfaz de terminal interactiva (TUI) con explorador de comandos, visor de catálogo, diagnósticos y glosario.",
			Flags:       []CommandFlag{},
			Examples: []string{
				"subikit tui",
			},
			Details: "Permite navegar fluidamente con teclado entre pestañas, buscar directrices, consultar el glosario de IA y ver el estado de salud de tu proyecto en tiempo real.",
		},
		{
			Name:        "version",
			Category:    "Información",
			Syntax:      "subikit version (o subikit -v, --version)",
			Description: "Muestra la versión instalada de SubiKit CLI.",
			Flags:       []CommandFlag{},
			Examples: []string{
				"subikit version",
			},
			Details: "Informa sobre la versión del binario compilado.",
		},
	}
}

// GetGlossaryTerms retorna los conceptos clave y explicaciones didácticas del ecosistema.
func GetGlossaryTerms() []GlossaryTerm {
	return []GlossaryTerm{
		{
			Term:     "SubiKit",
			Category: "Ecosistema",
			Summary:  "Dev-Kit especializado para desarrollo de software guiado por IA y Antigravity.",
			Explanation: "SubiKit estandariza la interacción entre desarrolladores y agentes de IA. Proporciona reglas de código adaptativas según el stack, roles de agentes especializados, flujos de ingeniería probados (SDD) y configuración automática de servidores MCP.",
		},
		{
			Term:     "Reglas & Convenciones (Rules)",
			Category: "Directrices",
			Summary:  "Instrucciones normativas y estándares de código que los agentes de IA deben seguir obligatoriamente.",
			Explanation: "Se almacenan en .agents/rules/ y definen convenciones de nombrado, arquitectura de carpetas, manejo de errores, principios SOLID y restricciones técnicas para lenguajes específicos (Go, React, TypeScript, Python, Tailwind, etc.).",
		},
		{
			Term:     "Skills (Habilidades & Workflows)",
			Category: "Directrices",
			Summary:  "Guías y metodologías avanzadas de ingeniería de software para tareas específicas.",
			Explanation: "A diferencia de las reglas (que son normativas cortas), las skills son guías completas que instruyen al agente sobre cómo diseñar interfaces UI de alta calidad (UI Craftsmanship), cómo ejecutar un flujo Spec-Driven Development o cómo interactuar con herramientas específicas.",
		},
		{
			Term:     "Agentes & Subagentes",
			Category: "Roles de IA",
			Summary:  "Especialistas con prompts de sistema optimizados para tareas concretas.",
			Explanation: "• Orquestador: Coordina tareas complejas, desglosa épicas y delega trabajo.\n• Arquitecto: Toma decisiones de diseño y estructura de software.\n• Especialistas (Frontend, Backend, QA, DB): Implementan con maestría su área técnica sin desviarse.",
		},
		{
			Term:     "Model Context Protocol (MCP)",
			Category: "Herramientas & Protocolos",
			Summary:  "Protocolo abierto que permite a los asistentes de IA conectarse a fuentes de datos y herramientas externas.",
			Explanation: "SubiKit gestiona servidores MCP esenciales como:\n• Context7: Obtiene documentación oficial y actualizada de cualquier framework en tiempo real.\n• CodeGraph: Indexa el árbol AST del repositorio para navegación semántica precisa.\n• Engram: Memoria persistente a largo plazo entre sesiones de trabajo.",
		},
		{
			Term:     "Spec-Driven Development (SDD)",
			Category: "Metodología",
			Summary:  "Metodología en 7 fases para desarrollar features con IA minimizando errores y retrabajos.",
			Explanation: "Las 7 Fases de SDD:\n1. Especificación (spec.md): Requisitos funcionales y de negocio.\n2. Clarificación: Preguntas interactivas para eliminar ambigüedades.\n3. Plan Técnico (tech-plan.md): Arquitectura, endpoints y modelos.\n4. Tareas Atómicas (tasks.md): Checklist granular de pasos de implementación.\n5. Implementación: Ejecución ordenada fase a fase.\n6. Verificación (verify.md): Pruebas manuales y automatizadas.\n7. Archivo (archive.md): Registro de cambios y cierre formal.",
		},
		{
			Term:     "Directrices Locales vs Globales",
			Category: "Configuración",
			Summary:  "Diferencia entre directrices aplicadas a un repositorio o a todo el entorno de usuario.",
			Explanation: "• Locales (.agents/): Reglas y skills específicas de la tecnología de ese proyecto. Se versionan en Git y aplican a todo el equipo.\n• Globales (~/.gemini/antigravity/): Preferencias y directrices universales del desarrollador que aplican a todos los proyectos de su máquina.",
		},
		{
			Term:     "AGENTS.md",
			Category: "Compatibilidad",
			Summary:  "Archivo estándar agnóstico en la raíz del proyecto para cualquier herramienta de IA.",
			Explanation: "Aunque Antigravity utiliza la estructura .agents/, SubiKit también genera un AGENTS.md en la raíz del proyecto para que otras herramientas (Claude Desktop, Cursor, Aider, Copilot, etc.) puedan leer y respetar automáticamente las mismas reglas.",
		},
	}
}
