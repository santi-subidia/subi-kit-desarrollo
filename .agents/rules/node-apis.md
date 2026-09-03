---
name: node-apis
title: APIs Backend, Manejo de Errores y Validación
category: backend
always_on: false
description: Estructura de APIs REST/Route Handlers, validación con Zod, códigos de estado HTTP y manejo centralizado de excepciones.
tags: [node, backend, api, rest, zod]
---

# Desarrollo de APIs Backend y Route Handlers

## 1. Validación de Payloads y Parámetros
- **Validación con Zod**: Todo endpoint debe validar `body`, `query` y `params` antes de ejecutar cualquier lógica de negocio.
- **Códigos de Error HTTP Apropiados**:
  - `400 Bad Request`: Parámetros inválidos o error de validación de esquema.
  - `401 Unauthorized`: Falta de sesión o token no provisto/expirado.
  - `403 Forbidden`: Usuario autenticado sin permisos suficientes.
  - `404 Not Found`: Recurso no encontrado.
  - `422 Unprocessable Entity`: Error semántico de validación.
  - `500 Internal Server Error`: Excepción no controlada (no exponer stack traces en producción).

## 2. Formato Consistente de Respuesta
- Respuestas exitosas:
  ```json
  {
    "success": true,
    "data": { ... }
  }
  ```
- Respuestas de error:
  ```json
  {
    "success": false,
    "error": {
      "code": "INVALID_INPUT",
      "message": "Descripción legible del error",
      "details": [ ... ]
    }
  }
  ```

## 3. Seguridad
- Sanitizar inputs para evitar inyecciones.
- Implementar rate limiting y CORS estricto en rutas públicas o de autenticación.
- Proteger secretos y variables de entorno usando `process.env` tipado o bibliotecas tipo `@t3-oss/env-nextjs`.
