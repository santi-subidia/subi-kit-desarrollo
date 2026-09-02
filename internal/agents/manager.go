package agents

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Manager administra la carga, consulta y filtrado de agentes y subagentes del dev-kit.
type Manager struct {
	agents []*Agent
	byName map[string]*Agent
}

// NewManager carga los agentes desde un filesystem embebido o local.
func NewManager(catalogFS fs.FS) (*Manager, error) {
	m := &Manager{
		agents: make([]*Agent, 0),
		byName: make(map[string]*Agent),
	}

	err := fs.WalkDir(catalogFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := fs.ReadFile(catalogFS, path)
		if err != nil {
			return fmt.Errorf("error al leer agente %s: %w", path, err)
		}

		normPath := filepath.ToSlash(path)
		agent, err := ParseAgent(data, normPath)
		if err != nil {
			return nil
		}

		m.agents = append(m.agents, agent)
		m.byName[agent.Metadata.Name] = agent
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(m.agents, func(i, j int) bool {
		// Orquestador primero, luego subagentes alfabéticamente
		if m.agents[i].Metadata.Type == "orchestrator" {
			return true
		}
		if m.agents[j].Metadata.Type == "orchestrator" {
			return false
		}
		return m.agents[i].Metadata.Name < m.agents[j].Metadata.Name
	})

	return m, nil
}

// GetAll retorna todos los agentes y subagentes.
func (m *Manager) GetAll() []*Agent {
	return m.agents
}

// GetOrchestrator retorna el agente principal orquestador.
func (m *Manager) GetOrchestrator() *Agent {
	for _, a := range m.agents {
		if a.Metadata.Type == "orchestrator" {
			return a
		}
	}
	return nil
}

// GetSubagents retorna solo los subagentes especializados.
func (m *Manager) GetSubagents() []*Agent {
	var subs []*Agent
	for _, a := range m.agents {
		if a.Metadata.Type != "orchestrator" {
			subs = append(subs, a)
		}
	}
	return subs
}

// GetAgent busca un agente o subagente por su nombre.
func (m *Manager) GetAgent(name string) (*Agent, bool) {
	a, ok := m.byName[name]
	return a, ok
}
