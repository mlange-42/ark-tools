package main

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

// Position component.
type Position struct {
	X float64
	Y float64
}

// Velocity component.
type Velocity struct {
	X float64
	Y float64
}

func main() {
	// Create a new, seeded app.
	app := app.New(1024).Seed(123)
	// Limit simulation speed.
	app.TPS = 30

	// Add systems to the app.
	app.AddSystem(&VelocitySystem{EntityCount: 1000})
	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 100})

	// Run the app.
	app.Run()
}

// VelocitySystem is an example system adding velocity to position.
// For simplicity, it also creates entities during initialization.
type VelocitySystem struct {
	EntityCount int
	filter      ecs.Filter2[Position, Velocity]
}

// Initialize the system.
func (s *VelocitySystem) Initialize(w *ecs.World) {
	s.filter = *ecs.NewFilter2[Position, Velocity](w)

	mapper := ecs.NewMap2[Position, Velocity](w)
	mapper.NewBatch(s.EntityCount, &Position{}, &Velocity{})
}

// Update the system.
func (s *VelocitySystem) Update(w *ecs.World) {
	query := s.filter.Query()

	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

// Finalize the system.
func (s *VelocitySystem) Finalize(w *ecs.World) {}
