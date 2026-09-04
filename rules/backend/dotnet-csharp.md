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

## 2. C# Moderno, Rendimiento y Asincronía
- **Nullable Reference Types**: Habilitar `<Nullable>enable</Nullable>` y evitar desreferenciar nulls sin verificación.
- **Pattern Matching e Inmutabilidad**: Utilizar `records`, `readonly struct`, expresiones switch y constructores primarios (`primary constructors`).
- **Result Pattern**: Preferir retornar tipos `Result<T>` / `ErrorOr<T>` para errores esperados de negocio en lugar de lanzar excepciones para control de flujo.
- **Asincronía Segura (Zero Sync-over-Async)**:
  - NUNCA bloquear sobre tareas asíncronas con `.Result`, `.Wait()` o `.GetAwaiter().GetResult()` (causan starvation del ThreadPool y deadlocks).
  - NUNCA hacer múltiples `await` sobre la misma instancia de `ValueTask` o `ValueTask<T>` (provoca comportamiento indefinido y corrupción silenciosa de estado).
  - Siempre propagar `CancellationToken` en todas las operaciones I/O asíncronas.
- **Gestión de Memoria en Rutas Críticas (Hot Paths)**:
  - Preferir `ReadOnlySpan<char>` / `AsSpan()` sobre `Substring()` para parsing y slicing de texto sin asignaciones en el Heap.
  - Usar `ArrayPool<T>.Shared.Rent()` y `Return()` para buffers temporales de I/O en lugar de instanciar arrays repetitivos.
  - Evitar `stackalloc` dentro de bucles cerrados para prevenir desbordes de stack.

## 3. Entity Framework Core de Alto Rendimiento
- **Fluent API**: Configurar mappings de entidades en clases separadas (`IEntityTypeConfiguration<T>`) en lugar de Data Annotations en las entidades de dominio.
- **Predicados Sargables (Index Seek vs Scan)**:
  - Mantener la columna indexada limpia en un lado de la comparación: NUNCA envolver columnas en funciones (`CreatedAt.Year == y`, `ToLower(Name) == n`, `Price * 1.1 > x`, `name.Contains(term)`).
  - Reescribir a rangos semiabiertos (ej: `CreatedAt >= start && CreatedAt < end`) o búsquedas ancladas (`name.StartsWith(term)`).
- **Prevención de N+1 y Proyecciones Eficientes**:
  - Prohibido habilitar Lazy Loading o Proxies en entornos de servidor.
  - Proyectar directamente a DTOs con `.Select(x => new ...)` para traer únicamente las columnas necesarias.
- **Explosión Cartesiana en Colecciones**:
  - Al incluir dos o más colecciones relacionadas con `.Include()`, usar obligatoriamente `.AsSplitQuery()` para evitar la multiplicación geométrica de filas.
- **Hot Paths & Consultas Compiladas**:
  - Para consultas muy frecuentes sobre el mismo `DbContext`, utilizar `EF.CompileQuery` o `EF.CompileAsyncQuery` para eliminar el coste de traducción de LINQ en cada llamada.
- **Consultas de solo lectura**: Utilizar `.AsNoTracking()` en queries que no requieran actualizar el estado de las entidades.

## 4. Higiene de MSBuild y Gestión de Proyectos
- **Central Package Management (CPM)**: Centralizar versiones de dependencias NuGet en `Directory.Packages.props` con `<ManagePackageVersionsCentrally>true</ManagePackageVersionsCentrally>` en soluciones multi-proyecto.
- **Compilaciones Incrementales Limpias**:
  - Evitar `CopyToOutputDirectory="Always"`; usar `PreserveNewest` o `Never` para no forzar rebuilds continuos.
  - Utilizar tareas nativas de MSBuild (`<MakeDir>`, `<Copy>`, `<Delete>`) en lugar de invocar comandos de shell con `<Exec>`.
  - Proteger propiedades condicionales con comillas simples (ej. `Condition="'$(Configuration)' == 'Release'"`).

