package plugin

import (
	"fmt"
	"log"
	"sync"
)

// PluginInfo describes metadata about a plugin.
type PluginInfo struct {
	Name        string
	Version     string
	Description string
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
var state = newRegistryState()

type pluginStatus struct {
	Name        string `json:"name"`
	InitOK      bool   `json:"init_ok"`
	StartOK     bool   `json:"start_ok"`
	LastError   string `json:"last_error,omitempty"`
	Initialized bool   `json:"initialized"`
	Started     bool   `json:"started"`
}

type registryState struct {
	mu      sync.RWMutex
	plugins map[string]pluginStatus
}

func newRegistryState() *registryState {
	return &registryState{plugins: make(map[string]pluginStatus)}
}

// Register registers a plugin. Call this from init() to self-register.
func Register(m Plugin) {
	plugins = append(plugins, m)
	state.ensure(m.Name())
	log.Printf("[plugin] registered: %s", m.Name())
}

// InitAll runs Init() on all registered plugins in registration order.
func InitAll() error {
	for _, m := range plugins {
		log.Printf("[plugin] init: %s", m.Name())
		if err := m.Init(); err != nil {
			state.setInit(m.Name(), err)
			return fmt.Errorf("plugin %s init failed: %w", m.Name(), err)
		}
		state.setInit(m.Name(), nil)
	}
	return nil
}

// StartAll runs Start() on all registered plugins.
func StartAll() error {
	var startErr error
	for _, m := range plugins {
		if err := m.Start(); err != nil {
			state.setStart(m.Name(), err)
			log.Printf("[plugin] %s start failed: %v", m.Name(), err)
			if startErr == nil {
				startErr = fmt.Errorf("plugin %s start failed: %w", m.Name(), err)
			}
			continue
		}
		state.setStart(m.Name(), nil)
	}
	return startErr
}

// StopAll runs Stop() on all registered plugins in reverse order.
func StopAll() {
	for i := len(plugins) - 1; i >= 0; i-- {
		m := plugins[i]
		if err := m.Stop(); err != nil {
			log.Printf("[plugin] %s stop failed: %v", m.Name(), err)
		}
		state.setStopped(m.Name())
	}
}

func Snapshot() []pluginStatus {
	state.mu.RLock()
	defer state.mu.RUnlock()

	items := make([]pluginStatus, 0, len(plugins))
	for _, m := range plugins {
		items = append(items, state.plugins[m.Name()])
	}
	return items
}

func Ready() (bool, []pluginStatus) {
	snapshot := Snapshot()
	for _, item := range snapshot {
		if !item.InitOK || !item.StartOK {
			return false, snapshot
		}
	}
	return true, snapshot
}

func (s *registryState) ensure(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.plugins[name]; !ok {
		s.plugins[name] = pluginStatus{Name: name}
	}
}

func (s *registryState) setInit(name string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := s.plugins[name]
	item.Name = name
	item.Initialized = true
	item.InitOK = err == nil
	if err != nil {
		item.LastError = err.Error()
	}
	s.plugins[name] = item
}

func (s *registryState) setStart(name string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := s.plugins[name]
	item.Name = name
	item.Started = true
	item.StartOK = err == nil
	if err != nil {
		item.LastError = err.Error()
	} else if item.InitOK {
		item.LastError = ""
	}
	s.plugins[name] = item
}

func (s *registryState) setStopped(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := s.plugins[name]
	item.Started = false
	item.StartOK = false
	s.plugins[name] = item
}
