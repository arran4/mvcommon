package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	// Use project scope for tests to be deterministic
	cwd, _ := os.Getwd()

	tests := []struct {
		scope    string
		agent    string
		expected string
	}{
		{"project", "copilot", filepath.Join(cwd, ".copilot", "skills")},
		{"project", "cursor", filepath.Join(cwd, ".cursor", "skills")},
		{"project", "", filepath.Join(cwd, ".agents", "skills")},
	}

	for _, tt := range tests {
		t.Run(tt.scope+"-"+tt.agent, func(t *testing.T) {
			path, err := ResolvePath(tt.scope, tt.agent)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if path != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, path)
			}
		})
	}
}

func TestLocalInstall(t *testing.T) {
	// Create a mock local skill
	srcDir := t.TempDir()
	skillDir := filepath.Join(srcDir, "my-test-skill")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Test Skill"), 0644)

	// Switch to a temp directory for "project" scope
	projDir := t.TempDir()
	origWd, _ := os.Getwd()
	os.Chdir(projDir)
	defer os.Chdir(origWd)

	err := HandleInstall(skillDir, "test-skill", "project", "", false)
	if err != nil {
		t.Fatalf("installation failed: %v", err)
	}

	expectedDest := filepath.Join(projDir, ".agents", "skills", "test-skill")
	if _, err := os.Stat(expectedDest); os.IsNotExist(err) {
		t.Fatalf("expected skill to be installed at %q", expectedDest)
	}

	// Verify metadata
	metaPath := filepath.Join(expectedDest, "metadata.json")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Errorf("expected metadata.json to be created")
	}

	// Verify remove
	err = HandleRemove("test-skill", false)
	if err != nil {
		t.Fatalf("removal failed: %v", err)
	}

	if _, err := os.Stat(expectedDest); !os.IsNotExist(err) {
		t.Errorf("expected skill directory to be removed")
	}
}
