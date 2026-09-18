---
name: fast-fix
description: >-
  Flujo ágil y científico para resolver bugs, regresiones y errores puntuales sin overhead ceremonial de SDD.
  Enfocado en: Reproducción con Test (ROJO) -> Fix Quirúrgico -> Verificación (VERDE) -> Commit.
---

# Fast-Fix: Protocolo Ágil de Corrección de Bugs ⚡

Protocolo de alta velocidad para corregir bugs, fallos en runtime o regresiones en 1 o 2 archivos sin la sobrecarga documental de Spec-Driven Development (SDD).

> [!IMPORTANT]
> **Cero Overhead Documental**: En el flujo Fast-Fix **no** se crean `spec.md`, `tech-plan.md` ni `tasks.md`. El contrato formal de aceptación es el test de reproducción automatizado.

---

## 🔁 Las 4 Etapas del Flujo Fast-Fix

```
  ┌────────────────────────────────────────────────────────┐
  │ 1. REPRODUCCIÓN (ROJO)                                 │
  │    Escribir test unitario o comando que falle (<2s)    │
  └──────────────────────────┬─────────────────────────────┘
                             │
                             ▼
  ┌────────────────────────────────────────────────────────┐
  │ 2. FIX QUIRÚRGICO                                      │
  │    Modificar exclusivamente el código indispensable    │
  └──────────────────────────┬─────────────────────────────┘
                             │
                             ▼
  ┌────────────────────────────────────────────────────────┐
  │ 3. DEMOSTRACIÓN (VERDE) & SUITE                        │
  │    El test pasa a VERDE + suite general sin romper     │
  └──────────────────────────┬─────────────────────────────┘
                             │
                             ▼
  ┌────────────────────────────────────────────────────────┐
  │ 4. CIERRE Y COMMIT ATÓMICO                             │
  │    Preservar test de regresión + commit semántico fix  │
  └────────────────────────────────────────────────────────┘
```

---

## 🛠️ Detalle de Ejecución por Etapa

### Etapa 1: Reproducción Determinística (ROJO) 🔴
Antes de editar cualquier línea del código fuente de producción:
1. **Identificar el Seam**: Localizar la función, método o endpoint donde se manifiesta el bug.
2. **Construir el Test de Regresión**:
   - **Go**: Crear función `TestBugName(t *testing.T)` en el archivo `_test.go` correspondiente.
   - **TypeScript / React**: Crear test en Vitest / Jest (`it('should handle edge case...', () => { ... })`).
   - **C# / .NET**: Crear test en xUnit/NUnit (`[Fact] public void Should_HandleEdgeCase()`).
3. **Ejecutar y Validar Fallo**:
   - El test debe fallar con el error exacto reportado por el usuario o detectado en runtime.
   - **Prohibido avanzar** si el test pasa en verde antes de aplicar el fix.

### Etapa 2: Fix Quirúrgico 🎯
1. Localizar la causa raíz exacta (ej. valor nulo no contemplado, condición de borde invertida, parseo erróneo de JSON, tipo incorrecto).
2. Aplicar el cambio mínimo necesario:
   - Mantener nombres y estilo preexistente.
   - No introducir refactorizaciones fuera de alcance (*no drive-by refactoring*).
   - Preservar comentarios y estructura circundante.

### Etapa 3: Demostración en Verde & Suite Completa 🟢
1. **Reejecutar el Test de Regresión**: Debe pasar de inmediato a **VERDE**.
2. **Ejecutar la Suite del Módulo o Proyecto**:
   ```bash
   # Go
   go test ./...
   # Node / TS
   npm test
   # .NET
   dotnet test
   ```
   Garantizar que ninguna prueba preexistente haya fallado.

### Etapa 4: Commit Semántico de Regresión 📝
1. El test de reproducción se conserva permanentemente en el repositorio como salvaguarda contra futuras regresiones.
2. Realizar commit atómico siguiendo la convención:
   ```
   fix(<módulo>): <descripción imperativa del arreglo>
   ```

---

## ⏱️ Presupuesto de Intentos (Attempt Budget) 🎯

Para evitar bucles infinitos de prueba y error en el flujo ágil:

1. **Límite Estricto de 3 Intentos en Etapa 2 (Fix Quirúrgico)**:
   - El desarrollador tiene hasta 3 oportunidades para que el test de regresión pase a VERDE.
2. **Parada Obligatoria y Rendición de Evidencia**:
   - Si al 3er intento el test continúa fallando, **queda prohibido realizar un 4to intento a ciegas**.
   - Reportar inmediatamente al Orquestador y a Subi:
     - Comando del test ejecutado y su código de salida.
     - Logs limpios del error.
     - Resumen de las hipótesis descartadas.
3. **Escalamiento a SDD o Intervención Humana**:
   - Si el bug no se resuelve en 3 intentos, significa que no es un fix superficial: se escala a una sesión de diagnóstico profundo (`diagnosing-bugs`) o se abre una spec formal en SDD.

---

## 🚫 Anti-Patrones Prohibidos en Fast-Fix
- **Modificar código sin test previo**: Tocar el código "a ojo" o adivinar el fix.
- **Superar los 3 intentos**: Caer en el bucle de prueba y error sin detenerse a analizar.
- **Crear documentos innecesarios**: Redactar `spec.md` para un bug de 5 líneas de código.
- **Eliminar el test tras el fix**: El test de reproducción debe quedarse en el repositorio para siempre.
- **Refactors masivos**: Cambiar la arquitectura o renombrar métodos no relacionados aprovechando el bugfix.

