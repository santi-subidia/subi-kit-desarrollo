---
name: tailwind-css
title: Estilos con Tailwind CSS, Tokens y Diseño Responsivo
category: frontend
always_on: false
description: Reglas de maquetación, tokens semánticos (OKLCH), diseño mobile-first, prevención de AI slop y accesibilidad visual.
tags: [tailwind, css, frontend, styling, ui, design-system, oklch]
---

# Estilos con Tailwind CSS y Diseño UI Artesanal

## 1. Filosofía Token-First y Sistema de Diseño
- **Fuente de Verdad (`DESIGN.md` / `@theme`)**: Utilizar siempre los tokens semánticos del sistema de diseño (`bg-background`, `text-foreground`, `bg-primary`, `border-border`, etc.). Evitar valores hexadecimales arbitrarios hardcodeados (`[#123456]`).
- **Espacio de Color OKLCH y Neutros Tintados**: Los fondos oscuros y textos secundarios deben emplear neutros tintados con el matiz (`hue`) de la marca. Prohibido el texto gris plano sobre fondos saturados o con tinte.
- **Mobile-First Real**: Construir primero para pantallas pequeñas y escalar progresivamente con breakpoints (`sm:`, `md:`, `lg:`, `xl:`).
- **Piso Mínimo en Inputs Móviles**: Los campos de texto en móvil deben tener `text-base` (≥16px / 1rem) para evitar que iOS Safari fuerce un zoom destructivo en la pantalla al enfocar.

## 2. Prevención de "AI Frontend Slop" (Anti-Patrones Prohibidos)
- **🚫 Prohibido anidar tarjetas dentro de tarjetas (*Cards in cards*)**: Usar variaciones de fondo sutiles, bordes tenues o espaciado en lugar de cajas dentro de cajas.
- **🚫 Prohibidas las cuadrículas monótonas de tarjetas idénticas**: Romper la monotonía combinando elementos de ancho completo, asimetría controlada y jerarquías claras.
- **🚫 Prohibido el *kicker/eyebrow* compulsivo**: No colocar etiquetas en mayúsculas flotando arriba de cada título en cada sección. El título debe sostenerse solo.
- **🚫 Prohibido el texto con gradientes decorativos**: El énfasis se expresa mediante peso tipográfico (`font-bold`), tamaño o contraste.
- **🚫 Prohibidas las sombras duras de bloque**: Evitar `box-shadow: 4px 4px 0` salvo en interfaces expresamente neobrutalistas.

## 3. Organización, Composición y Helper `cn()`
- **Agrupación lógica de clases**:
  `Layout/Display` (`flex`, `grid`) → `Posicionamiento` (`relative`, `absolute`) → `Spacing/Dimensiones` (`w-full`, `p-4`, `min-w-0`) → `Tipografía` (`text-base`, `font-semibold`) → `Colores/Fondos` (`bg-surface`, `text-foreground`) → `Bordes/Sombras` (`rounded-lg`, `border`) → `Estados e Interacciones` (`hover:`, `focus-visible:`, `disabled:`).
- **Uso de `clsx` + `tailwind-merge` (`cn()`)**: En componentes React dinámicos o reutilizables, resolver colisiones de clases siempre con la función helper `cn(...)`.

## 4. Accesibilidad (a11y), Superficies y Ergonomía
- **Contraste WCAG AA**: Ratio mínimo de 4.5:1 en texto regular y 3:1 en encabezados grandes.
- **Touch Targets Ergonómicos**: En dispositivos táctiles, botones y enlaces deben tener al menos **44x44px** de área interactiva (`min-h-[44px] min-w-[44px]` o `p-3`).
- **Anillos de Foco Visibles**: Todo elemento interactivo debe incluir `:focus-visible` accesible (`focus-visible:ring-2 focus-visible:ring-primary focus-visible:outline-none focus-visible:ring-offset-2`).
- **Superficies Nativas Personalizadas**:
  - Selección de texto: `selection:bg-primary/20 selection:text-primary-foreground`.
  - Números tabulares en tablas/métricas: `tabular-nums`.
  - Reducción de movimiento: `motion-reduce:transition-none` o `motion-reduce:animate-none`.
