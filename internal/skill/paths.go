package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolvePath determines the installation directory for a skill
// based on scope ("user" or "project") and agent.
func ResolvePath(scope string, agent string) (string, error) {
	if scope == "" {
		scope = "user"
	}

	var baseDir string
	if scope == "user" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not find user home directory: %w", err)
		}
		baseDir = home
	} else if scope == "project" {
		// Use current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not get current working directory: %w", err)
		}
		baseDir = cwd
	} else {
		return "", fmt.Errorf("invalid scope %q, expected 'user' or 'project'", scope)
	}

	// Based on standard agent skill paths conventions
	var agentDir string
	switch agent {
	case "copilot", "cursor":
		agentDir = fmt.Sprintf(".%s", agent)
	case "":
		agentDir = ".agents"
	default:
		agentDir = fmt.Sprintf(".%s", agent)
	}

	return filepath.Join(baseDir, agentDir, "skills"), nil
}
