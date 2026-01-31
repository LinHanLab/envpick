package e2e

import (
	"fmt"
	"os/exec"
)

// MockSelector simulates fzf selection for testing
type MockSelector struct {
	Selection string
	Err       error
	Called    bool
	Options   []string
	Prompt    string
}

func (m *MockSelector) Select(options []string, prompt string) (string, error) {
	m.Called = true
	m.Options = options
	m.Prompt = prompt
	if m.Err != nil {
		return "", m.Err
	}
	return m.Selection, nil
}

// MockURLOpener simulates browser opening for testing
type MockURLOpener struct {
	OpenedURL string
	Err       error
	Called    bool
}

func (m *MockURLOpener) Open(url string) error {
	m.Called = true
	m.OpenedURL = url
	if m.Err != nil {
		return m.Err
	}
	return nil
}

// MockFileEditor simulates editor opening for testing
type MockFileEditor struct {
	EditedFile string
	Err        error
	Called     bool
}

func (m *MockFileEditor) Edit(path string) error {
	m.Called = true
	m.EditedFile = path
	if m.Err != nil {
		return m.Err
	}
	return nil
}

// MockCommandRunner simulates command execution for testing
type MockCommandRunner struct {
	Commands []string
	Err      error
}

func (m *MockCommandRunner) Run(cmd *exec.Cmd) error {
	m.Commands = append(m.Commands, fmt.Sprintf("%s %v", cmd.Path, cmd.Args))
	return m.Err
}
