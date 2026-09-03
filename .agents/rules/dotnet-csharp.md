---
name: dotnet-csharp
title: Arquitectura .NET, C# y Clean Architecture
category: backend
always_on: false
description: Buenas prácticas de C# moderno, Clean Architecture, DDD, CQRS, Result Pattern, Entity Framework Core y testing.
tags: [dotnet, csharp, backend, clean-architecture, ddd]
---

# Arquitectura .NET, C# y Clean Architecture

## 1. Principios de Clean Architecture y DDD
- **Independencia del Dominio**: La capa `Domain` no debe tener dependencias externas (sin Entity Framework, sin HTTP, sin librerías de infraestructura).
- **Entidades y Value Objects**: Encapsular reglas de negocio e invariantes dentro de entidades ricas en comportamiento (evitar modelos anémicos cuando aplique).
- **Separación de Responsabilidades**:
  - `Domain`: Entidades, Value Objects, Enums, Errores de Dominio, Interfaces de Repositorio.
  - `Application`: Casos de uso / Handlers (CQRS/MediatR), DTOs, Validadores (FluentValidation), Mapeos.
  - `Infrastructure`: Persistencia (EF Core DbContext, migraciones), servicios externos, adaptadores.
  - `Api`: Controllers o Minimal APIs, configuración de DI, Middlewares, autenticación.

## 2. C# Moderno y Buenas Prácticas
- **Nullable Reference Types**: Habilitar `<Nullable>enable</Nullable>` y evitar desreferenciar nulls sin verificación.
- **Pattern Matching e Inmutabilidad**: Utilizar `records`, `readonly struct`, expresiones switch y constructores primarios (`primary constructors`).
- **Result Pattern**: Preferir retornar tipos `Result<T>` / `ErrorOr<T>` para errores esperados de negocio en lugar de lanzar excepciones para control de flujo.
- **Async/Await correcto**: Siempre propagar `CancellationToken` en operaciones I/O asíncronas y evitar `Task.Result` o `.Wait()` (anti-pattern de sync-over-async).

## 3. Entity Framework Core y Persistencia
- **Fluent API**: Configurar mappings de entidades en clases separadas (`IEntityTypeConfiguration<T>`) en lugar de Data Annotations en las entidades de dominio.
- **Consultas de solo lectura**: Utilizar `.AsNoTracking()` en queries que no requieran actualizar el estado de las entidades.
- **Índices y Migraciones**: Definir índices explícitos y generar migraciones reproducibles con `dotnet ef migrations add`.
