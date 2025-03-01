package app

import (
	"math/rand/v2"
	"time"

	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

// App is the top-level ECS entrypoint.
//
// App provides access to the ECS world, and manages the scheduling
// of [System] and [UISystem] instances via [Systems].
// [System] instances are updated with a frequency given by TPS.
// [UISystem] instances are updated independently of normal systems,
// with a frequency given by FPS.
//
// The [Systems] scheduler, the app's [resource.Tick], [resource.Termination]
// and a central [resource.Rand] PRNG source can be accessed by systems as resources.
type App struct {
	Systems             // Systems manager and scheduler
	World     ecs.World // The ECS world
	rand      resource.Rand
	time      resource.Tick
	terminate resource.Termination
}

// New creates a new app.
//
// See [ecs.NewWorld] for the arguments.
func New(initialCapacity uint32) *App {
	var app = App{
		World: ecs.NewWorld(initialCapacity),
	}
	app.FPS = 30
	app.TPS = 0
	app.Systems.world = &app.World

	app.rand = resource.Rand{
		Source: rand.NewPCG(0, uint64(time.Now().UnixNano())),
	}
	ecs.AddResource(&app.World, &app.rand)
	app.time = resource.Tick{}
	ecs.AddResource(&app.World, &app.time)
	app.terminate = resource.Termination{}
	ecs.AddResource(&app.World, &app.terminate)

	ecs.AddResource(&app.World, &app.Systems)

	return &app
}

// Seed sets the random seed of the app's [resource.Rand].
// Call without an argument to seed from the current time.
//
// Systems should always use the Rand resource for PRNGs.
func (app *App) Seed(seed ...uint64) *App {
	switch len(seed) {
	case 0:
		app.rand.Source = rand.NewPCG(0, uint64(time.Now().UnixNano()))
	case 1:
		app.rand.Source = rand.NewPCG(0, seed[0])
	default:
		panic("can only use a single random seed")
	}
	return app
}

// Run the app, updating systems and ui systems according to App.TPS and App.FPS, respectively.
// Initializes the app if it is not already initialized.
// Finalizes the app after the run.
//
// Runs until Terminate in the resource resource.Termination is set to true
// (see [resource.Termination]).
//
// To perform updates manually, see [App.Update] and [App.UpdateUI],
// as well as [App.Initialize] and [App.Finalize].
func (app *App) Run() {
	app.Systems.run()
}

// Initialize the app.
func (app *App) Initialize() {
	app.Systems.initialize()
}

// Update the app's systems.
// Return whether the run should continue.
//
// Ignores App.TPS.
//
// Panics if [App.Initialize] was not called.
func (app *App) Update() bool {
	return app.Systems.updateSystems()
}

// UpdateUI the app's UI systems.
//
// Ignores App.FPS.
//
// Panics if [App.Initialize] was not called.
func (app *App) UpdateUI() {
	app.Systems.updateUISystems()
}

// Finalize the app.
func (app *App) Finalize() {
	app.Systems.finalize()
}

// Reset resets the world and removes all systems.
//
// Can be used to run systematic simulations without the need to re-allocate memory for each run.
// Accelerates re-populating the world by a factor of 2-3.
func (app *App) Reset() {
	app.World.Reset()
	app.Systems.reset()

	app.rand = resource.Rand{
		Source: rand.NewPCG(0, uint64(time.Now().UnixNano())),
	}
	ecs.AddResource(&app.World, &app.rand)
	app.time = resource.Tick{}
	ecs.AddResource(&app.World, &app.time)
	app.terminate = resource.Termination{}
	ecs.AddResource(&app.World, &app.terminate)

	ecs.AddResource(&app.World, &app.Systems)
}
