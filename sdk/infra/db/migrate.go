package db

import (
	"fmt"
	"log"
)

var registeredModels []any

func RegisterModel(model any) {
	registeredModels = append(registeredModels, model)
}

func GetModels() []any {
	return registeredModels
}

type Seed struct {
	Name string
	Run  func() error
}

var seeds []Seed

func RegisterSeed(name string, fn func() error) {
	seeds = append(seeds, Seed{Name: name, Run: fn})
}

func RunSeeds() error {
	for _, s := range seeds {
		log.Printf("[Seed] Running: %s", s.Name)
		if err := s.Run(); err != nil {
			return fmt.Errorf("seed %q failed: %w", s.Name, err)
		}
	}
	return nil
}
