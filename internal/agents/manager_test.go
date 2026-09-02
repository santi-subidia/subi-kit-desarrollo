package agents

import (
	"testing"
	"testing/fstest"
)

func TestAgentsManager(t *testing.T) {
	mockFS := fstest.MapFS{
		"orchestrator.md": &fstest.MapFile{
			Data: []byte(`---
name: orchestrator
title: Tech Lead Orchestrator
type: orchestrator
description: Coordinator
tools: [all]
subagents: [architect, fullstack]
---
# Orchestrator instructions`),
		},
		"subagents/architect.md": &fstest.MapFile{
			Data: []byte(`---
name: architect
title: Software Architect
type: subagent
description: Architecture specialist
tools: [read, codegraph]
---
# Architect instructions`),
		},
		"subagents/fullstack.md": &fstest.MapFile{
			Data: []byte(`---
name: fullstack
title: Fullstack Dev
type: subagent
description: Developer
tools: [read, write]
---
# Fullstack instructions`),
		},
	}

	mgr, err := NewManager(mockFS)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	all := mgr.GetAll()
	if len(all) != 3 {
		t.Fatalf("expected 3 agents, got %d", len(all))
	}

	orch := mgr.GetOrchestrator()
	if orch == nil || orch.Metadata.Name != "orchestrator" {
		t.Errorf("expected orchestrator agent, got %v", orch)
	}

	subs := mgr.GetSubagents()
	if len(subs) != 2 {
		t.Errorf("expected 2 subagents, got %d", len(subs))
	}

	arch, ok := mgr.GetAgent("architect")
	if !ok || arch.Metadata.Title != "Software Architect" {
		t.Errorf("expected to find architect subagent, got %v", arch)
	}
}
