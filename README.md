# SubiKit: Dev-Kit para Desarrollo con IA 🚀

Kit de desarrollo personalizado, portable y multiplataforma para asistentes de código con Inteligencia Artificial (**Antigravity / Gemini CLI**, Cursor, Claude Code y entornos universales).

Compilado como un **único binario en Go** (sin dependencias de Node.js o Python), con catálogo embebido de **Rules & Conventions**, motor de **detección automática de stack**, catálogo de **Skills de Ingeniería Rigurosa & Frontend Craftsmanship**, flujo **SDD (Spec-Driven Development)**, sistema de tokens **`DESIGN.md`**, gestión de servidores **MCP** y arquitectura de **Agente Orquestador + Subagentes Especializados**.

---

## ⚡ Características Principales

- **📦 Binario Único y Autocontenido**: Reglas, plantillas, skills y agentes embebidos con `//go:embed`.
- **👑 Agente Orquestador & Subagentes**: Coordinador principal (Tech Lead) que delega tareas a subagentes especializados según la fase del desarrollo.
- **🎨 Frontend Craftsmanship & Anti-Slop**: Erradicación del "AI frontend slop" mediante 4 *Visitor Modes* (`Persuade`, `Operate`, `Read`, `Experience`), el *Craft Floor*, y tokens de color en `OKLCH`.
- **📐 Contrato Visual `DESIGN.md`**: Fuente única de verdad de diseño con tokens portables (paletas OKLCH, escala tipográfica con `clamp()`, espaciado, elevación y radios).
- **🛡️ UI Hardening & Auditoría 5D**: Blindaje ante datos reales extremos (i18n, desbordes de texto, zero-CLS skeletons, 4 estados de UI) y Finish Review en verificación.
- **🔄 Flujo SDD (Spec-Driven Development)**: Metodología estricta de 7 fases con **puertas de aprobación obligatorias**, `CONTEXT.md` (Domain Modeling) y **loop de feedback en verificación**.
- **🧠 Ingeniería de Software Rigurosa (Real Engineering)**: Módulos profundos (*Deep Modules*), *Seams* de desacoplamiento, *Design It Twice*, diagnóstico científico de bugs (*Tight Loops*) y mapas de decisión (*Wayfinder*).
- **🔌 Gestión de Servidores MCP**: Configuración y diagnóstico de `engram`, `context7` y `codegraph` con backups automáticos de `mcp_config.json`.
- **🔍 Detección Automática de Stack**: Inspecciona proyectos simples y monorepos (`package.json`, `tsconfig.json`, `.NET / C#`, `go.mod`, `supabase/`, etc.).
- **🌐 Modo Híbrido**: Local (`.agents/` y `AGENTS.md`) o Global (`~/.gemini/config/`).

---

## 👥 Arquitectura de Agentes y Subagentes

```
                       [ 👑 AGENTE ORQUESTADOR ]
                          (Tech Lead & Coordinador)
                                     │
         ┌──────────────┬────────────┴───────────┬──────────────┐
         ▼              ▼                        ▼              ▼
  [ 🏛️ Architect ] [ 🗄️ DBA ]              [ 💻 Fullstack ] [ 🎨 UI/UX ]
   Spec & Plan     SQL & RLS                Implementación   Art Direction
                                                 │           & Craftsmanship
                                                 ▼
                                        [ 🛡️ Code Reviewer ]
                                         Auditoría & Finish Review
```

| Agente / Subagente | Rol | Responsabilidad Principal |
| :--- | :---: | :--- |
| **`orchestrator`** | **Orquestador (Tech Lead)** | Interfaz principal con el usuario, control del flujo SDD, delegación, resolución metódica de bugs y validación de puertas de aprobación. |
| **`architect`** | Subagente | Diseño de sistemas, Clean Architecture, *Deep Modules*, *Seams*, *Design It Twice*, Fases 1 a 3 de SDD (`spec.md`, `CONTEXT.md`, `tech-plan.md` y `docs/adr/`). |
| **`fullstack-developer`**| Subagente | Implementación de código tipado (Next.js, React, TypeScript, C# .NET, Node APIs) y TDD con feedback loops. |
| **`code-reviewer`** | Subagente | Auditoría independiente de calidad, seguridad, profundidad de módulos, detección de gaps vs Spec y **UI Finish Review** (`rebuild`, `fix`, `ship`). |
| **`database-engineer`** | Subagente | Modelado relacional, políticas RLS en Supabase, migraciones SQL reproducibles y optimización EF Core. |
| **`ui-specialist`** | Subagente | **Director de Arte e Ingeniero UI**: Aplicación de *Visitor Modes*, cumplimiento de *Craft Floor*, contrato `DESIGN.md`, Tailwind CSS y accesibilidad WCAG AA. |

---

## 🧰 Catálogo de Skills Embebidas

| Skill | Propósito |
| :--- | :--- |
| **`ui-craftsmanship`** | Dirección de arte frontend, 4 Visitor Modes, Craft Floor, prohibición de anti-patrones de IA y sistema de tokens `DESIGN.md`. |
| **`ui-hardening-audit`** | Auditoría técnica en 5 dimensiones (A11y, Performance, Theming, Responsive, Integridad), blindaje contra datos extremos y checklist de Polish pre-ship. |
| **`sdd-workflow`** | Ciclo de 7 fases de desarrollo guiado por especificaciones con puertas de aprobación. |
| **`domain-modeling`** | Construcción continua del glosario de términos (`CONTEXT.md`) y registros de decisiones (`docs/adr/`). |
| **`diagnosing-bugs`** | Diagnóstico científico de bugs en 4 fases; **prohibido modificar código sin un Tight Feedback Loop en ROJO**. |
| **`codebase-design`** | Vocabulario de diseño de software: *Deep Modules*, *Seams* de desacoplamiento y práctica *Design It Twice*. |
| **`wayfinder`** | Planificación de iniciativas masivas (Epics) mediante mapas de decisiones y grafos de tareas desbloqueables. |
| **`wait-what`** | Comando de freno de mano y re-enfoque cuando el asistente genera explicaciones confusas o verbosas. |
| **`dotnet-hardening`** | Optimización de consultas EF Core (sargabilidad, zero N+1, AsSplitQuery), rendimiento C# (Zero Sync-over-Async, Span, ArrayPool), higiene MSBuild y aserciones rigurosas de tests. |

---

## 🔌 Servidores MCP Soportados

| Servidor MCP | Tipo | Propósito | Autenticación / Requisitos |
| :--- | :---: | :--- | :--- |
| **`engram`** | Stdio | Memoria persistente a largo plazo entre sesiones y proyectos. | Binario `engram` en `PATH` |
| **`context7`** | HTTP | Documentación oficial y actualizada en tiempo real de frameworks. | Token Bearer (`--token` o `CONTEXT7_API_KEY`) |
| **`codegraph`** | Stdio | Mapeo semántico de dependencias e impacto de cambios en el código. | Binario `codegraph` en `PATH` |

---

## 📋 Catálogo de Rules & Conventions

- **Core**: `agent-behavior` (concisión, lenguaje de `CONTEXT.md`, feedback loop), `git-conventions`, `clean-code` (*Deep Modules*, *Design It Twice*).
- **Frontend**: `typescript`, `react-nextjs` (resiliencia, 4 estados de UI, zero-CLS skeletons), `tailwind-css` (tokens semánticos, OKLCH, anti-slop, fluid type con clamp).
- **Backend**: `dotnet-csharp`, `supabase-sql`, `node-apis`.

---

---

## ⚡ Instalación Rápida (One-Liner)

Para instalar SubiKit en cualquier máquina nueva (PC, notebook, servidores) sin necesidad de tener Go o Git instalados:

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/santi-subidia/dev-kit-desarrollo/main/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://raw.githubusercontent.com/santi-subidia/dev-kit-desarrollo/main/install.sh | bash
```

### O con Go (Modo Desarrollador)
```bash
go install github.com/santi-subidia/dev-kit-desarrollo/cmd/subikit@latest
```

---

## 🔄 Auto-Actualización y Portabilidad

SubiKit incluye un motor de auto-actualización multiplataforma que se conecta a GitHub Releases:

```bash
# 1. Comprobar si hay una nueva versión disponible
subikit update --check

# 2. Actualizar el binario automáticamente a la última versión
subikit update

# 3. Forzar reinstalación de la versión actual
subikit update --force -y
```

> [!TIP]
> **Desde la TUI (`subikit tui`)**: Si hay una nueva versión disponible, verás una notificación en la cabecera y en la pestaña **Doctor**. Puedes presionar la tecla `u` para descargar e instalar la actualización directamente.

---

## 🚀 Uso del CLI (`subikit`)

### 0. TUI Interactiva (Dashboard & Glosario)
```bash
# Inicia la interfaz de terminal interactiva para explorar comandos, agentes, reglas, skills y glosario
subikit tui
```

### 1. Inicializar en un proyecto nuevo
```bash
# Auto-detecta el stack e inyecta reglas, skills y agentes en .agents/ y AGENTS.md
subikit init
```

### 2. Sincronizar directrices del proyecto
```bash
# Tras actualizar SubiKit, actualiza las reglas locales del proyecto con las últimas mejoras
subikit sync
```

### 3. Gestión de Skills y Agentes
```bash
# Listar todas las skills disponibles
subikit skill list

# Ver contenido de una skill
subikit skill show ui-craftsmanship
subikit skill show ui-hardening-audit
subikit skill show domain-modeling

# Listar y ver agentes
subikit agent list
subikit agent show ui-specialist
subikit agent show code-reviewer
```

### 4. Gestión de Servidores MCP
```bash
subikit mcp list
subikit mcp doctor
subikit mcp install context7 --token <tu-token>
subikit mcp install --all --token <tu-token>
```

### 5. Flujo SDD (Spec-Driven Development)
```bash
# Crear nueva feature en .specs/mi-feature/
subikit sdd new mi-feature

# Ver estado de features
subikit sdd status
```

### 6. Configurar globalmente en tu máquina
```bash
subikit init --global
```

### 7. Diagnóstico y catálogo unificado
```bash
subikit list
subikit doctor
```

---

## 🛠️ Compilación Multiplataforma Local

- **Windows (PowerShell)**: `.\build.ps1`
- **Linux / macOS**: `./build.sh`

Binarios resultantes en `/bin/` (`subikit.exe` y `subikit`).

