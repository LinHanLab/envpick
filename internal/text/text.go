package text

// CommandText represents the text for a command (use, short description, long description).
type CommandText struct {
	Use   string
	Short string
	Long  string
}

// CommandsText contains all command-related text.
type CommandsText struct {
	Root      CommandText
	Use       CommandText
	Env       CommandText
	EnvSelect CommandText
	Edit      CommandText
	Web       CommandText
	Init      CommandText
	InitZsh   CommandText
	Flags     FlagsText
}

// FlagsText contains flag descriptions.
type FlagsText struct {
	Namespace string
}

// ErrorsText contains all error messages.
type ErrorsText struct {
	ConfigHomeDir       string
	ConfigFileNotFound  string
	ConfigFileRead      string
	ConfigFileParse     string
	ConfigNotFound      string
	ConfigNoWebURL      string
	StateFileRead       string
	StateFileParse      string
	StateEncode         string
	StateFileWrite      string
	FzfNotFound         string
	FzfFailed           string
	SelectionCancelled  string
	NoSelectionMade     string
	NoOptionsAvailable  string
	NoConfigurations    string
	NoConfigurationsUse string
	BrowserOpenFailed   string
	UnsupportedPlatform string
}

// MessagesText contains informational messages.
type MessagesText struct {
	SwitchedToConfig   string
	SwitchedToConfigNS string
	OpenedURL          string
}

// FormatsText contains formatting strings.
type FormatsText struct {
	ErrorPrefix     string
	ActiveIndicator string
	ExportStatement string
	PromptSuffix    string
}

// PromptsText contains interactive prompts.
type PromptsText struct {
	SelectConfiguration string
	SelectWebURL        string
}

// TextData contains all user-facing text for the envpick application.
// This centralizes help text, error messages, informational messages,
// display formats, and interactive prompts in a single location.
type TextData struct {
	Commands CommandsText
	Errors   ErrorsText
	Messages MessagesText
	Formats  FormatsText
	Prompts  PromptsText
}

// Text is the global instance containing all application text.
var Text = TextData{
	Commands: CommandsText{
		Root: CommandText{
			Use:   "envpick",
			Short: "Manage multiple environment variable configurations",
			Long:  `Manage multiple environment variable configurations through a simple config file and interactive commands.`,
		},
		Use: CommandText{
			Use:   "use",
			Short: "Switch configuration persistently",
			Long:  `Select a configuration to persist across new terminal sessions.`,
		},
		Env: CommandText{
			Use:   "env",
			Short: "Output current config as exports",
			Long: `Output the current configuration's environment variables as shell export statements.

Usage in shell profile (.zshrc, .bashrc):
  eval "$(envpick env)"`,
		},
		EnvSelect: CommandText{
			Use:   "select [config-name]",
			Short: "Select a configuration and output its export statements",
			Long: `Output exports for a configuration without persisting.
Prompts interactively if config-name is omitted.

Usage:
  eval "$(envpick env select myconfig)"
  eval "$(envpick env select)"`,
		},
		Edit: CommandText{
			Use:   "edit",
			Short: "Edit the configuration file",
			Long:  `Open config file in $EDITOR (default: vi).`,
		},
		Web: CommandText{
			Use:   "web",
			Short: "Open config web URL",
			Long:  `Select a configuration and open its web URL.`,
		},
		Init: CommandText{
			Use:   "init",
			Short: "Generate shell configuration for envpick",
			Long: `Generate shell integration config (completion, auto-loading, helpers).

Usage:
  eval "$(envpick init zsh)"  # Add to ~/.zshrc
  source ~/.zshrc              # Reload shell`,
		},
		InitZsh: CommandText{
			Use:   "zsh",
			Short: "Generate zsh configuration",
			Long: `Generate zsh integration config.

Add to ~/.zshrc:
  eval "$(envpick init zsh)"

Sets up completion, auto-loading, and 'ep' helper:
  ep use [flags]        - Persistent selection
  ep tmp [flags] [name] - Temporary selection
  ep <other>            - Pass through to envpick

Reload shell after adding:
  source ~/.zshrc`,
		},
		Flags: FlagsText{
			Namespace: "filter configurations by namespace (e.g., 'db' for db.local, db.prod)",
		},
	},
	Errors: ErrorsText{
		ConfigHomeDir: "failed to get home directory: %w",
		ConfigFileNotFound: `config file not found: %s

Create it with your configurations:

  [personal]
  ANTHROPIC_API_KEY = "sk-ant-xxxxx"
  ANTHROPIC_MODEL = "claude-sonnet-4-5"

  [work]
  ANTHROPIC_AUTH_TOKEN = "sk-work-xxxxx"
  _web_url = "https://dashboard.company.com"

Variables with _ prefix are metadata.
Run 'envpick edit' to create the file.`,
		ConfigFileRead:      "failed to read config file: %w",
		ConfigFileParse:     "failed to parse config file: %w",
		ConfigNotFound:      "configuration %q not found",
		ConfigNoWebURL:      "configuration %q has no web URL",
		StateFileRead:       "failed to read state file: %w",
		StateFileParse:      "failed to parse state file: %w",
		StateEncode:         "failed to encode state: %w",
		StateFileWrite:      "failed to write state file: %w",
		FzfNotFound:         "fzf not found: install fzf for interactive selection",
		FzfFailed:           "fzf failed: %w",
		SelectionCancelled:  "selection cancelled",
		NoSelectionMade:     "no selection made",
		NoOptionsAvailable:  "no options available",
		NoConfigurations:    "no configurations found",
		NoConfigurationsUse: "no available configurations",
		BrowserOpenFailed:   "failed to open browser: %w",
		UnsupportedPlatform: "unsupported platform: %s",
	},
	Messages: MessagesText{
		SwitchedToConfig:   "Switched to configuration: %s\n",
		SwitchedToConfigNS: "Switched to configuration: %s (namespace: %s)\n",
		OpenedURL:          "Opened: %s\n",
	},
	Formats: FormatsText{
		ErrorPrefix:     "envpick: %v\n",
		ActiveIndicator: " [*]",
		ExportStatement: "export %s=%q",
		PromptSuffix:    " ",
	},
	Prompts: PromptsText{
		SelectConfiguration: "Select configuration:",
		SelectWebURL:        "Select configuration to open web URL:",
	},
}
