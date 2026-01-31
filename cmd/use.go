package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"envpick/internal/core"
	"envpick/internal/selector"
	"envpick/internal/text"
)

var useCmd = &cobra.Command{
	Use:   text.Text.Commands.Use.Use,
	Short: text.Text.Commands.Use.Short,
	Long:  text.Text.Commands.Use.Long,
	RunE: func(cmd *cobra.Command, args []string) error {
		engine, err := core.NewEngineWithNamespace(namespaceFlag)
		if err != nil {
			return err
		}

		options := engine.GetOptions()
		if len(options) == 0 {
			return errors.New(text.Text.Errors.NoConfigurationsUse)
		}

		selected, err := selector.Select(options, text.Text.Prompts.SelectConfiguration)
		if err != nil {
			return err
		}

		if err := engine.SetCurrentConfig(selected); err != nil {
			return err
		}

		// Show namespace in output if non-default
		if engine.GetNamespace() != "" {
			fmt.Printf(text.Text.Messages.SwitchedToConfigNS, selected, engine.GetNamespace())
		} else {
			fmt.Printf(text.Text.Messages.SwitchedToConfig, selected)
		}
		return nil
	},
}
