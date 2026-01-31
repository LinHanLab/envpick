package core

import (
	"fmt"

	"envpick/internal/config"
	"envpick/internal/selector"
	"envpick/internal/text"
)

// Engine handles the core logic of envpick
type Engine struct {
	config    *config.Config
	state     *config.State
	namespace string // current namespace for this engine instance
}

// NewEngine creates a new Engine with default namespace, loading config and state
func NewEngine() (*Engine, error) {
	return NewEngineWithNamespace("")
}

// NewEngineWithNamespace creates a new Engine with a specific namespace
func NewEngineWithNamespace(namespace string) (*Engine, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	state, err := config.LoadState()
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		config:    cfg,
		state:     state,
		namespace: namespace,
	}

	return engine, nil
}

func (e *Engine) GetConfig() *config.Config {
	return e.config
}

// GetNamespace returns the current namespace
func (e *Engine) GetNamespace() string {
	return e.namespace
}

// saveState saves the current state
func (e *Engine) saveState() error {
	return e.state.Save()
}

// GetCurrentConfig returns the current active configuration name (short form)
func (e *Engine) GetCurrentConfig() string {
	return e.state.GetCurrentConfig(e.namespace)
}

// GetCurrentConfigFull returns the full configuration name (with namespace prefix if applicable)
func (e *Engine) GetCurrentConfigFull() string {
	shortName := e.state.GetCurrentConfig(e.namespace)
	if shortName == "" {
		return ""
	}
	return config.BuildConfigName(e.namespace, shortName)
}

// SetCurrentConfig sets the current configuration (accepts short form name)
func (e *Engine) SetCurrentConfig(name string) error {
	// Build full config name for validation
	fullName := config.BuildConfigName(e.namespace, name)
	if _, ok := e.config.Configs[fullName]; !ok {
		return fmt.Errorf(text.Text.Errors.ConfigNotFound, name)
	}
	e.state.SetCurrentConfig(e.namespace, name)
	return e.saveState()
}

// GetOptions returns options for fzf selection (filtered by namespace)
func (e *Engine) GetOptions() []selector.Option {
	var options []selector.Option
	current := e.state.GetCurrentConfig(e.namespace)

	// Get configs for this namespace
	namespaceConfigs := e.config.GetNamespaceConfigs(e.namespace)

	for name := range namespaceConfigs {
		opt := selector.Option{
			Name: name,
		}

		// Set status
		if name == current {
			opt.Status = "active"
		}

		options = append(options, opt)
	}

	return options
}
