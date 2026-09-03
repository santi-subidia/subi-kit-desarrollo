---
name: task-routing
title: Enrutamiento de Flujos de Trabajo (SDD vs Fast-Fix vs Direct-Tweak)
category: core
always_on: true
description: Clasifica cada solicitud en tres vías (Macro SDD, Micro Fast-Fix, Nano Direct-Tweak) para eliminar la sobrecarga documental en bugs pequeños o tareas simples.
tags: [core, routing, sdd, fast-fix, direct-tweak, workflows]
---

# Enrutamiento de Flujos de Trabajo: Macro vs Micro vs Nano 🚦

El Agente Orquestador debe clasificar toda solicitud del usuario en una de las tres vías de ejecución antes de actuar. Esto evita la sobrecarga ceremonial de SDD en bugs pequeños o cambios triviales, reservando la planificación exhaustiva para iniciativas que realmente lo ameriten.

---

## 🧭 Matriz de Decisión y Vías de Ejecución

| Criterio | Vía Nano: Direct-Tweak | Vía Micro: Fast-Fix | Vía Macro: SDD Completo |
| :--- | :--- | :--- | :--- |
| **Tipo de Tarea** | Typos, textos, constantes, config menor | Bugs, excepciones, regresiones localizadas | Features nuevas, refactors estructurales |
| **Archivos Afectados** | 1 archivo | 1 a 2 archivos | 3+ archivos o módulos cruzados |
| **Riesgo Técnico** | Nulo | Medio (corrige comportamiento) | Alto (arquitectura, APIs, contratos) |
| **Incertidumbre** | Nula (la solución es evidente) | Baja (se identifica la causa o el seam) | Alta (requiere diseño de interfaces/seams) |
| **Documentación** | Cero artefactos | Test de regresión (cero Markdown) | `spec.md`, `tech-plan.md`, `tasks.md` |
| **Delegación** | In-line / Directo por el Orquestador | Subagente puntual (`fullstack`/`database`) | Equipo de subagentes por fases |

---

## 1. 🟢 Vía Nano: Direct-Tweak (Sin Ceremonia)
- **Cuándo aplica**:
  - Corrección de typos o cadenas de texto en UI.
  - Actualización de una constante numérica o booleana.
  - Bump de versión o cambio de un setting en archivo de configuración.
  - Renombrado simple de un identificador local.
- **Protocolo de ejecución**:
  1. Aplicar la edición quirúrgica directamente con la herramienta de código.
  2. Ejecutar comando rápido de compilación o linter (`go build`, `npm run build`, `dotnet build`) para verificar sintaxis.
  3. Realizar commit atómico inmediato (`style`, `chore`, `refactor`).
  4. Responder concisamente al usuario indicando el archivo modificado.

---

## 2. 🟡 Vía Micro: Fast-Fix (Científico, sin Markdown)
- **Cuándo aplica**:
  - Reportes de bugs o errores de ejecución.
  - Regresiones de funcionalidad existente.
  - Comportamiento inesperado o excepciones en un endpoint o handler puntual.
  - Cambios que tocan 1 o 2 archivos como máximo.
- **Protocolo de ejecución (`fast-fix`)**:
  - **Prohibido crear `spec.md`, `tech-plan.md` o `tasks.md`**. El contrato de verdad es el test.
  1. **Reproducción (ROJO)**: Escribir un test unitario, script cURL o invocación que demuestre el fallo de forma determinística en pantalla.
  2. **Fix Quirúrgico**: Modificar el código mínimo indispensable en el archivo afectado (directo o delegando a un subagente ejecutor).
  3. **Demostración (VERDE)**: Correr el test de reproducción hasta que pase a verde, seguido de la suite de pruebas del proyecto para asegurar cero regresiones colaterales.
  4. **Commit de Regresión**: Conservar el test como prueba permanente del repositorio y hacer commit: `fix(scope): ...`.

---

## 3. 🔴 Vía Macro: SDD Completo (Spec-Driven Development)
- **Cuándo aplica**:
  - Creación de nuevas características o endpoints desde cero.
  - Rediseños de base de datos, nuevas entidades o migraciones complejas.
  - Refactorizaciones de arquitectura o cambios que impactan 3 o más archivos.
  - Tareas con alta incertidumbre técnica o dependencias externas.
- **Protocolo de ejecución (`sdd-workflow`)**:
  1. **Fase 1**: Spec y Modelado (`spec.md` con Gherkin) -> Puerta de Aprobación.
  2. **Fase 2**: Clarificación de dependencias y edge cases.
  3. **Fase 3**: Plan Técnico (`tech-plan.md`, Deep Modules, Seams, ADRs) -> Puerta de Aprobación.
  4. **Fase 4**: Desglose atómico de tareas (`tasks.md`).
  5. **Fase 5**: Implementación delegada a subagentes (`fullstack`, `database`, `ui`).
  6. **Fase 6**: Verificación por `code-reviewer` (Finish Review) -> Puerta de Aprobación Final.
  7. **Fase 7**: Cierre, merge y commit.
