package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func HandleInspect(name string) error {
	if name == "" {
		return fmt.Errorf("skill name required for inspection")
	}

	scopes := []string{"project", "user"}
	var foundMeta *InstalledSkill
	var foundPath string

	for _, scope := range scopes {
		for _, agent := range []string{"", "copilot", "cursor"} {
			baseDest, err := ResolvePath(scope, agent)
			if err != nil {
				continue
			}

			destPath := filepath.Join(baseDest, name)
			metaPath := filepath.Join(destPath, "metadata.json")

			metaBytes, err := os.ReadFile(metaPath)
			if err == nil {
				var meta InstalledSkill
				if err := json.Unmarshal(metaBytes, &meta); err == nil {
					foundMeta = &meta
					foundPath = destPath
					break
				}
			}
		}
		if foundMeta != nil {
			break
		}
	}

	if foundMeta == nil {
		return fmt.Errorf("skill '%s' not found or lacks metadata", name)
	}

	fmt.Printf("Skill: %s\n", foundMeta.Name)
	fmt.Printf("Source: %s\n", foundMeta.Source)
	fmt.Printf("Source Type: %s\n", foundMeta.SourceType)
	fmt.Printf("Agent: %s\n", foundMeta.Agent)
	fmt.Printf("Scope: %s\n", foundMeta.Scope)
	fmt.Printf("Path: %s\n", foundPath)
	fmt.Printf("Revision: %s\n", foundMeta.Revision)
	fmt.Printf("Installed: %s\n", foundMeta.InstallTime.Format("2006-01-02 15:04:05"))

	// Check if SKILL.md exists and print its size
	skillMDPath := filepath.Join(foundPath, "SKILL.md")
	if info, err := os.Stat(skillMDPath); err == nil {
		fmt.Printf("\nSKILL.md: found (%d bytes)\n", info.Size())
	} else {
		// try to find it nested
		nestedPath := filepath.Join(foundPath, "skills", foundMeta.Name, "SKILL.md")
		if info, err := os.Stat(nestedPath); err == nil {
			fmt.Printf("\nSKILL.md: found in nested path (%d bytes)\n", info.Size())
		} else {
			fmt.Printf("\nSKILL.md: not found\n")
		}
	}

	return nil
}
