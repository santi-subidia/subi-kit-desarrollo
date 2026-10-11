---
name: react-nextjs
title: Arquitectura React, Next.js App Router y Resiliencia UI
category: frontend
always_on: false
description: Patrones de Next.js App Router, Server Components por defecto, Server Actions, los 4 estados de UI y resiliencia de componentes.
tags: [nextjs, react, frontend, web, server-components, resilience]
---

# Arquitectura React, Next.js App Router y Resiliencia UI

## 1. Server Components vs Client Components
- **Server Components por defecto**: Mantener los componentes en el servidor a menos que requieran interactividad (hooks como `useState`, `useActionState`, o listeners del DOM).
- **Mover `'use client'` a los extremos (hojas del árbol)**: Encapsular la interactividad en pequeños componentes hoja en lugar de marcar páginas o layouts enteros como client components.
- **Evitar fugas de código cliente**: No importar módulos que dependan de APIs del navegador (`window`, `document`, `localStorage`) en Server Components.

## 2. Mutaciones, Data Fetching y Server Actions
- **Server Actions tipadas**: Implementar mutaciones mediante Server Actions, validando entradas exhaustivamente con esquemas Zod y retornando resultados estructurados (`{ success: boolean, data?: T, error?: string }`).
- **Data Fetching directo**: En Server Components, consultar bases de datos o servicios directamente mediante `async / await` sin crear rutas API REST redundantes.
- **Revalidación precisa de caché**: Utilizar `revalidatePath` o `revalidateTag` inmediatamente tras mutaciones exitosas para mantener la UI reactiva y sincronizada.

## 3. Resiliencia de Componentes y los 4 Estados Mandatorios
Toda vista, tabla o sección que consuma datos asíncronos debe manejar de forma explícita los **4 estados de UI**:
1. **`Loading` (Zero-CLS Skeletons)**:
   - Utilizar `<Suspense fallback={<ComponentSkeleton />}>` o `loading.tsx`. El skeleton debe replicar con exactitud la altura y distribución del componente final para evitar saltos de layout (*Cumulative Layout Shift*).
2. **`Empty` (Sin datos)**:
   - Mostrar un estado vacío contextual con ilustración/ícono y un Call To Action (CTA) accionable para crear el primer registro.
3. **`Error` (Error Boundaries y Recuperación)**:
   - Implementar `error.tsx` y Error Boundaries locales que capturen fallos sin romper toda la pantalla, ofreciendo un botón de reintento (*Retry*).
4. **`Success` (Renderizado completo)**:
   - La interfaz con datos completos, protegida contra desbordes de texto (`truncate`, `line-clamp`, `min-w-0` en contenedores flex).

## 4. Rendimiento y Optimización de Medios
- **Next.js Image (`next/image`)**: Siempre utilizar el componente `Image` con `sizes` adecuados, `placeholder="blur"` o dimensiones fijas/aspect-ratio para optimización automática (WebP/AVIF) y cero CLS.
- **Next.js Font (`next/font`)**: Cargar fuentes mediante `next/font/google` o `next/font/local` para inyección de CSS optimizado con `display: 'swap'`.
- **Formularios Móviles**: Asegurar que los `<input>` tengan un tamaño de texto mínimo de `16px` (`text-base` en Tailwind) para evitar el zoom forzado en iOS Safari.
