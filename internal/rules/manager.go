package rules

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Manager administra la carga, consulta y filtrado del catálogo de reglas.
type Manager struct {
	rules []*Rule
	byName map[string]*Rule
}

// NewManager crea un nuevo gestor a partir de un sistema de archivos (embebido o local).
func NewManager(catalogFS fs.FS) (*Manager, error) {
	m := &Manager{
		rules:  make([]*Rule, 0),
		byName: make(map[string]*Rule),
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
			return fmt.Errorf("error al leer regla %s: %w", path, err)
		}

		// Normalizar separadores a '/'
		normPath := filepath.ToSlash(path)
		rule, err := ParseRule(data, normPath)
		if err != nil {
			// Si no tiene frontmatter YAML, ignoramos archivos markdown que no sean reglas
			return nil
		}

		m.rules = append(m.rules, rule)
		m.byName[rule.Metadata.Name] = rule
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Ordenar reglas alfabéticamente por categoría y nombre
	sort.Slice(m.rules, func(i, j int) bool {
		if m.rules[i].Metadata.Category == m.rules[j].Metadata.Category {
			return m.rules[i].Metadata.Name < m.rules[j].Metadata.Name
		}
		return m.rules[i].Metadata.Category < m.rules[j].Metadata.Category
	})

	return m, nil
}

// GetAll retorna todas las reglas disponibles en el catálogo.
func (m *Manager) GetAll() []*Rule {
	return m.rules
}

// GetRule busca una regla por su nombre único.
func (m *Manager) GetRule(name string) (*Rule, bool) {
	r, ok := m.byName[name]
	return r, ok
}

// GetCoreRules retorna las reglas marcadas con always_on o categoría 'core'.
func (m *Manager) GetCoreRules() []*Rule {
	var result []*Rule
	for _, r := range m.rules {
		if r.Metadata.AlwaysOn || strings.EqualFold(r.Metadata.Category, "core") {
			result = append(result, r)
		}
	}
	return result
}

// MatchRules filtra las reglas según una lista de tags detectadas o solicitadas.
// Si includeAlwaysOn es true, incluye siempre las reglas marcadas con always_on / core.
func (m *Manager) MatchRules(detectedTags []string, includeAlwaysOn bool) []*Rule {
	tagSet := make(map[string]bool)
	for _, t := range detectedTags {
		tagSet[strings.ToLower(strings.TrimSpace(t))] = true
	}

	seen := make(map[string]bool)
	var matched []*Rule

	for _, r := range m.rules {
		shouldInclude := false

		if includeAlwaysOn && (r.Metadata.AlwaysOn || strings.EqualFold(r.Metadata.Category, "core")) {
			shouldInclude = true
		} else {
			for _, tag := range r.Metadata.Tags {
				if tagSet[strings.ToLower(tag)] {
					shouldInclude = true
					break
				}
			}
		}

		if shouldInclude && !seen[r.Metadata.Name] {
			seen[r.Metadata.Name] = true
			matched = append(matched, r)
		}
	}

	return matched
}

// GetCategories retorna la lista única de categorías disponibles.
func (m *Manager) GetCategories() []string {
	catMap := make(map[string]bool)
	for _, r := range m.rules {
		if r.Metadata.Category != "" {
			catMap[r.Metadata.Category] = true
		}
	}

	var cats []string
	for c := range catMap {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	return cats
}
