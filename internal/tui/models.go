package tui

// Tab representa una sección en la barra de navegación superior.
type Tab int

const (
	TabCommands Tab = iota
	TabAgents
	TabRules
	TabSkills
	TabMCP
	TabDoctor
	TabSDD
	TabGlossary
)

var TabNames = []string{
	"1. Comandos",
	"2. Agentes",
	"3. Reglas",
	"4. Skills",
	"5. MCP",
	"6. Doctor",
	"7. SDD",
	"8. Glosario",
}

// FocusPane define qué panel tiene el foco del teclado (Lista o Detalle de lectura estilo Neovim).
type FocusPane int

const (
	FocusList FocusPane = iota
	FocusDetail
)

// CommandDoc documenta un comando de SubiKit para la TUI y el glosario.
type CommandDoc struct {
	Name        string
	Category    string
	Syntax      string
	Description string
	Flags       []CommandFlag
	Examples    []string
	Details     string
}

// CommandFlag documenta un flag de un comando.
type CommandFlag struct {
	Flag        string
	Description string
}

// GlossaryTerm documenta un concepto clave de desarrollo con IA en SubiKit.
type GlossaryTerm struct {
	Term        string
	Category    string
	Summary     string
	Explanation string
}

// SDDFeatureStatus representa el estado resumido de una feature en .specs/
type SDDFeatureStatus struct {
	Name      string
	Path      string
	Status    string
	HasSpec   bool
	HasPlan   bool
	HasTasks  bool
	HasVerify bool
	IsDone    bool
}
