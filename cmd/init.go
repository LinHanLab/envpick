package cmd

import (
	"fmt"

	"envpick/internal/text"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   text.Text.Commands.Init.Use,
	Short: text.Text.Commands.Init.Short,
	Long:  text.Text.Commands.Init.Long,
}

// initZshCmd represents the init zsh subcommand
var initZshCmd = &cobra.Command{
	Use:   text.Text.Commands.InitZsh.Use,
	Short: text.Text.Commands.InitZsh.Short,
	Long:  text.Text.Commands.InitZsh.Long,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(zshConfig)
	},
}

const zshConfig = `# envpick shell integration
if command -v envpick >/dev/null 2>&1; then
    # Load shell completion if completion system is initialized
    if (( $+functions[compdef] )); then
        source <(envpick completion zsh)
    fi

    # Load persisted environment on shell startup
    eval "$(envpick env 2>/dev/null)"

    # Helper function for envpick operations
    ep() {
        case "$1" in
            use)
                # Interactive selection with persistence
                shift
                if envpick use "$@"; then
                    eval "$(envpick env)"
                fi
                ;;
            tmp)
                # Temporary selection (no persistence)
                shift
                eval "$(envpick env select "$@")"
                ;;
            *)
                # Pass through all other commands
                envpick "$@"
                ;;
        esac
    }
fi
`

func init() {
	initCmd.AddCommand(initZshCmd)
}
