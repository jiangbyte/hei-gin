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
	Info() PluginInfo
	Name() string
	Init() error
	Start() error
	Stop() error
}

// NoopPlugin can be embedded to avoid implementing all methods.
type NoopPlugin struct{}

func (NoopPlugin) Info() PluginInfo { return PluginInfo{} }
func (NoopPlugin) Init() error      { return nil }
func (NoopPlugin) Start() error     { return nil }
func (NoopPlugin) Stop() error      { return nil }

type SnapshotItem struct {
	Name        string `json:"name"`
	InitOK      bool   `json:"init_ok"`
	StartOK     bool   `json:"start_ok"`
	LastError   string `json:"last_error,omitempty"`
	Initialized bool   `json:"initialized"`
	Started     bool   `json:"started"`
}

type registry struct {
	mu       sync.RWMutex
	items    []Plugin
	index    map[string]int
	status   map[string]SnapshotItem
	frozen   bool
	initDone bool
}

var global = newRegistry()

func newRegistry() *registry {
	return &registry{
		index:  make(map[string]int),
		status: make(map[string]SnapshotItem),
	}
}

// Register registers a plugin. Call this from init() to self-register.
func Register(m Plugin) {
	if err := global.register(m); err != nil {
		panic(err)
	}
}

func MustRegister(m Plugin) {
	Register(m)
}

// InitAll runs Init() on all registered plugins in registration order.
func InitAll() error {
	return global.initAll()
}

// StartAll runs Start() on all registered plugins.
func StartAll() error {
	return global.startAll()
}

// StopAll runs Stop() on all registered plugins in reverse order.
func StopAll() {
	global.stopAll()
}

func Freeze() {
	global.freeze()
}

func Snapshot() []SnapshotItem {
	return global.snapshot()
}

func Ready() (bool, []SnapshotItem) {
	snapshot := Snapshot()
	for _, item := range snapshot {
		if !item.InitOK || !item.StartOK {
			return false, snapshot
		}
	}
	return true, snapshot
}

func ResetForTest() {
	global = newRegistry()
}

func (r *registry) register(m Plugin) error {
	if m == nil {
		return fmt.Errorf("plugin register failed: nil plugin")
	}

	name := m.Name()
	if name == "" {
		return fmt.Errorf("plugin register failed: empty plugin name")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen {
		return fmt.Errorf("plugin register failed: registry frozen, plugin=%s", name)
	}
	if _, ok := r.index[name]; ok {
		return fmt.Errorf("plugin register failed: duplicate plugin=%s", name)
	}

	r.index[name] = len(r.items)
	r.items = append(r.items, m)
	r.status[name] = SnapshotItem{Name: name}
	log.Printf("[plugin] registered: %s", name)
	return nil
}

func (r *registry) initAll() error {
	r.freeze()

	r.mu.RLock()
	items := append([]Plugin(nil), r.items...)
	r.mu.RUnlock()

	for _, m := range items {
		log.Printf("[plugin] init: %s", m.Name())
		if err := m.Init(); err != nil {
			r.setInitStatus(m.Name(), err)
			return fmt.Errorf("plugin %s init failed: %w", m.Name(), err)
		}
		r.setInitStatus(m.Name(), nil)
	}

	r.mu.Lock()
	r.initDone = true
	r.mu.Unlock()
	return nil
}

func (r *registry) startAll() error {
	r.mu.RLock()
	items := append([]Plugin(nil), r.items...)
	initDone := r.initDone
	r.mu.RUnlock()

	if !initDone {
		return fmt.Errorf("plugin start failed: InitAll must run first")
	}

	var startErr error
	for _, m := range items {
		if err := m.Start(); err != nil {
			r.setStartStatus(m.Name(), err)
			log.Printf("[plugin] %s start failed: %v", m.Name(), err)
			if startErr == nil {
				startErr = fmt.Errorf("plugin %s start failed: %w", m.Name(), err)
			}
			continue
		}
		r.setStartStatus(m.Name(), nil)
	}
	return startErr
}

func (r *registry) stopAll() {
	r.mu.RLock()
	items := append([]Plugin(nil), r.items...)
	r.mu.RUnlock()

	for i := len(items) - 1; i >= 0; i-- {
		m := items[i]
		if err := m.Stop(); err != nil {
			log.Printf("[plugin] %s stop failed: %v", m.Name(), err)
		}
		r.setStopped(m.Name())
	}
}

func (r *registry) freeze() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frozen = true
}

func (r *registry) snapshot() []SnapshotItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]SnapshotItem, 0, len(r.items))
	for _, p := range r.items {
		items = append(items, r.status[p.Name()])
	}
	return items
}

func (r *registry) setInitStatus(name string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item := r.status[name]
	item.Name = name
	item.Initialized = true
	item.InitOK = err == nil
	if err != nil {
		item.LastError = err.Error()
	}
	r.status[name] = item
}

func (r *registry) setStartStatus(name string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item := r.status[name]
	item.Name = name
	item.Started = true
	item.StartOK = err == nil
	if err != nil {
		item.LastError = err.Error()
	} else if item.InitOK {
		item.LastError = ""
	}
	r.status[name] = item
}

func (r *registry) setStopped(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item := r.status[name]
	item.Started = false
	item.StartOK = false
	r.status[name] = item
}
