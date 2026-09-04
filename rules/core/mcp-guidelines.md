---
name: mcp-guidelines
title: Protocolos de Uso de MCPs (Context7, CodeGraph, Engram)
category: core
always_on: true
description: Protocolos estrictos para el uso de servidores MCP predeterminados para documentación, análisis de código y memoria persistente.
tags: [core, mcp, context7, codegraph, engram]
---

# Protocolos de Servidores MCP Predeterminados

Como agente, tienes acceso a servidores MCP (Model Context Protocol) que extienden tus capacidades. Es obligatorio seguir estas directrices al interactuar con ellos.

## 1. Context7 (Documentación Externa)

**Cuándo usarlo**: Siempre que el usuario pregunte sobre librerías, frameworks, SDKs, APIs, herramientas CLI o servicios cloud (ej. React, Next.js, Prisma, Tailwind, etc.). Úsalo incluso si crees saber la respuesta, ya que tu conocimiento base podría estar desactualizado.
**No usar para**: Refactorizaciones, escribir scripts desde cero, depurar lógica de negocio o conceptos generales de programación.

**Pasos obligatorios**:
1. Llama a `resolve-library-id` usando el nombre de la librería y lo que deseas buscar (a menos que el usuario provea el ID exacto en formato `/org/project`).
2. Elige el mejor resultado basándote en coincidencia exacta, relevancia, reputación y puntaje.
3. Llama a `query-docs` con el ID de librería seleccionado y el concepto específico a buscar. Si la pregunta abarca múltiples conceptos distintos, haz llamadas separadas a `query-docs` (no las mezcles en una sola).
4. Responde utilizando la documentación oficial obtenida.

## 2. CodeGraph (Análisis Estructural del Código)

**Cuándo usarlo**: Para preguntas estructurales o sobre el código (repo maps, arquitectura, flujos de llamadas, dependencias, referencias de símbolos, análisis de impacto o "cómo funciona X").
**Regla de orden**: Debes usar CodeGraph **antes** de realizar búsquedas amplias en el sistema de archivos (Read/Glob/Grep).

**Pasos obligatorios**:
1. Confirma que la raíz es un proyecto real (no ejecutes CodeGraph en tu `$HOME` o directorios temporales).
2. Verifica si existe el directorio `<project-root>/.codegraph/`.
3. Si el índice no existe, invoca la herramienta de inicialización (ej. `codegraph init <project-root>`) para crear el índice.
4. Llama a la herramienta `codegraph_explore` para explorar los símbolos y flujos de llamadas.
5. **Solo como respaldo**: Pasa a usar herramientas de sistema de archivos normales (Grep/Read) únicamente si la inicialización o consulta de CodeGraph falla, explicando brevemente el fallo.
6. **Aislamiento en Git**: El directorio `.codegraph/` generado para indexar el repositorio es de uso exclusivamente local. **Debe estar siempre incluido en `.gitignore` y jamás debe incluirse en commits**.

## 3. Engram (Memoria Persistente)

Engram es el sistema de memoria que sobrevive entre sesiones y compactaciones. Debes gestionar el estado del proyecto activamente.

**Cuándo GUARDAR (`mem_save`) - Obligatorio**:
Llama a `mem_save` inmediatamente después de:
- Completar la solución de un bug.
- Tomar una decisión de arquitectura o diseño.
- Hacer un descubrimiento no obvio sobre el código.
- Cambiar configuración o establecer un patrón/convención.

**Cuándo BUSCAR (`mem_search` / `mem_context`)**:
- Cuando el usuario pida recordar algo ("qué hicimos", "recordar", etc.), llama primero a `mem_context` (reciente) y luego `mem_search`.
- **Proactivamente**: Al iniciar trabajo en algo que podría haberse tocado antes, o cuando el usuario menciona un tema del que no tienes contexto actual.

**Protocolo de Cierre de Sesión - Obligatorio**:
Antes de dar por terminada una tarea o sesión, llama a `mem_session_summary` con:
- Goal (Objetivo)
- Instructions (Preferencias aprendidas)
- Discoveries (Hallazgos)
- Accomplished (Logros)
- Next Steps (Próximos pasos)
- Relevant Files (Archivos clave)

**Captura Pasiva**:
Al completar una tarea, puedes incluir una sección `## Key Learnings:` numerada al final de tu respuesta para que Engram la capture automáticamente, o llamar a `mem_capture_passive`.
