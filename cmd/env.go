package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envpick/internal/core"
	"envpick/internal/selector"
	"envpick/internal/text"
)

var envCmd = &cobra.Command{
	Use:   text.Text.Commands.Env.Use,
	Short: text.Text.Commands.Env.Short,
	Long:  text.Text.Commands.Env.Long,
	Run: func(cmd *cobra.Command, args []string) {
		engine, err := core.NewEngineWithNamespace(namespaceFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, text.Text.Formats.ErrorPrefix, err)
			return
		}

		// Get current config (full name with namespace)
		configName := engine.GetCurrentConfigFull()

		exports, err := engine.GetConfig().GetExportStatements(configName)
		if err != nil {
			fmt.Fprintf(os.Stderr, text.Text.Formats.ErrorPrefix, err)
			return
		}

		fmt.Println(strings.Join(exports, "\n"))
	},
}

var envSelectCmd = &cobra.Command{
	Use:   text.Text.Commands.EnvSelect.Use,
	Short: text.Text.Commands.EnvSelect.Short,
	Long:  text.Text.Commands.EnvSelect.Long,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		engine, err := core.NewEngineWithNamespace(namespaceFlag)
		if err != nil {
			return err
		}

		var selected string
		if len(args) > 0 {
			// Direct config selection (short form)
			selected = args[0]
			// Build full name and validate config exists
			var fullName string
			if engine.GetNamespace() != "" {
				fullName = engine.GetNamespace() + "." + selected
			} else {
				fullName = selected
			}
			if _, ok := engine.GetConfig().Configs[fullName]; !ok {
				return fmt.Errorf(text.Text.Errors.ConfigNotFound, selected)
			}
			selected = fullName
		} else {
			// Interactive selection
			options := engine.GetOptions()
			if len(options) == 0 {
				return errors.New(text.Text.Errors.NoConfigurationsUse)
			}

			shortName, err := selector.Select(options, text.Text.Prompts.SelectConfiguration)
			if err != nil {
				return err
			}
			// Build full name from short name
			if engine.GetNamespace() != "" {
				selected = engine.GetNamespace() + "." + shortName
			} else {
				selected = shortName
			}
		}

		exports, err := engine.GetConfig().GetExportStatements(selected)
		if err != nil {
			return err
		}

		fmt.Println(strings.Join(exports, "\n"))
		return nil
	},
}

func init() {
	envCmd.AddCommand(envSelectCmd)
}
