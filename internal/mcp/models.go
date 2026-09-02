package mcp

// ServerConfig representa la configuración de un servidor MCP individual según la especificación.
type ServerConfig struct {
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	ServerURL string            `json:"serverUrl,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}

// ConfigFile representa el archivo mcp_config.json de Antigravity / Gemini.
type ConfigFile struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// Definition define un servidor MCP soportado por el Dev-Kit.
type Definition struct {
	ID          string
	Name        string
	Description string
	Type        string // "stdio" o "http"
	DefaultCmd  string
	DefaultArgs []string
	DefaultURL  string
	RequiresAuth bool
	EnvKey      string
}

// Catálogo de MCPs soportados
var Catalog = []Definition{
	{
		ID:          "engram",
		Name:        "Engram Memory",
		Description: "Memoria persistente a largo plazo y recuperación de contexto entre sesiones y proyectos",
		Type:        "stdio",
		DefaultCmd:  "engram",
		DefaultArgs: []string{"mcp", "--tools=agent"},
	},
	{
		ID:           "context7",
		Name:         "Context7 Docs",
		Description:  "Documentación oficial y actualizada en tiempo real de librerías y frameworks",
		Type:         "http",
		DefaultURL:   "https://mcp.context7.com/mcp",
		RequiresAuth: true,
		EnvKey:       "CONTEXT7_API_KEY",
	},
	{
		ID:          "codegraph",
		Name:        "CodeGraph",
		Description: "Mapeo semántico de repositorios, navegación de grafos de dependencias e impacto",
		Type:        "stdio",
		DefaultCmd:  "codegraph",
		DefaultArgs: []string{"serve", "--mcp"},
	},
}

// HealthStatus representa el diagnóstico de salud de un servidor MCP.
type HealthStatus struct {
	ID          string
	Name        string
	Installed   bool
	Type        string
	Executable  string
	FoundInPath bool
	Details     string
}
