# envpick

[English](README.md) | 简体中文

一个通过 fzf 交互式选择管理多个环境变量配置的 CLI 工具。

<img src="docs/vhs/demo.gif" width="500" alt="Demo">

**[安装指南](docs/installation.zh-CN.md)** | 下方快速开始

## 快速开始

本指南帮助您通过 3 个简单步骤开始使用 envpick。

### 1. 设置你的 shell (zsh)

在 `~/.zshrc` 中添加:

```bash
# 初始化 zsh 补全系统（如果尚未完成）
autoload -U compinit; compinit

# 初始化 envpick
eval "$(envpick init zsh)"
```

重新加载你的 shell:

```bash
source ~/.zshrc
```

这将启用补全功能并添加 `ep` 快捷命令。

### 2. 创建你的配置

envpick 使用 `~/.envpick/config.toml` 存储你的环境配置:

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

以 `_` 开头的变量是元数据（例如，`_web_url` 用于 web URL）。

随时使用以下命令编辑配置:

```bash
ep edit
```

### 3. 在配置之间切换

使用交互式选择在你的配置之间切换:

```bash
ep use
```

这将打开 fzf 进行交互式选择。你的选择将在新的终端会话中保持。

### 使用命名空间（高级）

当你有多组相关的配置（例如，数据库、API）时，使用命名空间:

```toml
[db.local]
DATABASE_URL = "postgres://localhost/myapp"

[db.staging]
DATABASE_URL = "postgres://staging.example.com/myapp"

[db.prod]
DATABASE_URL = "postgres://prod.example.com/myapp"
```

在命名空间内切换:

```bash
ep use -n db
```

每个命名空间独立维护自己的状态。

### 临时配置（一次性使用）

对于不需要持久化的临时配置更改:

```bash
# 交互式选择（仅当前终端）
ep tmp

# 直接选择
ep tmp work

# 使用命名空间
ep tmp -n db staging
```

当你需要为单个终端会话使用不同的环境变量而不改变持久配置时，使用 `ep tmp`。

## 功能特性

- 使用 fzf 进行交互式配置切换
- 跨终端会话的持久状态
- 支持命名空间以组织配置
- Web URL 启动器: `envpick web`
- 临时配置选择: `envpick env select`
- 通过 `ep` 辅助函数进行 shell 集成

完整的命令文档请参考: `envpick --help`
