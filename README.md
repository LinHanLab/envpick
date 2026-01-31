# envpick

[English](README.md) | [简体中文](README.zh-CN.md)

A CLI tool for managing multiple environment variable configurations with interactive selection via fzf.

<img src="docs/vhs/demo.gif" width="500" alt="Demo">

**[Installation Guide](docs/installation.md)** | Quick Start below

## Quick Start

This guide helps you get started with envpick in 3 simple steps.

### 1. Set up your shell (zsh)

Add to `~/.zshrc`:

```bash
# Initialize zsh completion system (if not already done)
autoload -U compinit; compinit

# Initialize envpick
eval "$(envpick init zsh)"
```

Reload your shell:

```bash
source ~/.zshrc
```

This enables completion and adds the `ep` shortcut command. 

### 2. Create your configuration

envpick uses `~/.envpick/config.toml` to store your environment configurations:

```toml
[personal]
ANTHROPIC_BASE_URL = "https://api.anthropic.com"
ANTHROPIC_API_KEY = "sk-ant-personal-xxxxx"
ANTHROPIC_AUTH_TOKEN = ""
ANTHROPIC_MODEL = "claude-sonnet-4-5"
ANTHROPIC_SMALL_FAST_MODEL = "claude-haiku-4"
API_TIMEOUT_MS = "300000"
_web_url = "https://console.anthropic.com"

[work]
ANTHROPIC_BASE_URL = "https://api.company.com"
ANTHROPIC_API_KEY = ""
ANTHROPIC_AUTH_TOKEN = "sk-work-token-xxxxx"
ANTHROPIC_MODEL = "claude-opus-4-5"
ANTHROPIC_SMALL_FAST_MODEL = "claude-sonnet-4-5"
API_TIMEOUT_MS = "600000"
_web_url = "https://dashboard.company.com"
```

Variables starting with `_` are metadata (e.g., `_web_url` for web URLs).

Edit your config anytime with:

```bash
ep edit
```

### 3. Switch between configurations

Use interactive selection to switch between your configurations:

```bash
ep use
```

This opens fzf for interactive selection. Your choice persists across new terminal sessions.

### Using Namespaces (Advanced)

When you have multiple groups of related configurations (e.g., databases, APIs), use namespaces:

```toml
[db.local]
DATABASE_URL = "postgres://localhost/myapp"

[db.staging]
DATABASE_URL = "postgres://staging.example.com/myapp"

[db.prod]
DATABASE_URL = "postgres://prod.example.com/myapp"
```

Switch within a namespace:

```bash
ep use -n db
```

Each namespace maintains its own state independently.

### Temporary Configuration (One-time Use)

For temporary configuration changes that don't persist:

```bash
# Interactive selection (current terminal only)
ep tmp

# Direct selection
ep tmp work

# With namespace
ep tmp -n db staging
```

Use `ep tmp` when you need different environment variables for a single terminal session without changing your persistent configuration.

## Features

- Interactive configuration switching with fzf
- Persistent state across terminal sessions
- Namespace support for organized configs
- Web URL launcher: `envpick web`
- Temporary config selection: `envpick env select`
- Shell integration with `ep` helper function

For complete command documentation: `envpick --help`
