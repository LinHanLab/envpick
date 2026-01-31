package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateGetSetCurrentConfig(t *testing.T) {
	state := &State{
		Current: make(map[string]string),
	}

	// Test setting and getting configs for different namespaces
	state.SetCurrentConfig("", "dev")
	state.SetCurrentConfig("db", "local")
	state.SetCurrentConfig("deploy", "aws")

	assert.Equal(t, "dev", state.GetCurrentConfig(""), "default namespace should be dev")
	assert.Equal(t, "local", state.GetCurrentConfig("db"), "db namespace should be local")
	assert.Equal(t, "aws", state.GetCurrentConfig("deploy"), "deploy namespace should be aws")

	// Test non-existent namespace
	assert.Empty(t, state.GetCurrentConfig("nonexistent"), "non-existent namespace should return empty string")
}

func TestStateMigration(t *testing.T) {
	// Create a temporary directory for test state file
	tmpDir := t.TempDir()
	oldStatePath := filepath.Join(tmpDir, "state.toml")

	// Write old format state file
	oldContent := `current_config = "prod"`
	err := os.WriteFile(oldStatePath, []byte(oldContent), 0644)
	require.NoError(t, err, "Failed to write test state file")

	// Override GetStatePath for testing
	originalGetStatePath := GetStatePath
	GetStatePath = func() (string, error) {
		return oldStatePath, nil
	}
	defer func() { GetStatePath = originalGetStatePath }()

	// Load state and verify migration
	state, err := LoadState()
	require.NoError(t, err, "LoadState should succeed")

	assert.Equal(t, "prod", state.GetCurrentConfig(""), "after migration, default namespace should be prod")

	// Verify legacy field is cleared
	assert.Empty(t, state.CurrentConfig, "after migration, CurrentConfig should be empty")
}

func TestStateMigrationWithNamespace(t *testing.T) {
	// Create a temporary directory for test state file
	tmpDir := t.TempDir()
	oldStatePath := filepath.Join(tmpDir, "state.toml")

	// Write old format state file with namespaced config
	oldContent := `current_config = "db.local"`
	err := os.WriteFile(oldStatePath, []byte(oldContent), 0644)
	require.NoError(t, err, "Failed to write test state file")

	// Override GetStatePath for testing
	originalGetStatePath := GetStatePath
	GetStatePath = func() (string, error) {
		return oldStatePath, nil
	}
	defer func() { GetStatePath = originalGetStatePath }()

	// Load state and verify migration
	state, err := LoadState()
	require.NoError(t, err, "LoadState should succeed")

	assert.Equal(t, "local", state.GetCurrentConfig("db"), "after migration, db namespace should be local")
}

func TestStateSaveLoad(t *testing.T) {
	// Create a temporary directory for test state file
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.toml")

	// Override GetStatePath for testing
	originalGetStatePath := GetStatePath
	GetStatePath = func() (string, error) {
		return statePath, nil
	}
	defer func() { GetStatePath = originalGetStatePath }()

	// Create and save state
	state := &State{
		Current: map[string]string{
			"":       "prod",
			"db":     "staging",
			"deploy": "gcp",
		},
	}

	err := state.Save()
	require.NoError(t, err, "Save should succeed")

	// Load state and verify
	loadedState, err := LoadState()
	require.NoError(t, err, "LoadState should succeed")

	assert.Equal(t, "prod", loadedState.GetCurrentConfig(""), "default namespace should be prod")
	assert.Equal(t, "staging", loadedState.GetCurrentConfig("db"), "db namespace should be staging")
	assert.Equal(t, "gcp", loadedState.GetCurrentConfig("deploy"), "deploy namespace should be gcp")
}
