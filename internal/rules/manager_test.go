package rules

import (
	"testing"
	"testing/fstest"
)

func TestParseRuleAndManager(t *testing.T) {
	mockFS := fstest.MapFS{
		"core/git.md": &fstest.MapFile{
			Data: []byte(`---
name: git-rule
title: Git Rule
category: core
always_on: true
description: Git rules
tags: [git, core]
---
# Git content`),
		},
		"frontend/next.md": &fstest.MapFile{
			Data: []byte(`---
name: next-rule
title: Next Rule
category: frontend
always_on: false
description: Nextjs rules
tags: [nextjs, react]
---
# Next content`),
		},
	}

	mgr, err := NewManager(mockFS)
	if err != nil {
		t.Fatalf("unexpected error creating manager: %v", err)
	}

	all := mgr.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(all))
	}

	coreRules := mgr.GetCoreRules()
	if len(coreRules) != 1 || coreRules[0].Metadata.Name != "git-rule" {
		t.Errorf("expected 1 core rule 'git-rule', got %v", coreRules)
	}

	matched := mgr.MatchRules([]string{"nextjs"}, true)
	if len(matched) != 2 {
		t.Errorf("expected 2 matched rules (1 always_on + 1 nextjs), got %d", len(matched))
	}
}
