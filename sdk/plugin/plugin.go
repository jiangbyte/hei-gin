package plugin

import "log"

// PluginInfo describes metadata about a plugin.
type PluginInfo struct {
	Name         string
	Version      string
	Description  string
	Dependencies []string
}

// Plugin defines the lifecycle of an app plugin.
type Plugin interface {
	// Info returns metadata about the plugin.
	Info() PluginInfo
	// Name returns the plugin name for logging.
	Name() string
	// Init is called during app startup, after config and DB are ready, before the HTTP server starts.
	Init() error
	// Start is called after the HTTP server starts (for background tasks like cron).
	Start() error
	// Stop is called during graceful shutdown.
	Stop() error
}

// NoopPlugin can be embedded to avoid implementing all methods.
type NoopPlugin struct{}

func (NoopPlugin) Info() PluginInfo { return PluginInfo{} }
func (NoopPlugin) Init() error      { return nil }
func (NoopPlugin) Start() error     { return nil }
func (NoopPlugin) Stop() error      { return nil }

var plugins []Plugin

// Register registers a plugin. Call this from init() to self-register.
func Register(m Plugin) {
	plugins = append(plugins, m)
	log.Printf("[plugin] registered: %s", m.Name())
}

// InitAll runs Init() on all registered plugins in registration order.
func InitAll() error {
	for _, m := range plugins {
		log.Printf("[plugin] init: %s", m.Name())
		if err := m.Init(); err != nil {
			return err
		}
	}
	return nil
}

// StartAll runs Start() on all registered plugins.
func StartAll() {
	for _, m := range plugins {
		if err := m.Start(); err != nil {
			log.Printf("[plugin] %s start error: %v", m.Name(), err)
		}
	}
}

// StopAll runs Stop() on all registered plugins in reverse order.
func StopAll() {
	for i := len(plugins) - 1; i >= 0; i-- {
		m := plugins[i]
		if err := m.Stop(); err != nil {
			log.Printf("[plugin] %s stop error: %v", m.Name(), err)
		}
	}
}
