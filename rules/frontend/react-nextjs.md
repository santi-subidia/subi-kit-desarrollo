---
name: react-nextjs
title: Arquitectura React y Next.js App Router
category: frontend
always_on: false
description: Patrones de Next.js App Router, Server Components por defecto, Server Actions, optimización y manejo de estado.
tags: [nextjs, react, frontend, web]
---

# Arquitectura React y Next.js App Router

## 1. Server Components vs Client Components
- **Server Components por defecto**: Mantener los componentes en el servidor a menos que requieran interactividad (hooks como `useState`, `useEffect`, event listeners del DOM).
- **Mover `'use client'` a las hojas**: Encapsular la interactividad en pequeños componentes hoja en lugar de marcar páginas enteras como client components.
- **Evitar leaks de código cliente**: No importar módulos que dependan del navegador en Server Components.

## 2. Mutaciones y Data Fetching
- **Server Actions**: Usar Server Actions para mutaciones (formularios, acciones POST/PUT), validando entradas con Zod y manejando errores controlados (`{ success: boolean, error?: string }`).
- **Data Fetching directo**: En Server Components, consultar bases de datos o servicios directamente de forma asíncrona (`async / await`) sin necesidad de endpoints API intermedios.
- **Revalidación de caché**: Usar `revalidatePath` o `revalidateTag` después de mutaciones para mantener la UI sincronizada.

## 3. Rendimiento y UX
- **Next.js Image (`next/image`)**: Siempre utilizar el componente `Image` para optimización automática de formatos (WebP/AVIF) y dimensiones adaptativas.
- **Next.js Font (`next/font`)**: Utilizar Google Fonts / Local Fonts integrados para evitar layout shifts.
- **Suspense y Loadings**: Utilizar `loading.tsx` y `<Suspense fallback={<Skeleton />}>` para streaming de componentes lentos.
