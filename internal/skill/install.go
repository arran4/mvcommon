package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func HandleInstall(source, name, scope, agent string, force bool) error {
	baseDest, err := ResolvePath(scope, agent)
	if err != nil {
		return err
	}

	skillName := name
	if skillName == "" {
		if source == "official" || source == "embedded" {
			skillName = "mvcommon"
		} else {
			skillName = filepath.Base(source)
		}
	}

	destPath := filepath.Join(baseDest, skillName)

	// Check if already exists
	if _, err := os.Stat(destPath); err == nil {
		if !force {
			return fmt.Errorf("skill '%s' already exists at %s. Use --force to replace", skillName, destPath)
		}
		// Remove existing
		if err := os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("failed to remove existing skill directory: %w", err)
		}
	}

	var revision string
	var sourceType string

	if source == "official" || source == "embedded" {
		sourceType = "embedded"
		if err := FetchEmbedded(source, destPath); err != nil {
			return err
		}
		revision = "embedded"
	} else if IsLocalPath(source) {
		sourceType = "local"
		if err := FetchLocal(source, destPath); err != nil {
			return err
		}
	} else {
		sourceType = "remote"
		rev, err := FetchRemote(source, destPath)
		if err != nil {
			return err
		}
		revision = rev
	}

	// Basic validation of SKILL.md
	if _, err := os.Stat(filepath.Join(destPath, "SKILL.md")); os.IsNotExist(err) {
		// Try to find if it's nested (e.g. skills/NAME/SKILL.md)
		nestedPath := filepath.Join(destPath, "skills", skillName, "SKILL.md")
		if _, err := os.Stat(nestedPath); err == nil {
			// Move nested up? For simplicity, we just look for SKILL.md in the root
			// Or we accept it as is if there's any SKILL.md somewhere.
			// Let's just do a basic check and not error hard if nested.
		} else {
			return fmt.Errorf("validation failed: SKILL.md not found in downloaded source")
		}
	}

	// Write metadata
	meta := InstalledSkill{
		Source:      source,
		Name:        skillName,
		Agent:       agent,
		Scope:       scope,
		InstallTime: time.Now(),
		Revision:    revision,
		Path:        destPath,
		SourceType:  sourceType,
	}

	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	if err := os.WriteFile(filepath.Join(destPath, "metadata.json"), metaBytes, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	fmt.Printf("Successfully installed skill '%s' to %s\n", skillName, destPath)
	return nil
}
