# envpick

A CLI tool for managing multiple environment variable configurations with interactive selection via fzf.

## Installation

### Download Pre-built Binary (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/LinHanLab/envpick/releases).

**Linux (x86_64)**
```bash
curl -L https://github.com/LinHanLab/envpick/releases/latest/download/envpick_Linux_x86_64.tar.gz | tar xz
sudo mv envpick /usr/local/bin/
```

**Linux (ARM64)**
```bash
curl -L https://github.com/LinHanLab/envpick/releases/latest/download/envpick_Linux_arm64.tar.gz | tar xz
sudo mv envpick /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -L https://github.com/LinHanLab/envpick/releases/latest/download/envpick_Darwin_x86_64.tar.gz | tar xz
sudo mv envpick /usr/local/bin/
```

**macOS (Apple Silicon)**
```bash
curl -L https://github.com/LinHanLab/envpick/releases/latest/download/envpick_Darwin_arm64.tar.gz | tar xz
sudo mv envpick /usr/local/bin/
```

**Windows**

Download `envpick_Windows_x86_64.zip` or `envpick_Windows_arm64.zip` from the releases page, extract it, and add the directory to your PATH.

### Install via Go

If you have Go installed:

```bash
go install github.com/LinHanLab/envpick@latest
```

### Build from Source

```bash
git clone https://github.com/LinHanLab/envpick.git
cd envpick
make compile
sudo mv envpick /usr/local/bin/
```

## Requirements

- **fzf** must be installed and available in your PATH for interactive selection
  - macOS: `brew install fzf`
  - Linux: `sudo apt-get install fzf` or `sudo yum install fzf`
  - Windows: `choco install fzf` or download from [fzf releases](https://github.com/junegunn/fzf/releases)


## Configuration

Create `~/.envpick/config.toml`:

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

# Namespace configurations (grouped by prefix)
[db.local]
DATABASE_URL = "postgres://localhost/myapp"

[db.staging]
DATABASE_URL = "postgres://staging.example.com/myapp"

[db.prod]
DATABASE_URL = "postgres://prod.example.com/myapp"
```

Variables starting with `_` are metadata:
- `_web_url` - URL to open with `envpick web`

### Namespaces

Configurations can be organized into namespaces using dot notation (e.g., `db.local`, `db.staging`). This allows you to group related configurations together.

Use the `-n/--namespace` flag to filter operations to a specific namespace:

```bash
# List and select from db namespace only
envpick use -n db

# Use db namespace configuration
eval "$(envpick env -n db)"
```

Each namespace maintains its own state independently.

## Quick Setup

Add to your `~/.zshrc`:

```bash
# Initialize zsh completion system (if not already done)
autoload -U compinit; compinit

# Initialize envpick
eval "$(envpick init zsh)"
```

This sets up:
- Shell completion for envpick commands
- Automatic loading of persisted environment on shell startup
- Helper function `ep` for convenient operations

After adding, reload your shell:

```bash
source ~/.zshrc
```

### Helper Function

The `ep` function provides convenient shortcuts:

```bash
# Interactive selection with persistence
ep use
ep use -n db

# Temporary selection (no persistence)
ep tmp
ep tmp prod
ep tmp -n db local

# Pass through to envpick
ep edit
ep web
```

## Commands

### `envpick init zsh`

Generate shell configuration for envpick integration.

```bash
# Output configuration
envpick init zsh

# Add to your shell profile
eval "$(envpick init zsh)"
```

This command outputs shell-specific configuration that you can add to your shell profile. The configuration includes completion setup, automatic environment loading, and helper functions.

### `envpick use`

Interactively switch the persistent configuration. New terminal sessions will use the selected config.

```bash
envpick use

# Use with namespace
envpick use -n db
```

### `envpick env`

Output export statements for the current configuration. Add to your shell profile:

```bash
# In ~/.zshrc or ~/.bashrc
eval "$(envpick env)"

# Use with namespace
eval "$(envpick env -n db)"
```

### `envpick env select [config-name]`

Select a configuration and output its export statements. Does not change persistent state.

```bash
# Interactive selection
eval "$(envpick env select)"

# Direct selection (one-time, non-persistent)
eval "$(envpick env select prod)"

# Use with namespace
eval "$(envpick env select -n db local)"
```

### `envpick edit`

Open `~/.envpick/config.toml` in your editor (`$EDITOR`, defaults to `vi`).

```bash
envpick edit
```

### `envpick web`

Interactively select a configuration and open its `_web_url` in your browser.

```bash
envpick web

# Use with namespace
envpick web -n prod
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `EDITOR` | Editor to use for `envpick edit` (defaults to `vi`) |
