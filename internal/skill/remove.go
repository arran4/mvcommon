package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func HandleRemove(name string, force bool) error {
	if name == "" {
		return fmt.Errorf("skill name required for removal")
	}

	scopes := []string{"project", "user"}
	var foundPath string

	for _, scope := range scopes {
		baseDest, err := ResolvePath(scope, "")
		if err != nil {
			continue
		}

		destPath := filepath.Join(baseDest, name)
		if _, err := os.Stat(destPath); err == nil {
			foundPath = destPath

			// Optional: load metadata to report what was removed
			metaPath := filepath.Join(destPath, "metadata.json")
			metaBytes, err := os.ReadFile(metaPath)
			if err == nil {
				var meta InstalledSkill
				if err := json.Unmarshal(metaBytes, &meta); err == nil {
					fmt.Printf("Found skill '%s' (scope: %s, agent: %s)\n", meta.Name, meta.Scope, meta.Agent)
				}
			}
			break
		}
	}

	if foundPath == "" {
		return fmt.Errorf("skill '%s' not found in any scope", name)
	}

	if err := os.RemoveAll(foundPath); err != nil {
		if !force {
			return fmt.Errorf("failed to remove skill directory (use --force to ignore): %w", err)
		}
	}

	fmt.Printf("Successfully removed skill '%s' from %s\n", name, foundPath)
	return nil
}
