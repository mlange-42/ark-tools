package app_test

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleSystems() {
	// Create a new model.
	myApp := app.New(1024)

	// The app contains Systems as an embed, TPS and FPS are accessible through the model directly.
	myApp.TPS = 1000
	myApp.FPS = 60

	// Create a system
	sys := system.FixedTermination{
		Steps: 10,
	}

	// Add the system the usual way, through the model.
	// The model contains Systems as an embed, so actually [Systems.AddSystem] is called.
	myApp.AddSystem(&sys)

	// Inside systems, [Systems] can be accessed as a resource.
	systems := ecs.GetResource[app.Systems](&myApp.World)

	// Pause the simulation, e.g. based on user input.
	systems.Paused = true

	// Remove the system using the resource.
	systems.RemoveSystem(&sys)
	// Output:
}
