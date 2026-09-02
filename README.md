# SubiKit: Dev-Kit para Desarrollo con IA 🚀

Kit de desarrollo personalizado, portable y multiplataforma para asistentes de código con Inteligencia Artificial (**Antigravity / Gemini CLI**, Cursor, Claude Code y entornos universales).

Compilado como un **único binario en Go** (sin dependencias de Node.js o Python), con catálogo embebido de **Rules & Conventions**, motor de **detección automática de stack**, flujo **SDD (Spec-Driven Development)**, gestión de servidores **MCP** y arquitectura de **Agente Orquestador + Subagentes Especializados**.

---

## ⚡ Características Principales

- **📦 Binario Único y Autocontenido**: Reglas, plantillas, skills y agentes embebidos con `//go:embed`.
- **👑 Agente Orquestador & Subagentes**: Coordinador principal (Tech Lead) que delega tareas a subagentes especializados según la fase del desarrollo.
- **🔄 Flujo SDD (Spec-Driven Development)**: Metodología estricta de 7 fases con **puertas de aprobación obligatorias** y **loop de feedback en verificación**.
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
   Spec & Plan     SQL & RLS                Implementación   Tailwind & a11y
                                                 │
                                                 ▼
                                        [ 🛡️ Code Reviewer ]
                                         Auditoría & Verify
```

| Agente / Subagente | Rol | Responsabilidad Principal |
| :--- | :---: | :--- |
| **`orchestrator`** | **Orquestador (Tech Lead)** | Interfaz principal con el usuario, control del flujo SDD, delegación y validación de puertas de aprobación. |
| **`architect`** | Subagente | Diseño de sistemas, Clean Architecture, contratos de datos, Fases 1 a 3 de SDD (Spec y Plan Técnico). |
| **`fullstack-developer`**| Subagente | Implementación de código tipado (Next.js, React, TypeScript, C# .NET, Node APIs). |
| **`code-reviewer`** | Subagente | Auditoría independiente de calidad, seguridad, detección de gaps vs Spec y Fase 6 (Verificación). |
| **`database-engineer`** | Subagente | Modelado relacional, políticas RLS en Supabase, migraciones SQL reproducibles y optimización EF Core. |
| **`ui-specialist`** | Subagente | Composición con Tailwind CSS, diseño mobile-first, accesibilidad WCAG y micro-interacciones. |

---

## 🔌 Servidores MCP Soportados

| Servidor MCP | Tipo | Propósito | Autenticación / Requisitos |
| :--- | :---: | :--- | :--- |
| **`engram`** | Stdio | Memoria persistente a largo plazo entre sesiones y proyectos. | Binario `engram` en `PATH` |
| **`context7`** | HTTP | Documentación oficial y actualizada en tiempo real de frameworks. | Token Bearer (`--token` o `CONTEXT7_API_KEY`) |
| **`codegraph`** | Stdio | Mapeo semántico de dependencias e impacto de cambios en el código. | Binario `codegraph` en `PATH` |

---

## 📋 Catálogo de Rules & Conventions

- **Core**: `agent-behavior`, `git-conventions`, `clean-code`.
- **Frontend**: `typescript`, `react-nextjs`, `tailwind-css`.
- **Backend**: `dotnet-csharp`, `supabase-sql`, `node-apis`.

---

## 🚀 Uso del CLI (`subikit`)

### 1. Inicializar en un proyecto nuevo
```bash
# Auto-detecta el stack e inyecta reglas, skills y agentes en .agents/ y AGENTS.md
subikit init
```

### 2. Gestión de Agentes y Subagentes
```bash
# Listar todos los agentes disponibles
subikit agent list

# Ver el system prompt y directrices de un agente
subikit agent show architect
subikit agent show orchestrator
```

### 3. Gestión de Servidores MCP
```bash
subikit mcp list
subikit mcp doctor
subikit mcp install context7 --token <tu-token>
subikit mcp install --all --token <tu-token>
```

### 4. Flujo SDD (Spec-Driven Development)
```bash
# Crear nueva feature en .specs/mi-feature/
subikit sdd new mi-feature

# Ver estado de features
subikit sdd status
```

### 5. Configurar globalmente en tu máquina
```bash
subikit init --global
```

### 6. Diagnóstico y catálogo unificado
```bash
subikit list
subikit doctor
```

---

## 🛠️ Compilación Multiplataforma

- **Windows (PowerShell)**: `.\build.ps1`
- **Linux / macOS**: `./build.sh`

Binarios resultantes en `/bin/` (`subikit.exe` y `subikit`).
