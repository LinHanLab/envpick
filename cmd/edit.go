package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"envpick/internal/config"
	"envpick/internal/text"
)

var editCmd = &cobra.Command{
	Use:   text.Text.Commands.Edit.Use,
	Short: text.Text.Commands.Edit.Short,
	Long:  text.Text.Commands.Edit.Long,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Ensure config directory exists
		if err := config.EnsureConfigDir(); err != nil {
			return err
		}

		configPath, err := config.GetConfigPath()
		if err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		editorCmd := exec.Command(editor, configPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		return editorCmd.Run()
	},
}
