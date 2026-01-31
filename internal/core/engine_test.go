package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"envpick/internal/config"
)

func TestEngineNamespaceFiltering(t *testing.T) {
	// Create a mock config with multiple namespaces
	cfg := &config.Config{
		Configs: map[string]map[string]string{
			"dev":      {"API_KEY": "dev-key"},
			"prod":     {"API_KEY": "prod-key"},
			"db.local": {"DB_HOST": "localhost"},
			"db.prod":  {"DB_HOST": "prod.db"},
		},
	}

	state := &config.State{
		Current: map[string]string{
			"":   "dev",
			"db": "local",
		},
	}

	// Test default namespace
	engine := &Engine{
		config:    cfg,
		state:     state,
		namespace: "",
	}

	options := engine.GetOptions()
	assert.Len(t, options, 2, "default namespace should have 2 options")

	// Verify options are from default namespace
	optionNames := make(map[string]bool)
	for _, opt := range options {
		optionNames[opt.Name] = true
	}
	assert.True(t, optionNames["dev"], "should contain dev")
	assert.True(t, optionNames["prod"], "should contain prod")

	// Test db namespace
	engine.namespace = "db"
	options = engine.GetOptions()
	assert.Len(t, options, 2, "db namespace should have 2 options")

	// Verify options are from db namespace
	optionNames = make(map[string]bool)
	for _, opt := range options {
		optionNames[opt.Name] = true
	}
	assert.True(t, optionNames["local"], "should contain local")
	assert.True(t, optionNames["prod"], "should contain prod")
}

func TestEngineGetCurrentConfig(t *testing.T) {
	cfg := &config.Config{
		Configs: map[string]map[string]string{
			"dev":      {"API_KEY": "dev-key"},
			"db.local": {"DB_HOST": "localhost"},
		},
	}

	state := &config.State{
		Current: map[string]string{
			"":   "dev",
			"db": "local",
		},
	}

	// Test default namespace
	engine := &Engine{
		config:    cfg,
		state:     state,
		namespace: "",
	}

	assert.Equal(t, "dev", engine.GetCurrentConfig(), "current config should be dev")
	assert.Equal(t, "dev", engine.GetCurrentConfigFull(), "current config full should be dev")

	// Test db namespace
	engine.namespace = "db"
	assert.Equal(t, "local", engine.GetCurrentConfig(), "current config should be local")
	assert.Equal(t, "db.local", engine.GetCurrentConfigFull(), "current config full should be db.local")
}

func TestEngineSetCurrentConfig(t *testing.T) {
	cfg := &config.Config{
		Configs: map[string]map[string]string{
			"dev":      {"API_KEY": "dev-key"},
			"prod":     {"API_KEY": "prod-key"},
			"db.local": {"DB_HOST": "localhost"},
		},
	}

	state := &config.State{
		Current: map[string]string{
			"": "dev",
		},
	}

	// Mock GetStatePath to avoid file I/O
	originalGetStatePath := config.GetStatePath
	config.GetStatePath = func() (string, error) {
		return "/tmp/test-state.toml", nil
	}
	defer func() { config.GetStatePath = originalGetStatePath }()

	// Test setting config in default namespace
	engine := &Engine{
		config:    cfg,
		state:     state,
		namespace: "",
	}

	err := engine.SetCurrentConfig("prod")
	require.NoError(t, err, "SetCurrentConfig should succeed")

	assert.Equal(t, "prod", engine.GetCurrentConfig(), "current config should be prod")

	// Test setting config in db namespace
	engine.namespace = "db"
	err = engine.SetCurrentConfig("local")
	require.NoError(t, err, "SetCurrentConfig should succeed")

	assert.Equal(t, "local", engine.GetCurrentConfig(), "current config should be local")

	// Test setting non-existent config
	err = engine.SetCurrentConfig("nonexistent")
	assert.Error(t, err, "SetCurrentConfig with non-existent config should fail")
}

func TestEngineGetNamespace(t *testing.T) {
	engine := &Engine{
		namespace: "db",
	}

	assert.Equal(t, "db", engine.GetNamespace(), "namespace should be db")
}
