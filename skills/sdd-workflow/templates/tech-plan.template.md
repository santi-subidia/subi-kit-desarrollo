# Plan Técnico: {{FEATURE_NAME}}

- **Spec Asociada**: [spec.md](./spec.md)
- **Estado**: `[BORRADOR | EN_REVISION | APROBADO]`
- **Documentación de Referencia**: `DESIGN.md`, `PRODUCT.md`, `ARCHITECTURE.md`

---

## 1. Análisis Arquitectónico e Impacto
<!-- Cómo encaja este cambio en la arquitectura existente y qué capas impacta -->
- **Frontend**: Componentes, Hooks, State Management, Server Actions / Route Handlers.
- **Backend / Application**: Casos de uso, Handlers (MediatR / Services), Validadores (Zod / FluentValidation).
- **Backend / Domain**: Entidades, Value Objects, Domain Events.
- **Persistencia / Database**: Tablas, RLS, Migraciones SQL, Índices.

## 2. Contratos de Datos y Esquemas

### 2.1 Esquemas / DTOs
```typescript
// O C# DTO / Record
```

### 2.2 Endpoints / APIs
- `METHOD /api/v1/...`
  - **Request**: `{ ... }`
  - **Response 200**: `{ success: true, data: { ... } }`
  - **Response 4xx/5xx**: `{ success: false, error: { ... } }`

## 3. Decisiones de Diseño y Trade-offs
- **Decisión 1**: Por qué se elige el Enfoque A sobre el Enfoque B.
- **Seguridad & Performance**: Índices, transacciones atómicas, validación en fronteras.

## 4. Estrategia de Testing
- **Unit Tests**: Casos y funciones clave a probar.
- **Integration Tests**: Flujos con base de datos o APIs simuladas.

---

## 5. Puerta de Aprobación del Plan Técnico
- [ ] **Aprobación del Usuario**: Requiere confirmación explícita del diseño técnico antes de generar el desglose de tareas.
- **Notas / Ajustes del Revisor**:
