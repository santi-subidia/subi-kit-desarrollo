package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager administra la configuración, instalación y diagnóstico de servidores MCP.
type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// GetCatalog retorna la lista de MCPs soportados en el kit.
func (m *Manager) GetCatalog() []Definition {
	return Catalog
}

// GetConfigPath retorna la ruta al archivo de configuración global de MCPs (~/.gemini/config/mcp_config.json).
func (m *Manager) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("no se pudo determinar el directorio de usuario: %w", err)
	}
	return filepath.Join(homeDir, ".gemini", "config", "mcp_config.json"), nil
}

// ReadConfig lee el archivo de configuración actual de MCPs.
func (m *Manager) ReadConfig() (*ConfigFile, error) {
	configPath, err := m.GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return &ConfigFile{MCPServers: make(map[string]ServerConfig)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al leer configuración MCP %s: %w", configPath, err)
	}

	var cfg ConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error al parsear JSON de MCPs: %w", err)
	}

	if cfg.MCPServers == nil {
		cfg.MCPServers = make(map[string]ServerConfig)
	}

	return &cfg, nil
}

// SaveConfig guarda el archivo de configuración realizando un backup previo.
func (m *Manager) SaveConfig(cfg *ConfigFile) error {
	configPath, err := m.GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio de configuración %s: %w", dir, err)
	}

	// Backup si ya existe
	if _, err := os.Stat(configPath); err == nil {
		backupPath := configPath + ".backup"
		if oldData, err := os.ReadFile(configPath); err == nil {
			_ = os.WriteFile(backupPath, oldData, 0644)
		}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error al formatear JSON de MCPs: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error al escribir %s: %w", configPath, err)
	}

	return nil
}

// Install instala o actualiza un servidor MCP por su ID.
func (m *Manager) Install(id string, token string, customCmd string) error {
	var targetDef *Definition
	for _, d := range Catalog {
		if strings.EqualFold(d.ID, id) {
			targetDef = &d
			break
		}
	}

	if targetDef == nil {
		return fmt.Errorf("servidor MCP no reconocido: '%s'", id)
	}

	cfg, err := m.ReadConfig()
	if err != nil {
		return err
	}

	serverConfig := ServerConfig{}

	if targetDef.Type == "http" {
		serverConfig.ServerURL = targetDef.DefaultURL
		// Token de flag o de variable de entorno
		authToken := token
		if authToken == "" && targetDef.EnvKey != "" {
			authToken = os.Getenv(targetDef.EnvKey)
		}

		if authToken != "" {
			serverConfig.Headers = map[string]string{
				"Authorization": "Bearer " + strings.TrimPrefix(authToken, "Bearer "),
			}
		}
	} else {
		cmdName := targetDef.DefaultCmd
		if customCmd != "" {
			cmdName = customCmd
		} else {
			// Intentar resolver en PATH o buscar en GOPATH/bin si es engram
			if resolved, err := exec.LookPath(cmdName); err == nil {
				cmdName = resolved
			}
		}
		serverConfig.Command = cmdName
		serverConfig.Args = targetDef.DefaultArgs
	}

	cfg.MCPServers[targetDef.ID] = serverConfig
	return m.SaveConfig(cfg)
}

// InstallAll instala todos los servidores MCP del catálogo.
func (m *Manager) InstallAll(token string) ([]string, error) {
	var installed []string
	for _, def := range Catalog {
		if err := m.Install(def.ID, token, ""); err != nil {
			return installed, fmt.Errorf("error al instalar %s: %w", def.ID, err)
		}
		installed = append(installed, def.ID)
	}
	return installed, nil
}

// Doctor verifica el estado y la disponibilidad de los ejecutables de los MCPs.
func (m *Manager) Doctor() []*HealthStatus {
	cfg, _ := m.ReadConfig()
	var statuses []*HealthStatus

	for _, def := range Catalog {
		status := &HealthStatus{
			ID:   def.ID,
			Name: def.Name,
			Type: def.Type,
		}

		if cfg != nil {
			if srv, ok := cfg.MCPServers[def.ID]; ok {
				status.Installed = true
				if def.Type == "http" {
					if _, hasAuth := srv.Headers["Authorization"]; hasAuth {
						status.Details = "Configurado con Token de Autenticación"
					} else {
						status.Details = "Configurado sin Token (Público/Anónimo)"
					}
					status.FoundInPath = true
				} else {
					status.Executable = srv.Command
					if _, err := exec.LookPath(srv.Command); err == nil || fileExists(srv.Command) {
						status.FoundInPath = true
						status.Details = "Binario encontrado y ejecutable"
					} else {
						status.FoundInPath = false
						status.Details = "Binario no encontrado en el PATH"
					}
				}
			} else {
				status.Installed = false
				status.Details = "No configurado en mcp_config.json"
			}
		}

		statuses = append(statuses, status)
	}

	return statuses
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
