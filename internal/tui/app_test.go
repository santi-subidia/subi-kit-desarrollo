package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	subikit "github.com/santi-subidia/dev-kit-desarrollo"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/agents"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/mcp"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/rules"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/skills"
	"github.com/santi-subidia/dev-kit-desarrollo/internal/updater"
)

func TestCommandsData(t *testing.T) {
	cmds := GetCommandsData()
	if len(cmds) == 0 {
		t.Fatal("GetCommandsData no retornó comandos")
	}

	foundInit := false
	foundTUI := false
	for _, c := range cmds {
		if c.Name == "init" {
			foundInit = true
		}
		if c.Name == "tui" {
			foundTUI = true
		}
		if c.Syntax == "" {
			t.Errorf("Comando %s no tiene sintaxis definida", c.Name)
		}
		if c.Description == "" {
			t.Errorf("Comando %s no tiene descripción", c.Name)
		}
	}

	if !foundInit {
		t.Errorf("Comando 'init' no encontrado en el catálogo de comandos")
	}
	if !foundTUI {
		t.Errorf("Comando 'tui' no encontrado en el catálogo de comandos")
	}
}

func TestGlossaryData(t *testing.T) {
	terms := GetGlossaryTerms()
	if len(terms) == 0 {
		t.Fatal("GetGlossaryTerms no retornó términos")
	}

	for _, g := range terms {
		if g.Term == "" || g.Summary == "" || g.Explanation == "" {
			t.Errorf("Término incompleto en glosario: %+v", g)
		}
	}
}

func TestAppModelNavigation(t *testing.T) {
	rulesFS, err := subikit.GetRulesFS()
	if err != nil {
		t.Fatalf("GetRulesFS error: %v", err)
	}
	rulesMgr, err := rules.NewManager(rulesFS)
	if err != nil {
		t.Fatalf("rules.NewManager error: %v", err)
	}

	skillsFS, err := subikit.GetSkillsFS()
	if err != nil {
		t.Fatalf("GetSkillsFS error: %v", err)
	}
	skillsMgr, err := skills.NewManager(skillsFS)
	if err != nil {
		t.Fatalf("skills.NewManager error: %v", err)
	}

	agentsFS, err := subikit.GetAgentsFS()
	if err != nil {
		t.Fatalf("GetAgentsFS error: %v", err)
	}
	agentsMgr, err := agents.NewManager(agentsFS)
	if err != nil {
		t.Fatalf("agents.NewManager error: %v", err)
	}

	mcpMgr := mcp.NewManager()

	model := NewAppModel(rulesMgr, skillsMgr, mcpMgr, agentsMgr, "0.4.0")

	// Simular evento WindowSize
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	appModel, ok := updatedModel.(AppModel)
	if !ok {
		t.Fatalf("Esperaba tipo AppModel")
	}

	// Comprobar pestañas
	for i := range TabNames {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(string(rune('1' + i)))}
		updatedModel, _ = appModel.Update(keyMsg)
		appModel = updatedModel.(AppModel)
		if appModel.activeTab != Tab(i) {
			t.Errorf("Esperaba pestaña %d, obtuve %d", i, appModel.activeTab)
		}
	}

	// Comprobar toggle de ayuda
	updatedModel, _ = appModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	appModel = updatedModel.(AppModel)
	if !appModel.showHelp {
		t.Errorf("Esperaba showHelp=true al presionar '?'")
	}

	// Cerrar ayuda
	updatedModel, _ = appModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	appModel = updatedModel.(AppModel)
	if appModel.showHelp {
		t.Errorf("Esperaba showHelp=false al presionar '?' nuevamente")
	}

	// Comprobar alternancia de foco (Neovim mode)
	if appModel.focusPane != FocusList {
		t.Errorf("Esperaba foco inicial en FocusList")
	}

	// Presionar 'l' para enfocar panel de lectura
	updatedModel, _ = appModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	appModel = updatedModel.(AppModel)
	if appModel.focusPane != FocusDetail {
		t.Errorf("Esperaba foco en FocusDetail tras presionar 'l'")
	}

	// Presionar 'j' en modo lectura (LineDown)
	updatedModel, _ = appModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	appModel = updatedModel.(AppModel)

	// Presionar 'h' para volver a la lista
	updatedModel, _ = appModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	appModel = updatedModel.(AppModel)
	if appModel.focusPane != FocusList {
		t.Errorf("Esperaba foco en FocusList tras presionar 'h'")
	}

	// Renderizar vista no debe hacer panic
	viewOutput := appModel.View()
	if len(viewOutput) == 0 {
		t.Errorf("View retornó cadena vacía")
	}
}

func TestAppModelUpdateState(t *testing.T) {
	rulesFS, _ := subikit.GetRulesFS()
	rulesMgr, _ := rules.NewManager(rulesFS)
	skillsFS, _ := subikit.GetSkillsFS()
	skillsMgr, _ := skills.NewManager(skillsFS)
	agentsFS, _ := subikit.GetAgentsFS()
	agentsMgr, _ := agents.NewManager(agentsFS)
	mcpMgr := mcp.NewManager()

	model := NewAppModel(rulesMgr, skillsMgr, mcpMgr, agentsMgr, "0.4.0")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	appModel := updatedModel.(AppModel)

	// Simular recepción de update disponible
	fakeResult := &updateCheckMsg{
		result: &updater.UpdateResult{
			CurrentVersion: "0.4.0",
			LatestVersion:  "v0.5.0",
			UpdateAvail:    true,
			Release: &updater.ReleaseInfo{
				TagName: "v0.5.0",
				Name:    "SubiKit v0.5.0",
			},
		},
	}

	updatedModel, _ = appModel.Update(*fakeResult)
	appModel = updatedModel.(AppModel)

	if appModel.updateResult == nil || !appModel.updateResult.UpdateAvail {
		t.Fatalf("updateResult no se guardó correctamente")
	}

	// Verificar que la vista muestre el badge de actualización
	view := appModel.View()
	if len(view) == 0 {
		t.Errorf("View retornó cadena vacía")
	}

	// Probar atajo 'u'
	updatedModel, cmd := appModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	appModel = updatedModel.(AppModel)
	if !appModel.isUpdating {
		t.Errorf("Esperaba isUpdating=true tras presionar 'u'")
	}
	if cmd == nil {
		t.Errorf("Esperaba cmd no nulo para iniciar la actualización")
	}
}


