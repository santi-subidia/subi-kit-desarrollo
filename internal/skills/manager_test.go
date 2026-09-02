package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestSkillsManagerAndSDD(t *testing.T) {
	mockFS := fstest.MapFS{
		"sdd-workflow/SKILL.md": &fstest.MapFile{
			Data: []byte(`---
name: sdd-workflow
description: Spec-Driven Development skill
---
# SDD Content`),
		},
		"sdd-workflow/templates/spec.template.md": &fstest.MapFile{
			Data: []byte(`# Spec for {{FEATURE_NAME}} by {{AUTHOR}} on {{DATE}}`),
		},
		"sdd-workflow/templates/tasks.template.md": &fstest.MapFile{
			Data: []byte(`# Tasks for {{FEATURE_NAME}}`),
		},
	}

	mgr, err := NewManager(mockFS)
	if err != nil {
		t.Fatalf("error creating skills manager: %v", err)
	}

	all := mgr.GetAll()
	if len(all) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(all))
	}

	skill, ok := mgr.GetSkill("sdd-workflow")
	if !ok {
		t.Fatalf("expected to find 'sdd-workflow'")
	}
	if len(skill.Templates) != 2 {
		t.Errorf("expected 2 templates, got %d", len(skill.Templates))
	}

	// Test crear feature SDD
	tempDir, err := os.MkdirTemp("", "devkit-sdd-test-*")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	featureDir, err := mgr.CreateSDDFeature(tempDir, "auth-google-oauth", "Santiago")
	if err != nil {
		t.Fatalf("CreateSDDFeature returned error: %v", err)
	}

	specPath := filepath.Join(featureDir, "spec.md")
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("error reading generated spec.md: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "auth-google-oauth") || !strings.Contains(content, "Santiago") {
		t.Errorf("placeholders were not replaced correctly in spec.md: %s", content)
	}
}
