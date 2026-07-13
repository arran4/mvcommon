package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func HandleUpdate(name string, all bool) error {
	if all {
		fmt.Println("Updating all installed skills...")
		scopes := []string{"project", "user"}
		updatedCount := 0

		for _, scope := range scopes {
			baseDest, err := ResolvePath(scope, "")
			if err != nil {
				continue
			}

			entries, err := os.ReadDir(baseDest)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					err := updateSingleSkill(entry.Name())
					if err != nil {
						fmt.Printf("Failed to update skill %q: %v\n", entry.Name(), err)
					} else {
						updatedCount++
					}
				}
			}
		}

		if updatedCount == 0 {
			fmt.Println("No installed skills found to update.")
		} else {
			fmt.Printf("Successfully processed updates for %d skills.\n", updatedCount)
		}
		return nil
	}

	if name == "" {
		return fmt.Errorf("skill name required for update, or use --all")
	}

	return updateSingleSkill(name)
}

func updateSingleSkill(name string) error {
	// Determine where it's installed by searching scopes/agents
	scopes := []string{"project", "user"}
	var foundMeta *InstalledSkill

	for _, scope := range scopes {
		baseDest, err := ResolvePath(scope, "")
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
				break
			}
		}
	}

	if foundMeta == nil {
		return fmt.Errorf("skill '%s' not found locally or lacks metadata", name)
	}

	if foundMeta.SourceType == "embedded" {
		fmt.Printf("Skill '%s' is embedded. Updating from binary assets...\n", name)
		return HandleInstall(foundMeta.Source, name, foundMeta.Scope, foundMeta.Agent, true)
	}

	if foundMeta.SourceType == "local" {
		fmt.Printf("Skill '%s' was installed locally. Updating from local source: %s\n", name, foundMeta.Source)
		// Perform reinstall
		return HandleInstall(foundMeta.Source, name, foundMeta.Scope, foundMeta.Agent, true)
	}

	fmt.Printf("Updating skill '%s' from remote %s...\n", name, foundMeta.Source)

	// Create a temp dir to fetch the latest
	tmpDir, err := os.MkdirTemp("", "skill-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	latestRev, err := FetchRemote(foundMeta.Source, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to fetch remote update: %w", err)
	}

	if latestRev != "" && latestRev == foundMeta.Revision {
		fmt.Printf("Skill '%s' is already up to date (revision %s).\n", name, latestRev)
		return nil
	}

	fmt.Printf("New version found. Replacing locally installed skill...\n")

	// Reinstall forces replacement
	return HandleInstall(foundMeta.Source, name, foundMeta.Scope, foundMeta.Agent, true)
}
