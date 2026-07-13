package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
)

func HandleList() error {
	scopes := []string{"project", "user"}
	var allSkills []InstalledSkill

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
				metaPath := filepath.Join(baseDest, entry.Name(), "metadata.json")
				metaBytes, err := os.ReadFile(metaPath)
				if err == nil {
					var meta InstalledSkill
					if err := json.Unmarshal(metaBytes, &meta); err == nil {
						allSkills = append(allSkills, meta)
					}
				}
			}
		}
	}

	if len(allSkills) == 0 {
		fmt.Println("No agent skills installed.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSCOPE\tAGENT\tSOURCE\tREVISION")
	for _, s := range allSkills {
		rev := s.Revision
		if len(rev) > 8 {
			rev = rev[:8] // abbreviate for display
		} else if rev == "" {
			rev = "unknown"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", s.Name, s.Scope, s.Agent, s.Source, rev)
	}
	w.Flush()

	return nil
}
