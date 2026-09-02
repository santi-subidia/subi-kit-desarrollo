---
name: typescript
title: Estándares de TypeScript Estricto
category: frontend
always_on: false
description: Reglas de tipado estricto, buenas prácticas con interfaces, genéricos, type safety y validación runtime con Zod.
tags: [typescript, frontend, backend, types]
---

# Estándares y Buenas Prácticas de TypeScript

## 1. Type Safety Estricto
- **Prohibido `any`**: Usar `unknown` con type guards si el tipo es verdaderamente dinámico, o tipado genérico `<T>`.
- **Habilitar modo estricto**: Respetar `strict: true`, `noImplicitAny: true` y `strictNullChecks: true`.
- **Evitar Type Assertions inseguras (`as Type`)**: Preferir type predicates (`is`) o validación con esquemas en runtime (Zod / Valibot).

## 2. Declaración de Tipos e Interfaces
- **`interface` para contratos extensibles**: Usar `interface` para definir modelos de datos, contratos de clases o props de componentes.
- **`type` para composiciones**: Usar `type` para unions (`'A' | 'B'`), tuples, intersecciones o utilidades mapeadas (`Record`, `Pick`, `Omit`).
- **Nombres claros**: Nombrar tipos con PascalCase (`UserProfile`, `ApiResponse<T>`).

## 3. Validación en Fronteras (Runtime Validation)
- Los datos provenientes de requests HTTP, formularios, localStorage o APIs externas **deben** validarse en runtime con esquemas Zod.
- Deducir tipos directamente de los esquemas Zod con `z.infer<typeof schema>`.
