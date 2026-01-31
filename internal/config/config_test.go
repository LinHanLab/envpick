package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigName(t *testing.T) {
	tests := []struct {
		name              string
		fullName          string
		expectedNamespace string
		expectedConfig    string
	}{
		{
			name:              "default namespace",
			fullName:          "dev",
			expectedNamespace: "",
			expectedConfig:    "dev",
		},
		{
			name:              "named namespace",
			fullName:          "db.local",
			expectedNamespace: "db",
			expectedConfig:    "local",
		},
		{
			name:              "multi-dot config",
			fullName:          "db.prod.primary",
			expectedNamespace: "db",
			expectedConfig:    "prod.primary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns, cfg := ParseConfigName(tt.fullName)
			assert.Equal(t, tt.expectedNamespace, ns, "namespace should match")
			assert.Equal(t, tt.expectedConfig, cfg, "config should match")
		})
	}
}

func TestBuildConfigName(t *testing.T) {
	tests := []struct {
		name         string
		namespace    string
		config       string
		expectedName string
	}{
		{
			name:         "default namespace",
			namespace:    "",
			config:       "dev",
			expectedName: "dev",
		},
		{
			name:         "named namespace",
			namespace:    "db",
			config:       "local",
			expectedName: "db.local",
		},
		{
			name:         "multi-part config",
			namespace:    "db",
			config:       "prod.primary",
			expectedName: "db.prod.primary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildConfigName(tt.namespace, tt.config)
			assert.Equal(t, tt.expectedName, result, "config name should match")
		})
	}
}

func TestGetNamespaceConfigs(t *testing.T) {
	cfg := &Config{
		Configs: map[string]map[string]string{
			"dev":      {"API_KEY": "dev-key"},
			"prod":     {"API_KEY": "prod-key"},
			"db.local": {"DB_HOST": "localhost"},
			"db.prod":  {"DB_HOST": "prod.db"},
		},
	}

	tests := []struct {
		name              string
		namespace         string
		expectedConfigLen int
		expectedConfigs   []string
	}{
		{
			name:              "default namespace",
			namespace:         "",
			expectedConfigLen: 2,
			expectedConfigs:   []string{"dev", "prod"},
		},
		{
			name:              "db namespace",
			namespace:         "db",
			expectedConfigLen: 2,
			expectedConfigs:   []string{"local", "prod"},
		},
		{
			name:              "non-existent namespace",
			namespace:         "deploy",
			expectedConfigLen: 0,
			expectedConfigs:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GetNamespaceConfigs(tt.namespace)
			assert.Len(t, result, tt.expectedConfigLen, "should have expected number of configs")
			for _, expectedConfig := range tt.expectedConfigs {
				assert.Contains(t, result, expectedConfig, "should contain expected config")
			}
		})
	}
}

func TestGetNamespaces(t *testing.T) {
	cfg := &Config{
		Configs: map[string]map[string]string{
			"dev":        {"API_KEY": "dev-key"},
			"prod":       {"API_KEY": "prod-key"},
			"db.local":   {"DB_HOST": "localhost"},
			"db.prod":    {"DB_HOST": "prod.db"},
			"deploy.aws": {"AWS_REGION": "us-east-1"},
		},
	}

	namespaces := cfg.GetNamespaces()
	expectedNamespaces := []string{"", "db", "deploy"}

	assert.Len(t, namespaces, len(expectedNamespaces), "should have expected number of namespaces")
	for _, ns := range expectedNamespaces {
		assert.Contains(t, namespaces, ns, "should contain expected namespace")
	}
}

func TestExtractConfigs(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		expectedKeys   []string
		unexpectedKeys []string
	}{
		{
			name: "flat configs",
			input: map[string]interface{}{
				"dev": map[string]interface{}{
					"API_KEY": "dev-key",
				},
				"prod": map[string]interface{}{
					"API_KEY": "prod-key",
				},
			},
			expectedKeys:   []string{"dev", "prod"},
			unexpectedKeys: []string{},
		},
		{
			name: "nested configs",
			input: map[string]interface{}{
				"db": map[string]interface{}{
					"local": map[string]interface{}{
						"DB_HOST": "localhost",
					},
					"prod": map[string]interface{}{
						"DB_HOST": "prod.db",
					},
				},
			},
			expectedKeys:   []string{"db.local", "db.prod"},
			unexpectedKeys: []string{"db"},
		},
		{
			name: "mixed flat and nested",
			input: map[string]interface{}{
				"dev": map[string]interface{}{
					"API_KEY": "dev-key",
				},
				"db": map[string]interface{}{
					"local": map[string]interface{}{
						"DB_HOST": "localhost",
					},
				},
			},
			expectedKeys:   []string{"dev", "db.local"},
			unexpectedKeys: []string{"db"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configs := make(map[string]map[string]string)
			extractConfigs(configs, tt.input, "")

			for _, key := range tt.expectedKeys {
				assert.Contains(t, configs, key, "should contain expected key")
			}

			for _, key := range tt.unexpectedKeys {
				assert.NotContains(t, configs, key, "should not contain unexpected key")
			}
		})
	}
}
