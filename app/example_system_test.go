package app_test

import (
	"fmt"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

// TestSystem is an example for implementing [System].
type TestSystem struct {
	timeRes ecs.Resource[resource.Tick]
}

// Initialize the system.
func (s *TestSystem) Initialize(w *ecs.World) {
	s.timeRes = ecs.NewResource[resource.Tick](w)
}

// Update the system.
func (s *TestSystem) Update(w *ecs.World) {
	time := s.timeRes.Get()
	fmt.Println(time.Tick)
}

// Finalize the system.
func (s *TestSystem) Finalize(w *ecs.World) {}

func ExampleSystem() {
	// Create a new model.
	app := app.New(1024)

	// Add the test system.
	app.AddSystem(&TestSystem{})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 30})

	// Run the simulation.
	app.Run()
}
