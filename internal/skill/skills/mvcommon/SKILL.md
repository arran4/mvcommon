# mvcommon AI Agent Skill

This skill teaches AI coding agents how to use the `mvcommon` CLI tool correctly.
`mvcommon` is a CLI tool that automatically groups files into directories based on common prefixes.

## Core concepts

- The CLI uses `-` prefixes for flags, not `--`. (e.g. `-dry-run`, `-interactive`).
- By default, it runs non-interactively if `-interactive` is not provided, making it safe for automation scripts and agents.
- It is destructive (moves files) by default. Use `-dry-run` to test operations safely without mutating the file system.

## Important Flags

- `-dry-run`: Always use this first when running operations to see what changes would occur!
- `-interactive`: Only use this when communicating with a human user in a real TTY. Agents should NOT use this.
- `-stopword="<word>"`: Stop detecting common prefixes at these words. Defaults: ` - `, `] `, `[`.
- `-trim="<chars>"`: Characters to trim from the start/end of the folder prefix. Defaults: `-_ .`.
- `-min=<number>`: Minimum size of the shared prefix to trigger a match. Default is 3.

## Examples

### Safely testing a move operation

```bash
# Good: dry-run first
mvcommon -dry-run *.txt
```

### Applying a move operation

```bash
# Good: move the files
mvcommon *.txt
```

## Agent Skill Management

`mvcommon` supports agent skills. You can install skills for yourself to learn how to use other tools.

### Commands

- `mvcommon skill list`: Show installed skills
- `mvcommon skill install <source> [name] --scope project`: Install a skill from a GitHub repo (e.g. `owner/repo`) or local path.
- `mvcommon skill update <name>`: Update an installed skill
- `mvcommon skill inspect <name>`: View metadata about an installed skill
- `mvcommon skill remove <name>`: Remove an installed skill
