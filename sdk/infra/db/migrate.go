package db

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

type MigrationSnapshot struct {
	Models []string `json:"models"`
	Seeds  []string `json:"seeds"`
	Frozen bool     `json:"frozen"`
}

type Seed struct {
	Name string
	Run  func() error
}

type migrationRegistry struct {
	mu         sync.RWMutex
	models     []any
	modelNames []string
	modelIndex map[string]struct{}
	seeds      []Seed
	seedIndex  map[string]struct{}
	frozen     bool
}

var migrationGlobal = newMigrationRegistry()

func newMigrationRegistry() *migrationRegistry {
	return &migrationRegistry{
		modelIndex: make(map[string]struct{}),
		seedIndex:  make(map[string]struct{}),
	}
}

func RegisterModel(model any) {
	if err := migrationGlobal.registerModel(model); err != nil {
		panic(err)
	}
}

func GetModels() []any {
	return migrationGlobal.modelsSnapshot()
}

func RegisterSeed(name string, fn func() error) {
	if err := migrationGlobal.registerSeed(name, fn); err != nil {
		panic(err)
	}
}

func RunSeeds() error {
	return migrationGlobal.runSeeds()
}

func Freeze() {
	migrationGlobal.freeze()
}

func Snapshot() MigrationSnapshot {
	return migrationGlobal.snapshot()
}

func ResetForTest() {
	migrationGlobal = newMigrationRegistry()
}

func (r *migrationRegistry) registerModel(model any) error {
	if model == nil {
		return fmt.Errorf("model register failed: nil model")
	}
	name := modelTypeName(model)

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen {
		return fmt.Errorf("model register failed: registry frozen, model=%s", name)
	}
	if _, exists := r.modelIndex[name]; exists {
		return fmt.Errorf("model register failed: duplicate model=%s", name)
	}

	r.modelIndex[name] = struct{}{}
	r.models = append(r.models, model)
	r.modelNames = append(r.modelNames, name)
	return nil
}

func (r *migrationRegistry) modelsSnapshot() []any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]any(nil), r.models...)
}

func (r *migrationRegistry) registerSeed(name string, fn func() error) error {
	if name == "" {
		return fmt.Errorf("seed register failed: empty seed name")
	}
	if fn == nil {
		return fmt.Errorf("seed register failed: nil seed func, seed=%s", name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen {
		return fmt.Errorf("seed register failed: registry frozen, seed=%s", name)
	}
	if _, exists := r.seedIndex[name]; exists {
		return fmt.Errorf("seed register failed: duplicate seed=%s", name)
	}

	r.seedIndex[name] = struct{}{}
	r.seeds = append(r.seeds, Seed{Name: name, Run: fn})
	return nil
}

func (r *migrationRegistry) runSeeds() error {
	r.freeze()

	r.mu.RLock()
	seeds := append([]Seed(nil), r.seeds...)
	r.mu.RUnlock()

	for _, s := range seeds {
		log.Printf("[Seed] Running: %s", s.Name)
		if err := s.Run(); err != nil {
			return fmt.Errorf("seed %q failed: %w", s.Name, err)
		}
	}
	return nil
}

func (r *migrationRegistry) freeze() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frozen = true
}

func (r *migrationRegistry) snapshot() MigrationSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return MigrationSnapshot{
		Models: append([]string(nil), r.modelNames...),
		Seeds:  seedNames(r.seeds),
		Frozen: r.frozen,
	}
}

func modelTypeName(model any) string {
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.PkgPath() == "" {
		return t.String()
	}
	return t.PkgPath() + "." + t.Name()
}

func seedNames(items []Seed) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}
