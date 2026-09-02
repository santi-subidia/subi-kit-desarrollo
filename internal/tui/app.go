package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/agents"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/mcp"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/skills"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/updater"
)

// AppModel es el modelo central de la interfaz interactiva con Bubble Tea.
type AppModel struct {
	activeTab   Tab
	focusPane   FocusPane
	cursor      int
	viewport    viewport.Model
	ready       bool
	showHelp    bool
	width       int
	height      int
	version     string

	// Estado de auto-actualización
	updateResult *updater.UpdateResult
	updateStatus string
	isUpdating   bool

	// Gestores de datos
	rulesMgr  *rules.Manager
	skillsMgr *skills.Manager
	mcpMgr    *mcp.Manager
	agentsMgr *agents.Manager

	// Catálogos cargados
	commands []CommandDoc
	glossary []GlossaryTerm
	sddFeats []SDDFeatureStatus
}

// NewAppModel inicializa el modelo de la TUI con los catálogos y gestores.
func NewAppModel(rulesMgr *rules.Manager, skillsMgr *skills.Manager, mcpMgr *mcp.Manager, agentsMgr *agents.Manager, version string) AppModel {
	return AppModel{
		activeTab: TabCommands,
		focusPane: FocusList,
		cursor:    0,
		version:   version,
		rulesMgr:  rulesMgr,
		skillsMgr: skillsMgr,
		mcpMgr:    mcpMgr,
		agentsMgr: agentsMgr,
		commands:  GetCommandsData(),
		glossary:  GetGlossaryTerms(),
		sddFeats:  GetSDDFeatures(),
	}
}

type updateCheckMsg struct {
	result *updater.UpdateResult
	err    error
}

type updateApplyMsg struct {
	status string
	err    error
}

func checkUpdateCmd(version string) tea.Cmd {
	return func() tea.Msg {
		res, err := updater.CheckLatest(version)
		return updateCheckMsg{result: res, err: err}
	}
}

func applyUpdateCmd(release *updater.ReleaseInfo) tea.Cmd {
	return func() tea.Msg {
		err := updater.ApplyUpdate(release, nil)
		if err != nil {
			return updateApplyMsg{err: err}
		}
		return updateApplyMsg{status: "¡SubiKit actualizado con éxito!"}
	}
}

// Init inicializa el programa Bubble Tea con chequeo de actualizaciones en segundo plano.
func (m AppModel) Init() tea.Cmd {
	return checkUpdateCmd(m.version)
}

// Update maneja los mensajes y eventos de teclado con navegación estilo Neovim.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Teclas globales
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "?", "f1":
			m.showHelp = !m.showHelp
			return m, nil

		case "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			if m.focusPane == FocusDetail {
				m.focusPane = FocusList
				return m, nil
			}

		case "tab":
			// Alternar foco entre panel de lista y panel de lectura
			if m.focusPane == FocusList {
				m.focusPane = FocusDetail
			} else {
				m.focusPane = FocusList
			}
			return m, nil

		case "u", "U":
			if m.updateResult != nil && m.updateResult.UpdateAvail && !m.isUpdating && m.updateResult.Release != nil {
				m.isUpdating = true
				m.updateStatus = "Descargando e instalando actualización..."
				return m, applyUpdateCmd(m.updateResult.Release)
			}
			return m, nil

		case "1", "2", "3", "4", "5", "6", "7", "8":
			idx := int(msg.Runes[0] - '1')
			if idx >= 0 && idx < len(TabNames) {
				m.activeTab = Tab(idx)
				m.focusPane = FocusList
				m.cursor = 0
				m.updateViewportContent()
			}
			return m, nil
		}

		// Si el modal de ayuda está abierto
		if m.showHelp {
			if msg.String() == "q" || msg.String() == "esc" || msg.String() == "?" {
				m.showHelp = false
			}
			return m, nil
		}

		// Navegación según el panel con foco
		if m.focusPane == FocusList {
			// === MODO LISTA ===
			switch msg.String() {
			case "q":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.updateViewportContent()
				}

			case "down", "j":
				maxItems := m.getItemCount()
				if m.cursor < maxItems-1 {
					m.cursor++
					m.updateViewportContent()
				}

			case "g", "home":
				m.cursor = 0
				m.updateViewportContent()

			case "G", "end":
				maxItems := m.getItemCount()
				if maxItems > 0 {
					m.cursor = maxItems - 1
					m.updateViewportContent()
				}

			case "l", "right", "enter":
				// Entrar al modo lectura en el panel de detalles
				m.focusPane = FocusDetail

			case "h", "left":
				// Retroceder pestaña si estamos en la lista
				if m.activeTab == 0 {
					m.activeTab = Tab(len(TabNames) - 1)
				} else {
					m.activeTab--
				}
				m.cursor = 0
				m.updateViewportContent()
			}
		} else {
			// === MODO LECTURA DE DETALLES (ESTILO NEOVIM) ===
			switch msg.String() {
			case "q", "h", "left":
				// Volver al panel de lista
				m.focusPane = FocusList
				return m, nil

			case "j", "down":
				// Bajar línea por línea exacta
				m.viewport.LineDown(1)
				return m, nil

			case "k", "up":
				// Subir línea por línea exacta
				m.viewport.LineUp(1)
				return m, nil

			case "d", "ctrl+d":
				// Media página abajo
				m.viewport.HalfViewDown()
				return m, nil

			case "u", "ctrl+u":
				// Media página arriba
				m.viewport.HalfViewUp()
				return m, nil

			case " ", "pgdown":
				// Página completa abajo
				m.viewport.ViewDown()
				return m, nil

			case "pgup":
				// Página completa arriba
				m.viewport.ViewUp()
				return m, nil

			case "g", "home":
				// Ir al principio
				m.viewport.GotoTop()
				return m, nil

			case "G", "end":
				// Ir al final
				m.viewport.GotoBottom()
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 4
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight

		detailWidth := m.getRightPaneWidth()
		detailHeight := m.height - verticalMarginHeight - 2

		if detailHeight < 5 {
			detailHeight = 5
		}

		if !m.ready {
			m.viewport = viewport.New(detailWidth, detailHeight)
			m.ready = true
		} else {
			m.viewport.Width = detailWidth
			m.viewport.Height = detailHeight
		}

		m.updateViewportContent()

	case updateCheckMsg:
		if msg.err == nil && msg.result != nil {
			m.updateResult = msg.result
			if m.activeTab == TabDoctor {
				m.updateViewportContent()
			}
		}
		return m, nil

	case updateApplyMsg:
		m.isUpdating = false
		if msg.err != nil {
			m.updateStatus = fmt.Sprintf("Error al actualizar: %v", msg.err)
		} else {
			m.updateStatus = "¡SubiKit actualizado! Reinicia la app."
		}
		if m.activeTab == TabDoctor {
			m.updateViewportContent()
		}
		return m, nil
	}

	// Actualizar eventos adicionales del viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *AppModel) getItemCount() int {
	switch m.activeTab {
	case TabCommands:
		return len(m.commands)
	case TabAgents:
		return len(m.agentsMgr.GetAll())
	case TabRules:
		return len(m.rulesMgr.GetAll())
	case TabSkills:
		return len(m.skillsMgr.GetAll())
	case TabMCP:
		return len(m.mcpMgr.GetCatalog())
	case TabDoctor:
		return 1
	case TabSDD:
		if len(m.sddFeats) == 0 {
			return 1
		}
		return len(m.sddFeats)
	case TabGlossary:
		return len(m.glossary)
	default:
		return 0
	}
}

func (m *AppModel) getLeftPaneWidth() int {
	leftWidth := m.width / 3
	if leftWidth < 28 {
		leftWidth = 28
	}
	if leftWidth > 38 {
		leftWidth = 38
	}
	return leftWidth
}

func (m *AppModel) getRightPaneWidth() int {
	leftWidth := m.getLeftPaneWidth()
	rightWidth := m.width - leftWidth - 5
	if rightWidth < 30 {
		rightWidth = 30
	}
	return rightWidth
}

func (m *AppModel) updateViewportContent() {
	if !m.ready {
		return
	}

	leftWidth := m.getLeftPaneWidth()
	rightWidth := m.getRightPaneWidth()
	var detailContent string

	switch m.activeTab {
	case TabCommands:
		_, detailContent = RenderCommandsView(m.commands, m.cursor, leftWidth, rightWidth)
	case TabAgents:
		_, detailContent = RenderAgentsView(m.agentsMgr.GetAll(), m.cursor, leftWidth, rightWidth)
	case TabRules:
		_, detailContent = RenderRulesView(m.rulesMgr.GetAll(), m.cursor, leftWidth, rightWidth)
	case TabSkills:
		_, detailContent = RenderSkillsView(m.skillsMgr.GetAll(), m.cursor, leftWidth, rightWidth)
	case TabMCP:
		_, detailContent = RenderMCPView(m.mcpMgr, m.cursor, leftWidth, rightWidth)
	case TabDoctor:
		detailContent = RenderDoctorView(m.rulesMgr, m.skillsMgr, m.mcpMgr, m.agentsMgr, m.updateResult, rightWidth)
	case TabSDD:
		m.sddFeats = GetSDDFeatures()
		_, detailContent = RenderSDDView(m.sddFeats, m.cursor, leftWidth, rightWidth)
	case TabGlossary:
		_, detailContent = RenderGlossaryView(m.glossary, m.cursor, leftWidth, rightWidth)
	}

	m.viewport.SetContent(detailContent)
	m.viewport.GotoTop()
}

// View renderiza el frame completo de la TUI.
func (m AppModel) View() string {
	if !m.ready {
		return "\n  Iniciando SubiKit TUI..."
	}

	if m.showHelp {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			RenderHelpModal(),
		)
	}

	var sb strings.Builder

	// 1. Header & Title Bar con Badge de Modo Activo y Notificación de Actualización
	headerTitle := HeaderStyle.Render("⚡ SUBIKIT TUI")
	versionBadge := VersionBadgeStyle.Render("v" + m.version)

	var updateBadge string
	if m.isUpdating {
		updateBadge = " " + BadgeUpdateStyle.Render("⏳ "+m.updateStatus)
	} else if m.updateResult != nil && m.updateResult.UpdateAvail {
		updateBadge = " " + BadgeUpdateStyle.Render(fmt.Sprintf("⚡ v%s disponible [u: Actualizar]", m.updateResult.LatestVersion))
	} else if m.updateStatus != "" {
		updateBadge = " " + BadgeSuccessStyle.Render(m.updateStatus)
	}

	var modeBadge string
	if m.focusPane == FocusList {
		modeBadge = BadgeModeListStyle.Render("MODO LISTA • [l/Enter: Enfocar Texto]")
	} else {
		modeBadge = BadgeModeDetailStyle.Render("MODO LECTURA NEOVIM • [j/k: línea • d/u: 1/2 pág • h/Esc: volver]")
	}

	helpHint := lipgloss.NewStyle().Foreground(ColorDim).Render("[?: Ayuda / Atajos] [Tab: Foco] [q: Salir]")
	headerLine := lipgloss.JoinHorizontal(lipgloss.Center, headerTitle, versionBadge, updateBadge, "  ", modeBadge)
	headerRow := lipgloss.JoinHorizontal(lipgloss.Left, headerLine, "   ", helpHint)
	sb.WriteString(headerRow + "\n")

	// 2. Navigation Tabs
	var tabs []string
	for i, name := range TabNames {
		if Tab(i) == m.activeTab {
			tabs = append(tabs, ActiveTabStyle.Render(name))
		} else {
			tabs = append(tabs, TabStyle.Render(name))
		}
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...) + "\n\n")

	// 3. Main Content (Split Panes con resalte del panel enfocado)
	leftWidth := m.getLeftPaneWidth()
	rightWidth := m.getRightPaneWidth()
	paneHeight := m.height - 7
	if paneHeight < 5 {
		paneHeight = 5
	}

	var leftContent string
	switch m.activeTab {
	case TabCommands:
		leftContent, _ = RenderCommandsView(m.commands, m.cursor, leftWidth, rightWidth)
	case TabAgents:
		leftContent, _ = RenderAgentsView(m.agentsMgr.GetAll(), m.cursor, leftWidth, rightWidth)
	case TabRules:
		leftContent, _ = RenderRulesView(m.rulesMgr.GetAll(), m.cursor, leftWidth, rightWidth)
	case TabSkills:
		leftContent, _ = RenderSkillsView(m.skillsMgr.GetAll(), m.cursor, leftWidth, rightWidth)
	case TabMCP:
		leftContent, _ = RenderMCPView(m.mcpMgr, m.cursor, leftWidth, rightWidth)
	case TabDoctor:
		doctorStatus := "🩺 Diagnóstico en vivo\n  Escaneo en tiempo real"
		if m.updateResult != nil && m.updateResult.UpdateAvail {
			doctorStatus += "\n\n⚡ Actualización:\n  " + m.updateResult.LatestVersion + " lista"
		}
		leftContent = SelectedListItemStyle.Width(leftWidth - 4).Render(doctorStatus)
	case TabSDD:
		leftContent, _ = RenderSDDView(m.sddFeats, m.cursor, leftWidth, rightWidth)
	case TabGlossary:
		leftContent, _ = RenderGlossaryView(m.glossary, m.cursor, leftWidth, rightWidth)
	}

	// Aplicar estilo de borde activo según el foco
	leftStyle := PanelStyle
	rightStyle := PanelStyle

	if m.focusPane == FocusList {
		leftStyle = ActivePanelStyle
	} else {
		rightStyle = FocusedPanelStyle
	}

	leftPane := leftStyle.
		Width(leftWidth).
		Height(paneHeight).
		Render(leftContent)

	rightPane := rightStyle.
		Width(rightWidth).
		Height(paneHeight).
		Render(m.viewport.View())

	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, " ", rightPane)
	sb.WriteString(mainLayout + "\n")

	// 4. Footer & Shortcut Bar
	scrollPercent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	var footerText string

	updateShortcut := ""
	if m.updateResult != nil && m.updateResult.UpdateAvail {
		updateShortcut = fmt.Sprintf("  •  %s %s", KeyBadgeStyle.Render("u"), "Actualizar SubiKit")
	}

	if m.focusPane == FocusList {
		footerText = fmt.Sprintf(
			" %s %s  •  %s %s  •  %s %s  •  %s %s  •  %s %s%s  •  Doc: %s",
			KeyBadgeStyle.Render("1-8"), "Pestañas",
			KeyBadgeStyle.Render("j/k"), "Moverse",
			KeyBadgeStyle.Render("l / Enter"), "Enfocar lectura",
			KeyBadgeStyle.Render("Tab"), "Cambiar foco",
			KeyBadgeStyle.Render("?"), "Ayuda",
			updateShortcut,
			lipgloss.NewStyle().Foreground(ColorHighlight).Render(scrollPercent),
		)
	} else {
		footerText = fmt.Sprintf(
			" %s %s  •  %s %s  •  %s %s  •  %s %s  •  %s %s  •  Líneas: %s",
			KeyBadgeVioletStyle.Render("j/k"), "Línea arriba/abajo",
			KeyBadgeVioletStyle.Render("d/u"), "1/2 Página",
			KeyBadgeVioletStyle.Render("g/G"), "Inicio/Fin",
			KeyBadgeVioletStyle.Render("h / Esc"), "Volver a lista",
			KeyBadgeVioletStyle.Render("Tab"), "Cambiar foco",
			lipgloss.NewStyle().Foreground(ColorSecondary).Render(scrollPercent),
		)
	}

	sb.WriteString(FooterStyle.Width(m.width).Render(footerText))

	return sb.String()
}

// RunTUI inicia la aplicación TUI interactiva de SubiKit.
func RunTUI(rulesMgr *rules.Manager, skillsMgr *skills.Manager, mcpMgr *mcp.Manager, agentsMgr *agents.Manager, version string) error {
	model := NewAppModel(rulesMgr, skillsMgr, mcpMgr, agentsMgr, version)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
