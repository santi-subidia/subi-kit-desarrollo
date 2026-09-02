---
name: tailwind-css
title: Estilos con Tailwind CSS y Diseño Responsivo
category: frontend
always_on: false
description: Reglas de maquetación, utility classes, tokens semánticos, diseño mobile-first y accesibilidad visual.
tags: [tailwind, css, frontend, styling, ui]
---

# Estilos con Tailwind CSS y Diseño UI

## 1. Filosofía Utility-First
- **Clases estándar**: Usar clases estándar de Tailwind en lugar de estilos inline o valores arbitrarios `[#123456]`, a menos que sea estrictamente necesario.
- **Mobile First**: Diseñar primero para dispositivos móviles y aplicar breakpoints ascendentes (`sm:`, `md:`, `lg:`, `xl:`).
- **Consistencia en Espaciado y Tipografía**: Utilizar la escala predeterminada de spacing (`p-2`, `p-4`, `gap-4`, `space-y-4`) y tamaños de texto coherentes (`text-sm`, `text-base`, `text-lg`).

## 2. Organización y Composición de Clases
- **Agrupación lógica**: Ordenar mentalmente las clases: Layout/Display (`flex`, `grid`) -> Posicionamiento (`relative`, `absolute`) -> Spacing/Dimensiones (`w-full`, `p-4`) -> Tipografía (`text-lg`, `font-semibold`) -> Colores/Efectos (`bg-primary`, `rounded-lg`, `shadow-sm`) -> Estados (`hover:`, `focus-visible:`).
- **Uso de `clsx` / `tailwind-merge` (`cn()`)**: Siempre que se combinen clases condicionales o dinámicas en componentes React, utilizar una función helper `cn(...)` para resolver conflictos.

## 3. Accesibilidad y Modo Oscuro
- **Contraste adecuado**: Garantizar que el texto cumpla con ratios WCAG AA de contraste sobre los fondos.
- **Focus Rings**: Incluir estilos visibles de foco (`focus-visible:ring-2 focus-visible:outline-none`) en elementos interactivos como botones y enlaces.
