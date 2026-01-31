package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"envpick/internal/text"
	"envpick/internal/version"
)

var namespaceFlag string

var rootCmd = &cobra.Command{
	Use:     text.Text.Commands.Root.Use,
	Short:   text.Text.Commands.Root.Short,
	Long:    text.Text.Commands.Root.Long,
	Version: version.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespaceFlag, "namespace", "n", "", text.Text.Commands.Flags.Namespace)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(initCmd)
}
