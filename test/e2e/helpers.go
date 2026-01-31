package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnv represents an isolated test environment
type TestEnv struct {
	HomeDir    string
	ConfigDir  string
	ConfigFile string
	StateFile  string
	T          *testing.T
}

// NewTestEnv creates an isolated test environment
func NewTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	homeDir := t.TempDir()
	configDir := filepath.Join(homeDir, ".envpick")

	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err, "Failed to create config dir")

	return &TestEnv{
		HomeDir:    homeDir,
		ConfigDir:  configDir,
		ConfigFile: filepath.Join(configDir, "config.toml"),
		StateFile:  filepath.Join(configDir, "state.toml"),
		T:          t,
	}
}

// WriteConfig writes a config file
func (e *TestEnv) WriteConfig(content string) {
	e.T.Helper()
	err := os.WriteFile(e.ConfigFile, []byte(content), 0644)
	require.NoError(e.T, err, "Failed to write config")
}

// WriteState writes a state file
func (e *TestEnv) WriteState(content string) {
	e.T.Helper()
	err := os.WriteFile(e.StateFile, []byte(content), 0644)
	require.NoError(e.T, err, "Failed to write state")
}

// ReadState reads the state file
func (e *TestEnv) ReadState() string {
	e.T.Helper()
	content, err := os.ReadFile(e.StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		require.NoError(e.T, err, "Failed to read state")
	}
	return string(content)
}

// ConfigExists checks if config file exists
func (e *TestEnv) ConfigExists() bool {
	_, err := os.Stat(e.ConfigFile)
	return err == nil
}

// StateExists checks if state file exists
func (e *TestEnv) StateExists() bool {
	_, err := os.Stat(e.StateFile)
	return err == nil
}

// SetHome sets HOME environment variable to test directory
func (e *TestEnv) SetHome() func() {
	e.T.Helper()
	oldHome := os.Getenv("HOME")
	err := os.Setenv("HOME", e.HomeDir)
	require.NoError(e.T, err, "Failed to set HOME")
	return func() {
		if err := os.Setenv("HOME", oldHome); err != nil {
			e.T.Logf("Warning: Failed to restore HOME: %v", err)
		}
	}
}

// AssertStateContains checks if state file contains expected string
func (e *TestEnv) AssertStateContains(expected string) {
	e.T.Helper()
	state := e.ReadState()
	assert.Contains(e.T, state, expected, "State should contain expected string")
}
