package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"envpick/internal/core"
	"envpick/internal/selector"
	"envpick/internal/text"
)

var webCmd = &cobra.Command{
	Use:   text.Text.Commands.Web.Use,
	Short: text.Text.Commands.Web.Short,
	Long:  text.Text.Commands.Web.Long,
	RunE: func(cmd *cobra.Command, args []string) error {
		engine, err := core.NewEngineWithNamespace(namespaceFlag)
		if err != nil {
			return err
		}

		options := engine.GetOptions()
		if len(options) == 0 {
			return errors.New(text.Text.Errors.NoConfigurations)
		}

		selected, err := selector.Select(options, text.Text.Prompts.SelectWebURL)
		if err != nil {
			return err
		}

		// Build full config name
		fullName := selected
		if engine.GetNamespace() != "" {
			fullName = engine.GetNamespace() + "." + selected
		}

		url, err := engine.GetConfig().GetWebURL(fullName)
		if err != nil {
			return err
		}

		if err := openBrowser(url); err != nil {
			return fmt.Errorf(text.Text.Errors.BrowserOpenFailed, err)
		}

		fmt.Printf(text.Text.Messages.OpenedURL, url)
		return nil
	},
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf(text.Text.Errors.UnsupportedPlatform, runtime.GOOS)
	}

	return cmd.Start()
}
