package app_test

import (
	"fmt"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

// TestUISystem is an example for implementing [UISystem].
type TestUISystem struct {
	timeRes ecs.Resource[resource.Tick]
}

// Initialize the system.
func (s *TestUISystem) InitializeUI(w *ecs.World) {
	s.timeRes = ecs.NewResource[resource.Tick](w)
}

// Update the system.
func (s *TestUISystem) UpdateUI(w *ecs.World) {
	time := s.timeRes.Get()
	fmt.Println(time.Tick)
}

// PostUpdate the system.
func (s *TestUISystem) PostUpdateUI(w *ecs.World) {}

// Finalize the system.
func (s *TestUISystem) FinalizeUI(w *ecs.World) {}

func ExampleUISystem() {
	// Create a new model.
	app := app.New(1024)

	// Add the test ui system.
	app.AddUISystem(&TestUISystem{})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 30})

	// Run the simulation.
	app.Run()
}
