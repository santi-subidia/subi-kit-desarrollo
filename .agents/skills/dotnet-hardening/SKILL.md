---
name: dotnet-hardening
description: "Diagnóstico y optimización de rendimiento en .NET y C#, consultas lentas en Entity Framework Core (EF Core), asincronía segura (Zero Sync-over-Async), gestión de memoria de alto rendimiento (Span, ArrayPool), higiene de proyectos MSBuild y rigor de aserciones en testing. Usar cuando se reporte lentitud en APIs .NET, queries lentas en EF Core, timeouts en DbContext, deadlocks por Task.Result/Wait, saturación de ThreadPool o GC, problemas de build incremental o al auditar la calidad de tests en C#."
---

# .NET Hardening & High-Performance Engineering 🛡️⚡

Guía prescriptiva de optimización de rendimiento, diagnóstico de cuellos de botella y buenas prácticas de ingeniería en el ecosistema .NET, C#, Entity Framework Core y MSBuild.

> [!IMPORTANT]
> **Regla de Medición**: Nunca optimices a ciegas. Captura siempre la métrica base (tiempo de ejecución, SQL generado, consumo de memoria o allocations) antes de modificar el código, y valida la mejora cuantificable tras el cambio.

---

## 🔍 Protocolo de Diagnóstico Inicial

### 1. Capturar el SQL Real de EF Core
No asumas lo que EF Core está generando. Habilita el logging en desarrollo o pruebas:
```csharp
optionsBuilder.LogTo(Console.WriteLine, LogLevel.Information);
// O en appsettings.Development.json:
// "Logging": { "LogLevel": { "Microsoft.EntityFrameworkCore.Database.Command": "Information" } }
```
Etiqueta consultas complejas con `.TagWith("GetOrderDetails_Optimized")` para ubicarlas instantáneamente en los logs o profiler.

### 2. Verificar el Plan de Ejecución
- En SQL Server: Activar *Include Actual Execution Plan*.
- En PostgreSQL: Ejecutar `EXPLAIN ANALYZE <consulta_generada>`.
- **Objetivo**: Asegurar que las búsquedas utilicen **Index Seek** en lugar de **Clustered Index Scan** o **Seq Scan**.

---

## 🗄️ Módulo 1: Entity Framework Core de Alto Rendimiento

### 1.1 Predicados Sargables (Index Seek vs Full Table Scan)
Un índice únicamente puede aprovecharse si la columna aparece **limpia** en un lado de la comparación.

| Escenario | ❌ Anti-patrón (Non-Sargable - Full Scan) | ✅ Solución Óptima (Sargable - Index Seek) |
| :--- | :--- | :--- |
| **Filtro por año / fecha** | `db.Orders.Where(o => o.CreatedAt.Year == 2026)` | `var start = new DateTime(2026, 1, 1);`<br>`var end = start.AddYears(1);`<br>`db.Orders.Where(o => o.CreatedAt >= start && o.CreatedAt < end)` |
| **Filtro insensible a mayúsculas** | `db.Users.Where(u => u.Email.ToLower() == email)` | Normalizar al persistir o usar intercalación (*collation*) insensible a mayúsculas sin funciones en LINQ. |
| **Búsqueda de texto** | `db.Products.Where(p => p.Sku.Contains(term))` *(Genera `%term%`, scan completo)* | `db.Products.Where(p => p.Sku.StartsWith(term))` *(Genera `term%`, permite Seek)* o Full-Text Search. |
| **Conversión de tipos** | `db.Logs.Where(l => l.StatusCode.ToString().StartsWith("5"))` | Filtrar sobre la columna tipada: `l.StatusCode >= 500 && l.StatusCode < 600`. |

**Verificación**: El plan de ejecución pasa de *Index Scan* a *Index Seek* y la duración cae drásticamente.

### 1.2 Explosión Cartesiana & Split Queries
Incluir dos o más colecciones relacionadas en una sola consulta multiplica las filas exponencialmente y satura la red y memoria del servidor.

```csharp
// ❌ BAD: Genera una explosión cartesiana masiva (Orders x Items x Shipments)
var order = await db.Orders
    .Include(o => o.Items)
    .Include(o => o.Shipments)
    .FirstOrDefaultAsync(o => o.Id == id);

// ✅ GOOD: Cada colección se carga en una consulta SQL limpia e independiente
var order = await db.Orders
    .Include(o => o.Items)
    .Include(o => o.Shipments)
    .AsSplitQuery()
    .FirstOrDefaultAsync(o => o.Id == id);
```
**Verificación**: Verificar en el log de EF que se emitan consultas separadas y que el volumen total de bytes transferidos disminuya radicalmente.

### 1.3 Erradicación de N+1 y Proyecciones DTO
- **Prohibido**: Marcar navigations como `virtual` o instalar proxies de Lazy Loading en aplicaciones web/APIs.
- **Regla**: Proyectar directamente con `.Select()` para traer únicamente las columnas requeridas por el caso de uso:

```csharp
// ✅ GOOD: Trae solo las 3 columnas necesarias en 1 solo viaje a la base de datos
var summary = await db.Users
    .Where(u => u.TenantId == tenantId)
    .Select(u => new UserSummaryDto(u.Id, u.FullName, u.Email))
    .ToListAsync(cancellationToken);
```

### 1.4 Consultas Compiladas en Hot Paths
Para endpoints o procesos que ejecutan la misma forma de consulta miles de veces por segundo, compila el delegado una sola vez para eliminar el coste de traducción de LINQ:

```csharp
private static readonly Func<AppDbContext, int, Task<UserDto?>> GetUserByIdQuery =
    EF.CompileAsyncQuery((AppDbContext db, int id) =>
        db.Users
          .Where(u => u.Id == id)
          .Select(u => new UserDto(u.Id, u.Name))
          .FirstOrDefault());

public Task<UserDto?> GetUserAsync(int id) => GetUserByIdQuery(_dbContext, id);
```

---

## ⚡ Módulo 2: Rendimiento C#, Asincronía y Gestión de Memoria

### 2.1 Cero Sync-over-Async (Prevención de Deadlocks)
NUNCA bloquees tareas asíncronas de forma sincrónica. Destruye la escalabilidad del ThreadPool y genera deadlocks cuando hay contextos de sincronización.

| ❌ Prohibido (Sync-over-Async) | ✅ Obligatorio (Asincronía Pura) |
| :--- | :--- |
| `var res = GetDataAsync().Result;` | `var res = await GetDataAsync(ct);` |
| `GetDataAsync().Wait();` | `await GetDataAsync(ct);` |
| `var res = GetDataAsync().GetAwaiter().GetResult();` | `var res = await GetDataAsync(ct);` |

### 2.2 Higiene Estricta de `ValueTask`
`ValueTask` y `ValueTask<T>` están optimizados para asignación cero cuando la operación completa sincrónicamente. **Tienen restricciones estrictas**:
- **NUNCA** hagas múltiples `await` sobre la misma instancia de `ValueTask`. Provoca comportamiento indefinido y corrupción silenciosa.
- **NUNCA** llames a `.Result` o `.GetAwaiter().GetResult()` antes de que la operación haya completado.
- Si necesitas almacenar o compartir el resultado pendiente entre múltiples consumidores, conviértelo inmediatamente a `Task`:
  ```csharp
  Task<int> task = SomeValueTaskOperation().AsTask();
  ```

### 2.3 Cero Asignaciones en Rutas Críticas (Hot Paths)

#### Parsing de Cadenas con `ReadOnlySpan<char>`
Evita `.Substring()` en parsing intensivo de texto; genera nuevas instancias en el Heap en cada llamada.
```csharp
// ❌ BAD: Asigna un nuevo string en el Heap
string code = payload.Substring(0, 4);

// ✅ GOOD: Cero asignaciones en memoria (slicing sobre el buffer existente)
ReadOnlySpan<char> code = payload.AsSpan(0, 4);
```

#### Buffers Temporales con `ArrayPool<T>`
Para operaciones de red, compresión o I/O con arrays de bytes temporales:
```csharp
byte[] buffer = ArrayPool<byte>.Shared.Rent(4096);
try
{
    int bytesRead = await stream.ReadAsync(buffer.AsMemory(0, 4096), ct);
    Process(buffer.AsSpan(0, bytesRead));
}
finally
{
    ArrayPool<byte>.Shared.Return(buffer);
}
```

#### Prevención de `stackalloc` en Bucles
- **Peligro**: `stackalloc` no se libera al salir del bloque del bucle, sino al salir de la función entera.
- **Regla**: Prohibido usar `stackalloc` dentro de bucles `for`/`while`/`foreach`. Si necesitas memoria en stack en un bucle, asígnala antes del bucle o usa un buffer de `ArrayPool<T>`.

---

## 🛠️ Módulo 3: Higiene de MSBuild y Compilaciones Incrementales

### 3.1 Eliminación de `CopyToOutputDirectory="Always"`
- **Problema**: `CopyToOutputDirectory="Always"` invalida la detección de cambios de MSBuild y fuerza un rebuild completo de todos los proyectos dependientes.
- **Solución**: Usar siempre `PreserveNewest` o `Never`:
  ```xml
  <!-- ❌ BAD -->
  <None Update="config.json" CopyToOutputDirectory="Always" />

  <!-- ✅ GOOD -->
  <None Update="config.json" CopyToOutputDirectory="PreserveNewest" />
  ```

### 3.2 Tareas Nativas vs Comandos `<Exec>`
No uses llamadas a la terminal (`mkdir`, `copy`, `del`) que rompen la portabilidad entre Windows y Linux y ciegan a MSBuild:
```xml
<!-- ❌ BAD -->
<Target Name="CustomStep">
  <Exec Command="mkdir $(OutputPath)logs" />
  <Exec Command="copy file.txt $(OutputPath)" />
</Target>

<!-- ✅ GOOD -->
<Target Name="CustomStep">
  <MakeDir Directories="$(OutputPath)logs" />
  <Copy SourceFiles="file.txt" DestinationFolder="$(OutputPath)" />
</Target>
```

### 3.3 Central Package Management (CPM)
En soluciones con 2 o más proyectos, no disperses versiones de NuGet en cada `.csproj`. Centralízalas en `Directory.Packages.props`:
```xml
<Project>
  <PropertyGroup>
    <ManagePackageVersionsCentrally>true</ManagePackageVersionsCentrally>
  </PropertyGroup>
  <ItemGroup>
    <PackageVersion Include="Microsoft.EntityFrameworkCore" Version="9.0.2" />
    <PackageVersion Include="FluentValidation" Version="11.11.0" />
  </ItemGroup>
</Project>
```
En el `.csproj` solo se declara `<PackageReference Include="FluentValidation" />` sin especificar `Version`.

---

## 🧪 Módulo 4: Rigor y Diversidad de Aserciones en Testing

### 4.1 Erradicar Falsos Positivos de Cobertura
Un test que solo verifica que el resultado "no es nulo" da una falsa sensación de seguridad mientras los bugs pasan desapercibidos en producción.

```csharp
// ❌ BAD: Aserción trivial y superficial (falsa confianza)
[TestMethod]
public async Task CreateOrder_ShouldSucceed()
{
    var order = await _service.CreateOrderAsync(request);
    Assert.IsNotNull(order); // Pasa incluso si el total es $0 o faltan items
}

// ✅ GOOD: Aserción profunda de invariantes, estado y reglas de negocio
[TestMethod]
public async Task CreateOrder_ShouldCalculateTotalsAndApplyDiscount()
{
    var order = await _service.CreateOrderAsync(request);

    Assert.IsNotNull(order);
    Assert.AreEqual(OrderStatus.Pending, order.Status);
    Assert.AreEqual(90.00m, order.FinalPrice); // 100 - 10%
    Assert.AreEqual(2, order.Items.Count);
    Assert.IsTrue(order.DiscountApplied);
}
```

### 4.2 Matriz de Validación de Tests
Toda suite de tests rigurosa debe cubrir:
1. **Valores concretos**: Comprobar los datos esperados, no solo su existencia.
2. **Casos de error esperados**: Verificar excepciones de dominio con `Assert.ThrowsExceptionAsync<DomainException>(...)` o fluent assertions equivalentes.
3. **Efectos secundarios e invariantes**: Verificar que las entidades modificadas hayan mutado su estado en la persistencia.
4. **Pruebas de límites**: Rango mínimo, nulos, colecciones vacías y duplicados.
