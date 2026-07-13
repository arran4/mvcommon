package skill

import "time"

// SkillManifest represents the logical information about a skill,
// often conceptually parsed from SKILL.md or the source provenance.
type SkillManifest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// InstalledSkill tracks the provenance of a skill once it is installed locally.
// It is typically saved alongside the skill in a metadata file (e.g. metadata.json)
type InstalledSkill struct {
	Source      string    `json:"source"`
	Name        string    `json:"name"`
	Agent       string    `json:"agent"`
	Scope       string    `json:"scope"`
	InstallTime time.Time `json:"install_time"`
	Revision    string    `json:"revision"`     // commit hash, tree hash, digest, or version
	Path        string    `json:"path"`         // the local directory path
	SourceType  string    `json:"source_type"`  // "local" or "remote"
}
