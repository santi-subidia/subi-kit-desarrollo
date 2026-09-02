package subikit

import (
	"embed"
	"io/fs"
)

// EmbeddedCatalogFS contiene todas las reglas, skills y agentes predeterminados embebidos en el binario.
//
//go:embed rules/* skills/* agents/*
var EmbeddedCatalogFS embed.FS

// GetRulesFS retorna el sub-filesystem apuntando al directorio de reglas embebidas.
func GetRulesFS() (fs.FS, error) {
	return fs.Sub(EmbeddedCatalogFS, "rules")
}

// GetSkillsFS retorna el sub-filesystem apuntando al directorio de skills embebidas.
func GetSkillsFS() (fs.FS, error) {
	return fs.Sub(EmbeddedCatalogFS, "skills")
}

// GetAgentsFS retorna el sub-filesystem apuntando al directorio de agentes embebidos.
func GetAgentsFS() (fs.FS, error) {
	return fs.Sub(EmbeddedCatalogFS, "agents")
}
