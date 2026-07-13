package main

import (
	"github.com/arran4/mvcommon/internal/skill"
)

// SkillsCmd is a subcommand `mvcommon skill` -- Manage agent skills
// Agent skills teach AI coding agents how to use this CLI correctly.
func SkillsCmd() error {
    return nil
}

// SkillsInstall is a subcommand `mvcommon skill install` -- Install a new agent skill
// Install a skill from a remote repository or a local directory.
//
// Flags:
//
//	source: @1 The source to install from (repository or local path)
//	scope: --scope (default: "user") The scope to install to (user or project)
//	agent: --agent The agent to install for (e.g. copilot, cursor)
//	force: --force Force installation, replacing locally modified skills
//	optionalNames: ... The name of the skill to install (optional, used to disambiguate multi-skill repos)
func SkillsInstall(source string, scope string, agent string, force bool, optionalNames ...string) error {
    name := ""
    if len(optionalNames) > 0 {
        name = optionalNames[0]
    }
	return skill.HandleInstall(source, name, scope, agent, force)
}

// SkillsUpdate is a subcommand `mvcommon skill update` -- Update an agent skill
// Update an installed skill.
//
// Flags:
//
//	all: --all Update all installed skills
//	optionalNames: ... The name of the skill to update (optional if --all is used)
func SkillsUpdate(all bool, optionalNames ...string) error {
    name := ""
    if len(optionalNames) > 0 {
        name = optionalNames[0]
    }
	return skill.HandleUpdate(name, all)
}

// SkillsRemove is a subcommand `mvcommon skill remove` -- Remove an agent skill
// Remove an installed skill.
//
// Flags:
//
//	name: @1 The name of the skill to remove
//	force: --force Force removal
func SkillsRemove(name string, force bool) error {
	return skill.HandleRemove(name, force)
}

// SkillsList is a subcommand `mvcommon skill list` -- List installed agent skills
// List all installed skills.
func SkillsList() error {
	return skill.HandleList()
}

// SkillsInspect is a subcommand `mvcommon skill inspect` -- Inspect an agent skill
// Inspect the metadata and provenance of an installed skill.
//
// Flags:
//
//	name: @1 The name of the skill to inspect
func SkillsInspect(name string) error {
	return skill.HandleInspect(name)
}
