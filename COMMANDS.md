# CLI Commands Reference

This document lists all available CLI commands and their flags, organized by category as shown in the `help` command output.

## Common Flags

Many tools support these standard flags:

- `--json` - Output in JSON format (default: human-friendly format)
- `--clip` - Copy output to clipboard
- `--file <path>` - Write output to file
- `--compact` - Use compact JSON format

---

## AI

### `cld` - Claude

Claude AI chat interface.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format

**Usage:**

```bash
cld "your message" [flags]
echo "message" | cld [flags]
```

### `gem` - Gemini

Google Gemini AI chat interface.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format

**Usage:**

```bash
gem "your message" [flags]
echo "message" | gem [flags]
```

### `gro` - Grok

Grok AI chat interface.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format
- `--prompt <path>` - Path to prompt file (skips interactive selection)
- `--test` - Test mode - use translate.md prompt

**Usage:**

```bash
gro "your message" [flags]
echo "message" | gro [flags]
```

### `grop` - Grok w/ prompts

Grok AI with prompt selection.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format
- `--prompt <path>` - Prompt file path (for testing, bypasses fzf)
- `--test` - Test mode - use translate.md prompt and default message

**Usage:**

```bash
grop "your message" [flags]
```

### `haik` - Haikus

Generate haikus using AI.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format

**Usage:**

```bash
haik "your prompt" [flags]
```

### `j` - ChatGPT

ChatGPT chat interface.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format

**Usage:**

```bash
j "your message" [flags]
echo "message" | j [flags]
```

### `ji` - ChatGPT w/ input.md

ChatGPT with input.md file context.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format

**Usage:**

```bash
ji "your message" [flags]
```

### `jj` - ChatGPT (JSON)

ChatGPT with JSON output format.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format
- `--prompt <path>` - Path to prompt file

**Usage:**

```bash
jj "your message" [flags]
```

### `jp` - ChatGPT w/ prompts

ChatGPT with prompt selection.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format
- `--prompt <path>` - Path to prompt file
- `--test` - Test mode - use translate.md prompt

**Usage:**

```bash
jp "your message" [flags]
```

### `prompts` - Open prompts in Cursor

Open prompts directory in editor.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--editor <name>` - Editor to use (cursor, code, vim, nvim, etc.) (default: cursor)
- `--json` - Output in JSON format

**Usage:**

```bash
prompts [flags]
prompts --editor code
```

---

## Git

### `gaff` - Show current files vs branched off

Show files changed since branching off from base branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format
- `--test` - Test mode - pre-select first changed file

**Usage:**

```bash
gaff [flags]
```

### `gbd` - Git branch delete

Delete a git branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gbd <branch-name> [flags]
```

### `gcb` - Git checkout branch

Create and checkout a new branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gcb <branch-name> [flags]
```

### `gcd` - Git checkout develop

Checkout to develop branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gcd [flags]
```

### `gcm` - Checkout to forkpoint branch

Checkout to the forkpoint branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gcm [flags]
```

### `gco` - Git checkout

Checkout a git branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gco <branch-name> [flags]
```

**Note:** Accepts JSON input via stdin with `{"branch": "branch-name"}` format.

### `gcommit` - Create commit message w/ AI

Generate commit message using AI based on changes.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gcommit [flags]
```

### `ginstall` - Install repo

Install repository dependencies (yarn → pnpm → npm).

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
ginstall [flags]
```

### `gmain` - Check main branch name

Check the name of the main branch.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gmain [flags]
```

### `gname` - Check repo name

Check the repository name.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gname [flags]
```

### `gprs` - Search for PRs by ticket ID

Search for pull requests containing the specified ticket ID.

**Flags:**

- `--compact` - Use compact JSON format
- `--single <path>` - Operate on specific repository path
- `--main` - Operate on main repositories (orbit + rasch-stack)
- `--all` - Operate on all repositories from config
- `--json` - Output in JSON format
- `--o` / `--open` - Open PR in browser

**Usage:**

```bash
gprs <ticket-id> [flags]
gprs PNT-123 --open
gprs PNT-123 --main
```

### `greinstall` - Reinstall repo

Reinstall repository dependencies.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
greinstall [flags]
```

### `grt` - Go to repo root

Navigate to repository root directory.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
grt [flags]
```

### `gs` - Git stash (smart)

Smart git stash operation.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gs [flags]
```

### `gsp` - Git stash pop

Pop the most recent stash.

**Flags:**

- `--compact` - Use compact JSON format
- `--json` - Output in JSON format

**Usage:**

```bash
gsp [flags]
```

### `gstats` - Show file / LOC stats of repo

Show file and lines of code statistics for repository.

**Flags:**

- `--compact` - Use compact JSON format
- `--single <path>` - Operate on specific repository path
- `--main` - Operate on main repositories (orbit + rasch-stack)
- `--all` - Operate on all repositories from config
- `--json` - Output in JSON format (default: formatted)

**Usage:**

```bash
gstats [flags]
gstats --main
gstats --all
gstats --single /path/to/repo
```

### `ghistory` - Show commit history across repositories

Show commit history with filtering by days and author.

**Flags:**

- `--compact` - Use compact JSON format
- `--single <path>` - Operate on specific repository path
- `--main` - Operate on main repositories (orbit + rasch-stack)
- `--all` - Operate on all repositories from config
- `--days <n>` - Number of days to look back (default: 7)
- `--author <name>` - Filter by author name
- `--json` - Output in JSON format

**Usage:**

```bash
ghistory [flags]
ghistory --days 14 --author "John Doe"
ghistory --main --days 30
```

### `gactivity` - Show user activity across repositories

Show user activity statistics across repositories.

**Flags:**

- `--compact` - Use compact JSON format
- `--single <path>` - Operate on specific repository path
- `--main` - Operate on main repositories (orbit + rasch-stack)
- `--all` - Operate on all repositories from config
- `--json` - Output in JSON format

**Usage:**

```bash
gactivity [flags]
gactivity --main
```

---

## Ringier

### `smart_start` - Smart Start

Run ginstall && \_smart_start.

**Usage:**

```bash
smart_start
```

---

## Core

### `check_alias` - Check what's behind an alias

Check what command an alias points to.

**Usage:**

```bash
check_alias <alias-name>
```

### `perf` - Measure time (5x)

Measure execution time of a command (runs 5 times).

**Usage:**

```bash
perf <command>
```

### `zsh` - Edit config (if no args, otherwise run)

Edit zsh configuration or run zsh commands.

**Usage:**

```bash
zsh [args]
```

### `zss` - Apply zshrc

Apply zshrc configuration.

**Usage:**

```bash
zss
```

---

## Tools

### `edit` - CLI tools config

Edit CLI tools configuration file.

**Flags:**

- `--editor <name>` - Editor to use (textedit, cursor, code, vim) (default: textedit)

**Usage:**

```bash
edit [flags]
edit --editor cursor
```

### `figma` - Figma component search and management tool

Figma component search and management tool.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--compact` - Use compact JSON output
- `--json` - Output as JSON (for piping)

**Commands:**

- `figma <query>` - Search for components by name (default)
- `figma search <query>` - Search for components
- `figma init [fileKey]` - Initialize cache with components from file
- `figma list` - List all cached components
- `figma cache <stats|clear>` - Manage cache
- `figma -- <query>` - Get full metadata for component

**Usage:**

```bash
figma "button" [flags]
figma search "button" --json
figma init
figma list
figma cache stats
figma cache clear
figma -- "button-main"
```

### `help` - Display formatted help

Display formatted help for all CLI tools (what you see now).

**Usage:**

```bash
help
```

### `jira` - Jira CLI tool

Jira issue management CLI tool.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--json` - Output in JSON format
- `--o` / `--open` - Open in browser

**Commands:**

- `jira <issue-key> [flags]` - View issue details
  - `-c` - Show comments
  - `-t` - Show testing instructions
  - `-o` - Open in browser (minimal display)
  - `-cc` - Comments only
- `jira history <issue-key>` - View issue history
- `jira user [username]` - View user activity
- `jira help` - Show help

**Usage:**

```bash
jira PNT-123
jira PNT-123 c
jira PNT-123 t
jira PNT-123 o
jira history PNT-123
jira user
jira user john.doe@company.com
```

### `test` - Run e2e tests from config.yml

Run end-to-end tests from config.yml.

**Usage:**

```bash
test
```

### `web` - Web search tool using Perplexity AI

Web search tool using Perplexity AI with intelligent caching.

**Flags:**

- `--clip` - Copy to clipboard
- `--file <path>` - Write to file
- `--compact` - Use compact JSON output
- `--json` - Output in JSON format

**Commands:**

- `web <query>` - Search the web
- `web cache stats` - Show cache statistics
- `web cache clear` - Clear cache

**Usage:**

```bash
web "your search query" [flags]
web cache stats
web cache clear
```

---

## Notes

### Input/Output

- Most AI tools accept input via command-line arguments or stdin
- JSON input can be piped to tools that support it (e.g., `gco` accepts `{"branch": "name"}`)
- Use `--json` flag for structured output when piping between tools
- Use `--clip` to copy results to clipboard for quick access

### Repository Operations

Many Git tools support repository scope flags:

- `--single <path>` - Operate on a specific repository
- `--main` - Operate only on main repositories (orbit + rasch-stack)
- `--all` - Operate on all repositories from config
- Default behavior varies by tool (usually current directory or all repos)

### Caching

Several tools use intelligent caching:

- `web` - Caches search results
- `figma` - Caches component metadata

Use `cache stats` and `cache clear` subcommands to manage cache.
