package mcp

import (
	"os"
	"testing"
)

func TestMCPManagerInstallAndDoctor(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "devkit-mcp-test-*")
	if err != nil {
		t.Fatalf("error creating temp home: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Simular HOME
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	mgr := NewManager()

	// 1. Catálogo
	cat := mgr.GetCatalog()
	if len(cat) != 3 {
		t.Errorf("expected 3 MCP definitions in catalog, got %d", len(cat))
	}

	// 2. Instalar Context7 con Token
	token := "ctx7sk-test-token-12345"
	if err := mgr.Install("context7", token, ""); err != nil {
		t.Fatalf("Install context7 failed: %v", err)
	}

	// 3. Instalar CodeGraph
	if err := mgr.Install("codegraph", "", ""); err != nil {
		t.Fatalf("Install codegraph failed: %v", err)
	}

	// 4. Leer configuración generada
	cfg, err := mgr.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	if len(cfg.MCPServers) != 2 {
		t.Errorf("expected 2 servers configured, got %d", len(cfg.MCPServers))
	}

	ctx7, ok := cfg.MCPServers["context7"]
	if !ok || ctx7.Headers["Authorization"] != "Bearer ctx7sk-test-token-12345" {
		t.Errorf("expected context7 with bearer token header, got %v", ctx7)
	}

	// 5. Verificar que se creó backup al reinstalar
	if err := mgr.Install("engram", "", "engram-test-cmd"); err != nil {
		t.Fatalf("Install engram failed: %v", err)
	}

	configPath, _ := mgr.GetConfigPath()
	if _, err := os.Stat(configPath + ".backup"); os.IsNotExist(err) {
		t.Errorf("expected backup file %s.backup to exist", configPath)
	}

	// 6. Doctor
	statuses := mgr.Doctor()
	if len(statuses) != 3 {
		t.Errorf("expected 3 doctor statuses, got %d", len(statuses))
	}
}
