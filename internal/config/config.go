package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"envpick/internal/text"
)

// Config represents the main configuration file
type Config struct {
	Configs map[string]map[string]string `toml:"-"`
}

// ConfigEntry represents a single configuration with its variables and metadata
type ConfigEntry struct {
	Vars   map[string]string
	WebURL string
}

// GetConfigDir returns the envpick configuration directory
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(text.Text.Errors.ConfigHomeDir, err)
	}
	return filepath.Join(home, ".envpick"), nil
}

// GetConfigPath returns the path to config.toml
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// LoadConfig loads the configuration from config.toml
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Read raw TOML
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(text.Text.Errors.ConfigFileNotFound, configPath)
		}
		return nil, fmt.Errorf(text.Text.Errors.ConfigFileRead, err)
	}

	// Parse into generic map first
	var raw map[string]interface{}
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return nil, fmt.Errorf(text.Text.Errors.ConfigFileParse, err)
	}

	config := &Config{
		Configs: make(map[string]map[string]string),
	}

	// Extract config sections (recursively handle nested tables)
	extractConfigs(config.Configs, raw, "")

	return config, nil
}

// extractConfigs recursively extracts configuration sections from TOML data
// prefix is used to build the full config name (e.g., "db" for nested tables)
func extractConfigs(configs map[string]map[string]string, data map[string]interface{}, prefix string) {
	for key, val := range data {
		if key == "default" {
			continue // Skip legacy default key
		}

		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if section, ok := val.(map[string]interface{}); ok {
			// Check if this is a config section (contains string values) or a namespace (contains nested maps)
			hasStringValues := false
			hasNestedMaps := false

			for _, v := range section {
				if _, ok := v.(string); ok {
					hasStringValues = true
				}
				if _, ok := v.(map[string]interface{}); ok {
					hasNestedMaps = true
				}
			}

			if hasStringValues && !hasNestedMaps {
				// This is a config section
				configs[fullKey] = make(map[string]string)
				for k, v := range section {
					if s, ok := v.(string); ok {
						configs[fullKey][k] = s
					}
				}
			} else if hasNestedMaps {
				// This is a namespace, recurse into it
				extractConfigs(configs, section, fullKey)
			}
		}
	}
}

// GetEntry returns a ConfigEntry for the given config name
func (c *Config) GetEntry(name string) (*ConfigEntry, error) {
	vars, ok := c.Configs[name]
	if !ok {
		return nil, fmt.Errorf(text.Text.Errors.ConfigNotFound, name)
	}

	entry := &ConfigEntry{
		Vars: make(map[string]string),
	}

	for k, v := range vars {
		switch k {
		case "_web_url":
			entry.WebURL = v
		default:
			if len(k) > 0 && k[0] != '_' {
				entry.Vars[k] = v
			}
		}
	}

	return entry, nil
}

// GetConfigNames returns all configuration names
func (c *Config) GetConfigNames() []string {
	var names []string
	for name := range c.Configs {
		names = append(names, name)
	}
	return names
}

// ParseConfigName splits a full config name into namespace and config parts.
// Returns ("", "dev") for "dev" and ("db", "local") for "db.local"
func ParseConfigName(fullName string) (namespace, config string) {
	parts := strings.SplitN(fullName, ".", 2)
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0], parts[1]
}

// BuildConfigName joins namespace and config into a full name.
// Returns "dev" for ("", "dev") and "db.local" for ("db", "local")
func BuildConfigName(namespace, config string) string {
	if namespace == "" {
		return config
	}
	return namespace + "." + config
}

// GetNamespaceConfigs returns all configs in a specific namespace.
// For default namespace (""), returns configs without dots.
// For named namespace, returns configs with that prefix.
func (c *Config) GetNamespaceConfigs(namespace string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for fullName, vars := range c.Configs {
		ns, configName := ParseConfigName(fullName)
		if ns == namespace {
			result[configName] = vars
		}
	}
	return result
}

// GetNamespaces returns a list of all unique namespaces in the config.
func (c *Config) GetNamespaces() []string {
	namespaceSet := make(map[string]bool)
	for fullName := range c.Configs {
		ns, _ := ParseConfigName(fullName)
		namespaceSet[ns] = true
	}

	var namespaces []string
	for ns := range namespaceSet {
		namespaces = append(namespaces, ns)
	}
	return namespaces
}

// GetExportStatements returns shell export statements for a configuration
func (c *Config) GetExportStatements(name string) ([]string, error) {
	entry, err := c.GetEntry(name)
	if err != nil {
		return nil, err
	}

	var exports []string
	for k, v := range entry.Vars {
		exports = append(exports, fmt.Sprintf(text.Text.Formats.ExportStatement, k, v))
	}

	return exports, nil
}

// GetWebURL returns the web URL for a configuration
func (c *Config) GetWebURL(name string) (string, error) {
	entry, err := c.GetEntry(name)
	if err != nil {
		return "", err
	}

	if entry.WebURL == "" {
		return "", fmt.Errorf(text.Text.Errors.ConfigNoWebURL, name)
	}

	return entry.WebURL, nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}
