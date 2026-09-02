# Plan Técnico: {{FEATURE_NAME}}

- **Spec Asociada**: [spec.md](./spec.md)
- **Estado**: `[BORRADOR | EN_REVISION | APROBADO]`
- **Documentación de Referencia**: `CONTEXT.md`, `docs/adr/`, `DESIGN.md`, `PRODUCT.md`

---

## 1. Análisis Arquitectónico e Impacto
<!-- Cómo encaja este cambio en la arquitectura existente y qué capas impacta -->
- **Frontend**: Componentes, Hooks, State Management, Server Actions / Route Handlers.
- **Backend / Application**: Casos de uso, Handlers (MediatR / Services), Validadores (Zod / FluentValidation).
- **Backend / Domain**: Entidades, Value Objects, Domain Events (sincronizados con `CONTEXT.md`).
- **Persistencia / Database**: Tablas, RLS, Migraciones SQL, Índices.

---

## 2. Deep Modules & Costuras (Seams)
<!-- Principios de John Ousterhout y Michael Feathers -->
- **Módulos Profundos**: Definición de interfaces limpias que ocultan la complejidad de implementación.
- **Costuras (Seams) Identificadas**: Puntos de desacoplamiento para inyección de dependencias y pruebas unitarias aisladas sin mocks frágiles.

---

## 3. Contratos de Datos y Esquemas

### 3.1 Esquemas / DTOs
```typescript
// O C# DTO / Record / Zod Schema
```

### 3.2 Endpoints / APIs
- `METHOD /api/v1/...`
  - **Request**: `{ ... }`
  - **Response 200**: `{ success: true, data: { ... } }`
  - **Response 4xx/5xx**: `{ success: false, error: { ... } }`

---

## 4. Opciones de Diseño Evaluadas (Design It Twice)
<!-- Evaluar siempre 2 alternativas antes de congelar el diseño -->
1. **Alternativa A (Seleccionada)**: [Descripción, ventajas y por qué se elige].
2. **Alternativa B (Descartada)**: [Descripción, por qué se descartó].
- **ADR Vinculado**: `docs/adr/XXXX-{{FEATURE_NAME}}.md` *(si aplica para decisiones estructurales)*.

---

## 5. Estrategia de Testing (Tight Feedback Loops)
- **Unit Tests**: Casos determinísticos y funciones clave con ejecución < 2s.
- **Integration Tests**: Flujos con base de datos o APIs simuladas.
- **Loop de Reproducción de Regresión**: Verificación automatizada de criterios de aceptación.

---

## 6. Puerta de Aprobación del Plan Técnico
- [ ] **Aprobación del Usuario**: Requiere confirmación explícita del diseño técnico antes de generar el desglose de tareas.
- **Notas / Ajustes del Revisor**:
