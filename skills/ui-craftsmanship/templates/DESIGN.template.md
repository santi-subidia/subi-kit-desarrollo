---
name: [Nombre del Proyecto / Sistema de Diseño]
description: [Descripción conceptual de la estética y dirección de arte del proyecto]

# Tokens de Diseño Portables (Exportación en OKLCH)
# Este archivo es la fuente de verdad del sistema de diseño.

colors:
  # Anclas de Marca (Brand Anchors)
  brand-primary: "oklch(62% 0.19 255)"         # Color de acento principal / CTA
  brand-secondary: "oklch(75% 0.14 195)"       # Acento secundario / estados activos
  brand-foreground: "oklch(98% 0.01 255)"      # Texto contrastante sobre color primario

  # Superficies (Surfaces) - Modo Oscuro
  ground-dark: "oklch(14% 0.015 260)"          # Fondo principal de página
  surface-dark: "oklch(18% 0.02 260)"         # Paneles, tarjetas y modales
  surface-raised-dark: "oklch(22% 0.025 260)"  # Inputs, tooltips, dropdowns

  # Superficies (Surfaces) - Modo Claro
  ground-light: "oklch(98% 0.005 260)"         # Fondo de página claro
  surface-light: "oklch(100% 0 0)"             # Tarjetas y paneles claros
  surface-raised-light: "oklch(95% 0.008 260)" # Inputs y controles claros

  # Tipografía y Jerarquía de Texto (Neutros Tintados)
  text-primary-dark: "oklch(96% 0.005 260)"    # Titulares y texto principal en modo oscuro
  text-muted-dark: "oklch(70% 0.015 260)"      # Subtítulos, captions y metadatos
  text-faint-dark: "oklch(50% 0.02 260)"       # Texto deshabilitado o decorativo

  text-primary-light: "oklch(18% 0.02 260)"   # Texto principal en modo claro
  text-muted-light: "oklch(45% 0.025 260)"     # Texto secundario en modo claro
  text-faint-light: "oklch(65% 0.02 260)"      # Metadatos tenues en modo claro

  # Bordes y Divisores
  border-subtle: "oklch(30% 0.02 260 / 0.4)"   # Borde estándar en oscuro
  border-subtle-light: "oklch(88% 0.01 260)"   # Borde estándar en claro

  # Estados Semánticos
  success: "oklch(62% 0.17 145)"
  warning: "oklch(75% 0.16 75)"
  error: "oklch(58% 0.22 25)"
  info: "oklch(65% 0.15 240)"

typography:
  fonts:
    display: "'Plus Jakarta Sans', 'Inter', system-ui, sans-serif"
    body: "'Plus Jakarta Sans', system-ui, sans-serif"
    mono: "'JetBrains Mono', 'Fira Code', monospace"
  
  scale:
    "10": "0.625rem"     # Overlines, badges micro
    "12": "0.75rem"      # Metadatos, etiquetas secundarias, tags
    "14": "0.875rem"     # Texto secundario denso, tablas, controles UI
    "16": "1rem"         # Cuerpo de texto base (Piso mínimo móvil)
    "18": "1.125rem"     # Párrafo introductorio (lead), subtítulo de tarjeta
    "20": "1.25rem"      # Título de tarjeta, encabezado h4
    "24": "1.5rem"       # Encabezado h3
    "30": "1.875rem"     # Encabezado h2
    "36": "2.25rem"      # Encabezado h1
    "48": "3rem"         # Sub-display
    "60": "3.75rem"      # Display hero

  fluid:
    hero: "clamp(2.5rem, 5vw + 1rem, 4.5rem)"
    heading: "clamp(1.75rem, 3vw + 1rem, 2.75rem)"

rounded:
  none: "0px"
  xs: "2px"
  sm: "4px"
  md: "6px"
  lg: "8px"
  xl: "12px"
  "2xl": "16px"
  pill: "9999px"

spacing:
  xs: "0.25rem"   # 4px
  sm: "0.5rem"    # 8px
  md: "1rem"      # 16px
  lg: "1.5rem"    # 24px
  xl: "2rem"      # 32px
  "2xl": "3rem"   # 48px
  "3xl": "4rem"   # 64px
---

# Sistema de Diseño: {{PROJECT_NAME}}

## 1. Filosofía y Dirección de Arte
- **Estilo Visual**: [Ej. Editorial moderno, Técnico utilitario, Minimalismo cálido]
- **Tono y Voz**: [Ej. Preciso, sofisticado, directo, sin ornamentos innecesarios]
- **Uso de Sombras y Elevación**: [Ej. Elevación mediante sutil variación de fondo y bordes de 1px en lugar de sombras pesadas]

## 2. Reglas de Implementación
- **Tipografía**:
  - Párrafos de lectura delimitados entre `45ch` y `75ch` (`max-w-prose`).
  - Inputs móviles obligatoriamente con `text-base` (≥16px) para evitar auto-zoom en iOS Safari.
  - Tracking de titulares calibrado a `-0.025em`.
- **Paleta y Contraste**:
  - Texto secundario siempre tintado con el matiz (`oklch hue`) del fondo; nunca gris puro `#888` en superficies oscuras.
  - Cumplimiento estricto de ratios WCAG AA (≥4.5:1 para texto normal, ≥3:1 para titulares).
- **Superficies y Feedback**:
  - `::selection` estilizado con `bg-brand-primary/20`.
  - `:focus-visible` con anillo contrastante de 2px y offset.
  - Cero transiciones con rebote (`bounce`); usar desaceleración `cubic-bezier(0.16, 1, 0.3, 1)`.
