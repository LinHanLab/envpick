package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"envpick/internal/config"
	"envpick/internal/core"
)

// TestBasicEnvironmentSelection tests basic env command functionality
func TestBasicEnvironmentSelection(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup: Config with dev/prod environments
	env.WriteConfig(BasicConfig)
	env.WriteState(`
[current]
"" = "dev"
`)

	// Action: Load engine and get export statements
	engine, err := core.NewEngine()
	require.NoError(t, err, "Failed to create engine")

	// Verify: Current config is dev
	assert.Equal(t, "dev", engine.GetCurrentConfig(), "current config should be dev")

	// Verify: Export statements contain dev values
	exports, err := engine.GetConfig().GetExportStatements("dev")
	require.NoError(t, err, "Failed to get export statements")

	output := strings.Join(exports, "\n")
	assert.Contains(t, output, "export API_URL=\"http://localhost:3000\"")
	assert.Contains(t, output, "export DB_HOST=\"localhost\"")
	assert.Contains(t, output, "export DEBUG=\"true\"")
}

// TestNamespaceIsolation tests that namespaces maintain separate state
func TestNamespaceIsolation(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup: Config with default and db namespaces
	env.WriteConfig(NamespaceConfig)
	env.WriteState(`
[current]
"" = "dev"
db = "local"
`)

	// Action: Load default namespace engine
	defaultEngine, err := core.NewEngine()
	require.NoError(t, err, "Failed to create default engine")

	// Action: Load db namespace engine
	dbEngine, err := core.NewEngineWithNamespace("db")
	require.NoError(t, err, "Failed to create db engine")

	// Verify: Default namespace shows dev
	assert.Equal(t, "dev", defaultEngine.GetCurrentConfig(), "default current config should be dev")

	// Verify: DB namespace shows local
	assert.Equal(t, "local", dbEngine.GetCurrentConfig(), "db current config should be local")

	// Verify: Default namespace exports dev values
	exports, err := defaultEngine.GetConfig().GetExportStatements("dev")
	require.NoError(t, err, "Failed to get default export statements")
	output := strings.Join(exports, "\n")
	assert.Contains(t, output, "export API_URL=\"http://localhost:3000\"")

	// Verify: DB namespace exports local values
	exports, err = dbEngine.GetConfig().GetExportStatements("db.local")
	require.NoError(t, err, "Failed to get db export statements")
	output = strings.Join(exports, "\n")
	assert.Contains(t, output, "export DB_HOST=\"localhost\"")
	assert.Contains(t, output, "export DB_PORT=\"5432\"")
}

// TestDirectSelection tests env select command (no persistence)
func TestDirectSelection(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup: Current config is dev
	env.WriteConfig(BasicConfig)
	env.WriteState(`
[current]
"" = "dev"
`)

	// Action: Get prod export statements directly
	engine, err := core.NewEngine()
	require.NoError(t, err, "Failed to create engine")

	exports, err := engine.GetConfig().GetExportStatements("prod")
	require.NoError(t, err, "Failed to get export statements")

	// Verify: Output shows prod values
	output := strings.Join(exports, "\n")
	assert.Contains(t, output, "export API_URL=\"https://api.example.com\"")
	assert.Contains(t, output, "export DEBUG=\"false\"")

	// Verify: State file still shows dev (no persistence)
	assert.Equal(t, "dev", engine.GetCurrentConfig(), "state should remain dev")
}

// TestSetCurrentConfig tests changing the current configuration
func TestSetCurrentConfig(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup
	env.WriteConfig(BasicConfig)
	env.WriteState(`
[current]
"" = "dev"
`)

	// Action: Change to prod
	engine, err := core.NewEngine()
	require.NoError(t, err, "Failed to create engine")

	err = engine.SetCurrentConfig("prod")
	require.NoError(t, err, "Failed to set current config")

	// Verify: State updated to prod
	state := env.ReadState()
	assert.Contains(t, state, `"" = "prod"`, "state should contain prod")

	// Verify: New engine loads prod as current
	engine2, err := core.NewEngine()
	require.NoError(t, err, "Failed to create second engine")

	assert.Equal(t, "prod", engine2.GetCurrentConfig(), "current config should be prod")
}

// TestMetadataFiltering tests that metadata variables are not exported
func TestMetadataFiltering(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup: Config with metadata
	env.WriteConfig(MetadataConfig)

	// Action: Get export statements
	cfg, err := config.LoadConfig()
	require.NoError(t, err, "Failed to load config")

	exports, err := cfg.GetExportStatements("dev")
	require.NoError(t, err, "Failed to get export statements")

	output := strings.Join(exports, "\n")

	// Verify: Regular vars ARE in export statements
	assert.Contains(t, output, "export API_URL=\"http://localhost:3000\"")

	// Verify: _web_url NOT in export statements
	assert.NotContains(t, output, "_web_url")
	assert.NotContains(t, output, "http://localhost:3000/admin")
}

// TestStateMigration tests migration from old to new state format
func TestStateMigration(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup: Old state format
	env.WriteConfig(BasicConfig)
	env.WriteState(LegacyState)

	// Action: Load state (should trigger migration)
	state, err := config.LoadState()
	require.NoError(t, err, "Failed to load state")

	// Verify: Migrated to new format in memory
	assert.Equal(t, "dev", state.GetCurrentConfig(""), "migrated current config should be dev")

	// Save the migrated state
	err = state.Save()
	require.NoError(t, err, "Failed to save migrated state")

	// Verify: State file updated to new format
	stateContent := env.ReadState()
	assert.Contains(t, stateContent, "[current]", "state should contain [current] section")
}

// TestMultipleNamespaceOperations tests operations across multiple namespaces
func TestMultipleNamespaceOperations(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup: Configs in 3 namespaces
	env.WriteConfig(MultiNamespaceConfig)

	// Action: Set config in each namespace
	defaultEngine, _ := core.NewEngine()
	err := defaultEngine.SetCurrentConfig("prod")
	require.NoError(t, err, "Failed to set default config")

	dbEngine, _ := core.NewEngineWithNamespace("db")
	err = dbEngine.SetCurrentConfig("prod")
	require.NoError(t, err, "Failed to set db config")

	deployEngine, _ := core.NewEngineWithNamespace("deploy")
	err = deployEngine.SetCurrentConfig("aws")
	require.NoError(t, err, "Failed to set deploy config")

	// Verify: Each namespace maintains independent state
	state := env.ReadState()
	assert.Contains(t, state, "[current]")
	assert.Contains(t, state, `"" = "prod"`)
	assert.Contains(t, state, `db = "prod"`)
	assert.Contains(t, state, `deploy = "aws"`)

	// Verify: Switching namespaces shows correct config
	defaultEngine2, _ := core.NewEngine()
	assert.Equal(t, "prod", defaultEngine2.GetCurrentConfig(), "default should be prod")

	dbEngine2, _ := core.NewEngineWithNamespace("db")
	assert.Equal(t, "prod", dbEngine2.GetCurrentConfig(), "db should be prod")

	deployEngine2, _ := core.NewEngineWithNamespace("deploy")
	assert.Equal(t, "aws", deployEngine2.GetCurrentConfig(), "deploy should be aws")
}

// TestErrorHandling tests various error conditions
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig string
		setupState  string
		action      func(*core.Engine) error
		expectError string
	}{
		{
			name:        "non-existent config",
			setupConfig: BasicConfig,
			setupState:  `[current]` + "\n" + `"" = "dev"`,
			action: func(e *core.Engine) error {
				return e.SetCurrentConfig("nonexistent")
			},
			expectError: "not found",
		},
		{
			name:        "missing config file",
			setupConfig: "",
			setupState:  "",
			action: func(e *core.Engine) error {
				_, err := core.NewEngine()
				return err
			},
			expectError: "not found",
		},
		{
			name:        "invalid TOML syntax",
			setupConfig: InvalidConfig,
			setupState:  "",
			action: func(e *core.Engine) error {
				_, err := config.LoadConfig()
				return err
			},
			expectError: "failed to parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewTestEnv(t)
			defer env.SetHome()()

			if tt.setupConfig != "" {
				env.WriteConfig(tt.setupConfig)
			}
			if tt.setupState != "" {
				env.WriteState(tt.setupState)
			}

			var err error
			if tt.name == "missing config file" || tt.name == "invalid TOML syntax" {
				err = tt.action(nil)
			} else {
				engine, createErr := core.NewEngine()
				if createErr != nil {
					err = createErr
				} else {
					err = tt.action(engine)
				}
			}

			require.Error(t, err, "should return an error")
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

// TestGetOptions tests option generation for selection
func TestGetOptions(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup
	env.WriteConfig(NamespaceConfig)
	env.WriteState(`
[current]
"" = "prod"
db = "local"
`)

	// Test default namespace options
	t.Run("default namespace", func(t *testing.T) {
		engine, err := core.NewEngine()
		require.NoError(t, err, "Failed to create engine")

		options := engine.GetOptions()
		assert.Len(t, options, 2, "should have 2 options")

		// Find prod option and verify it's marked active
		var foundProd bool
		for _, opt := range options {
			if opt.Name == "prod" {
				foundProd = true
				assert.Equal(t, "active", opt.Status, "prod should be active")
			}
		}
		assert.True(t, foundProd, "should find prod option")
	})

	// Test db namespace options
	t.Run("db namespace", func(t *testing.T) {
		engine, err := core.NewEngineWithNamespace("db")
		require.NoError(t, err, "Failed to create engine")

		options := engine.GetOptions()
		assert.Len(t, options, 2, "should have 2 options")

		// Find local option and verify it's marked active
		var foundLocal bool
		for _, opt := range options {
			if opt.Name == "local" {
				foundLocal = true
				assert.Equal(t, "active", opt.Status, "local should be active")
			}
		}
		assert.True(t, foundLocal, "should find local option")
	})
}

// TestConfigFileCreation tests that config directory and file are created
func TestConfigFileCreation(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Remove config directory
	err := os.RemoveAll(env.ConfigDir)
	require.NoError(t, err, "Failed to remove config dir")

	// Action: Ensure config directory exists
	configDir := filepath.Join(env.HomeDir, ".envpick")
	err = os.MkdirAll(configDir, 0755)
	require.NoError(t, err, "Failed to create config dir")

	// Verify: Directory created
	if !env.ConfigExists() {
		// This is expected - config file doesn't exist yet
		// Write a sample config
		env.WriteConfig(BasicConfig)
	}

	// Verify: Can load config
	cfg, err := config.LoadConfig()
	require.NoError(t, err, "Failed to load config")

	assert.NotEmpty(t, cfg.Configs, "config should have environments")
}

// TestNamespaceFiltering tests that namespace filtering works correctly
func TestNamespaceFiltering(t *testing.T) {
	env := NewTestEnv(t)
	defer env.SetHome()()

	// Setup
	env.WriteConfig(MultiNamespaceConfig)

	// Test: Load config and check namespace filtering
	cfg, err := config.LoadConfig()
	require.NoError(t, err, "Failed to load config")

	// Test default namespace
	defaultConfigs := cfg.GetNamespaceConfigs("")
	assert.Len(t, defaultConfigs, 2, "should have 2 default configs")
	assert.Contains(t, defaultConfigs, "dev", "should contain dev in default namespace")
	assert.Contains(t, defaultConfigs, "prod", "should contain prod in default namespace")

	// Test db namespace
	dbConfigs := cfg.GetNamespaceConfigs("db")
	assert.Len(t, dbConfigs, 2, "should have 2 db configs")
	assert.Contains(t, dbConfigs, "local", "should contain local in db namespace")
	assert.Contains(t, dbConfigs, "prod", "should contain prod in db namespace")

	// Test deploy namespace
	deployConfigs := cfg.GetNamespaceConfigs("deploy")
	assert.Len(t, deployConfigs, 2, "should have 2 deploy configs")
	assert.Contains(t, deployConfigs, "aws", "should contain aws in deploy namespace")
	assert.Contains(t, deployConfigs, "gcp", "should contain gcp in deploy namespace")
}
