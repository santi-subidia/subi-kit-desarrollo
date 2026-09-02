package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Paleta de colores refinada (Tokyo Night / Catppuccin Slate)
var (
	ColorPrimary    = lipgloss.Color("#7aa2f7") // Azul Tokyo Suave
	ColorSecondary  = lipgloss.Color("#bb9af7") // Lavanda Suave
	ColorSuccess    = lipgloss.Color("#9ece6a") // Verde Matcha
	ColorWarning    = lipgloss.Color("#e0af68") // Ámbar Cálido
	ColorDanger     = lipgloss.Color("#f7768e") // Coral Suave
	ColorCyan       = lipgloss.Color("#7dcfff") // Cian Suave
	ColorOrange     = lipgloss.Color("#ff9e64") // Durazno Cálido
	ColorText       = lipgloss.Color("#c0caf5") // Texto Claro Suave (Alto contraste y descanso visual)
	ColorTextBright = lipgloss.Color("#f7768e") // Resaltado
	ColorDim        = lipgloss.Color("#787c99") // Gris Pizarra
	ColorSubtle     = lipgloss.Color("#3b4261") // Bordes Inactivos
	ColorBgCard     = lipgloss.Color("#1f2335") // Fondo Tarjetas
	ColorBgSelected = lipgloss.Color("#292e42") // Fondo Selección
	ColorHighlight  = lipgloss.Color("#73daca") // Verde Menta / Código
)

// Estilos de la aplicación
var (
	// Header & Banner
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	VersionBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1a1b26")).
				Background(ColorSecondary).
				Padding(0, 1).
				MarginLeft(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true)

	// Pestañas (Tabs)
	TabStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Padding(0, 1).
			MarginRight(1)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1a1b26")).
			Background(ColorPrimary).
			Padding(0, 1).
			MarginRight(1)

	// Paneles & Contenedores
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtle).
			Padding(0, 1)

	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	FocusedPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSecondary).
				Padding(0, 1)

	// Listas y Elementos
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(ColorText)

	SelectedListItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				Background(ColorBgSelected).
				PaddingLeft(1).
				BorderLeft(true).
				BorderStyle(lipgloss.ThickBorder()).
				BorderForeground(ColorPrimary)

	ItemTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText)

	ItemDescStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	// Badges y Etiquetas
	BadgeCoreStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1a1b26")).
			Background(ColorPrimary).
			Padding(0, 1)

	BadgeSuccessStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1a1b26")).
				Background(ColorSuccess).
				Padding(0, 1)

	BadgeWarnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1a1b26")).
			Background(ColorWarning).
			Padding(0, 1)

	BadgeAgentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1a1b26")).
			Background(ColorSecondary).
			Padding(0, 1)

	BadgeCommandStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1b26")).
		Background(ColorCyan).
		Padding(0, 1)

	BadgeUpdateStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1b26")).
		Background(ColorWarning).
		Padding(0, 1)

	BadgeModeListStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1a1b26")).
				Background(ColorPrimary).
				Padding(0, 1)

	BadgeModeDetailStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1a1b26")).
				Background(ColorSecondary).
				Padding(0, 1)

	// Sección de detalles
	DetailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				MarginBottom(1)

	DetailSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSecondary).
				MarginTop(1).
				MarginBottom(0)

	CodeBlockStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Background(ColorBgCard).
			Padding(0, 1).
			MarginTop(0).
			MarginBottom(0)

	DocTextStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	// Barra de estado y ayuda (Footer)
	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(ColorBgCard).
			Padding(0, 1)

	KeyBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	KeyBadgeVioletStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSecondary)

	// Modal de Ayuda
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2).
			Background(ColorBgCard)
)
