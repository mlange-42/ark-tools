package app_test

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
)

func Example() {
	// Create a new, seeded model.
	app := app.New(1024).Seed(123)

	// Add systems.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	app.Run()
	// Output:
}
