package app_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	app := app.New(1024)

	for i := 0; i < 3; i++ {
		app.Reset()
		app.Seed(123)

		app.AddSystem(&system.FixedTermination{
			Steps: 10,
		})

		app.Run()
	}
}

func TestAppStep(t *testing.T) {
	app := app.New(1024)

	for i := 0; i < 3; i++ {
		app.Reset()
		app.Seed(123)

		app.AddSystem(&system.FixedTermination{
			Steps: 10,
		})

		assert.Panics(t, func() { app.Update() })
		assert.Panics(t, func() { app.UpdateUI() })

		app.Initialize()

		app.Paused = true
		app.Update()
		app.Paused = false

		for app.Update() {
			app.UpdateUI()
		}
		app.Finalize()
	}
}

func TestAppSeed(t *testing.T) {
	app := app.New(1024)
	app.Seed(123)

	rand := ecs.GetResource[resource.Rand](&app.World)
	r1 := rand.Uint64()

	app.Seed(123)
	assert.Equal(t, r1, rand.Uint64())

	app.Seed()
	assert.NotEqual(t, r1, rand.Uint64())

	assert.Panics(t, func() { app.Seed(1, 2, 3) })
}

func ExampleApp() {
	// Create a new, seeded app.
	app := app.New(1024).Seed(123)

	// Add systems.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation.
	app.Run()
	// Output:
}

func ExampleApp_manualUpdate() {
	// Create a new, seeded app.
	app := app.New(1024).Seed(123)

	// Add systems.
	app.AddSystem(&system.FixedTermination{
		Steps: 100,
	})

	// Run the simulation manually.
	app.Initialize()
	for app.Update() {
		app.UpdateUI()
	}
	app.Finalize()
	// Output:
}

func ExampleApp_Reset() {
	// Create a new app.
	app := app.New(1024)

	// Do many simulations.
	for i := 0; i < 10; i++ {
		// Reset the app to clear entities, systems etc. before the run.
		app.Reset()

		// Seed the app for the run.
		app.Seed(uint64(i))

		// Add systems.
		app.AddSystem(&system.FixedTermination{
			Steps: 100,
		})

		// Run the simulation.
		app.Run()

	}
	// Output:
}
