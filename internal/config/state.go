package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"envpick/internal/text"
)

// State represents the state file
type State struct {
	// Map of namespace -> current config name (short form, without namespace prefix)
	Current map[string]string `toml:"current"`

	// Legacy field for backward compatibility (deprecated)
	CurrentConfig string `toml:"current_config,omitempty"`
}

// GetStatePath returns the path to state.toml
// This is a variable to allow overriding in tests
var GetStatePath = func() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.toml"), nil
}

// LoadState loads the state from state.toml
func LoadState() (*State, error) {
	statePath, err := GetStatePath()
	if err != nil {
		return nil, err
	}

	state := &State{
		Current: make(map[string]string),
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}
		return nil, fmt.Errorf(text.Text.Errors.StateFileRead, err)
	}

	if _, err := toml.Decode(string(data), state); err != nil {
		return nil, fmt.Errorf(text.Text.Errors.StateFileParse, err)
	}

	// Migrate from old format if needed
	if state.CurrentConfig != "" && len(state.Current) == 0 {
		// Parse the old current_config to determine namespace
		ns, cfg := ParseConfigName(state.CurrentConfig)
		state.Current[ns] = cfg
		state.CurrentConfig = "" // Clear legacy field
	}

	return state, nil
}

// Save saves the state to state.toml
func (s *State) Save() error {
	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	statePath, err := GetStatePath()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(s); err != nil {
		return fmt.Errorf(text.Text.Errors.StateEncode, err)
	}

	if err := os.WriteFile(statePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf(text.Text.Errors.StateFileWrite, err)
	}

	return nil
}

// GetCurrentConfig returns the current config for the given namespace.
// Returns empty string if no config is set for the namespace.
func (s *State) GetCurrentConfig(namespace string) string {
	if s.Current == nil {
		return ""
	}
	return s.Current[namespace]
}

// SetCurrentConfig sets the current config for the given namespace.
func (s *State) SetCurrentConfig(namespace, config string) {
	if s.Current == nil {
		s.Current = make(map[string]string)
	}
	s.Current[namespace] = config
}

// CreateDefaultState creates a default state.toml if it doesn't exist
func CreateDefaultState(defaultConfig string) error {
	statePath, err := GetStatePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(statePath); err == nil {
		return nil // File already exists
	}

	state := &State{
		Current: make(map[string]string),
	}

	// Parse the default config to determine namespace
	ns, cfg := ParseConfigName(defaultConfig)
	state.SetCurrentConfig(ns, cfg)

	return state.Save()
}
