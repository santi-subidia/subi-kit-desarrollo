# Verificación de Aceptación: {{FEATURE_NAME}}

- **Spec Original**: [spec.md](./spec.md)
- **Plan Técnico**: [tech-plan.md](./tech-plan.md)
- **Estado de Verificación**: `[EN_EVALUACION | OBSERVACIONES_DETECTADAS | APROBADO_TOTAL]`

---

## 1. Matriz de Trazabilidad (Criterios vs Implementación)

| Criterio de Aceptación (Spec) | Estado | Evidencia / Método de Prueba | Observaciones / Gap |
| :--- | :---: | :--- | :--- |
| **Criterio 1**: [Descripción] | `[CUMPLE / NO CUMPLE]` | Test unitario `xTest` / Verificación UI | - |
| **Criterio 2**: [Manejo de Errores] | `[CUMPLE / NO CUMPLE]` | Simulación de payload inválido | - |

---

## 2. Reporte de Gaps y Discrepancias (Feedback Loop)

> [!WARNING]
> Si algún criterio **NO CUMPLE**, se debe documentar aquí la causa exacta y no se permite archivar la feature hasta resolver el ciclo.

### Gap Detectado:
- **Qué falló**: 
- **Por qué no cumple con la Spec**: 
- **Acción Correctiva Requerida**:
  - [ ] ¿Requiere ajustar el Plan Técnico? `(Sí/No)`
  - [ ] Nuevas tareas añadidas a `tasks.md` (Sección 4).

---

## 3. Pruebas Automatizadas y Calidad
- **Tests Unitarios**: `[PASS | FAIL]` (Comando: `...`)
- **Tests de Integración / Build**: `[PASS | FAIL]` (Comando: `...`)
- **Linter / Typecheck**: `[PASS | FAIL]` (Comando: `...`)

---

## 4. Puerta de Aprobación de la Verificación
- [ ] **Aprobación Final del Usuario**: Confirmación explícita de que todos los criterios han sido satisfechos satisfactoriamente antes de archivar.
