package selector

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"envpick/internal/text"
)

// Option represents a selectable option in fzf
type Option struct {
	Name   string
	Status string // "active" or empty
}

// runFzf executes fzf with the given input and prompt, returning the selected line
func runFzf(input, prompt string) (string, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", errors.New(text.Text.Errors.FzfNotFound)
	}

	cmd := exec.Command("fzf", "--prompt", prompt+text.Text.Formats.PromptSuffix, "--no-multi", "--height=40%", "--reverse")
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
			return "", errors.New(text.Text.Errors.SelectionCancelled)
		}
		return "", fmt.Errorf(text.Text.Errors.FzfFailed, err)
	}

	selected := strings.TrimSpace(string(output))
	if selected == "" {
		return "", errors.New(text.Text.Errors.NoSelectionMade)
	}

	return selected, nil
}

// Select presents options via fzf and returns the selected option name
func Select(options []Option, prompt string) (string, error) {
	if len(options) == 0 {
		return "", errors.New(text.Text.Errors.NoOptionsAvailable)
	}

	var lines []string
	for _, opt := range options {
		lines = append(lines, formatOption(opt))
	}

	selected, err := runFzf(strings.Join(lines, "\n"), prompt)
	if err != nil {
		return "", err
	}

	return extractName(selected), nil
}

// formatOption formats an option for display in fzf
func formatOption(opt Option) string {
	line := opt.Name
	if opt.Status == "active" {
		line += text.Text.Formats.ActiveIndicator
	}

	return line
}

// extractName extracts the configuration name from a formatted fzf line
func extractName(line string) string {
	parts := strings.Fields(line)
	if len(parts) > 0 {
		return parts[0]
	}
	return line
}
