package resource_test

import (
	"fmt"

	"math/rand/v2"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

func ExampleRand() {
	app := app.New(1024)

	src := ecs.GetResource[resource.Rand](&app.World)
	rng := rand.New(src.Source)
	_ = rng.NormFloat64()
	// Output:
}

func ExampleTick() {
	app := app.New(1024)

	tick := ecs.GetResource[resource.Tick](&app.World)

	fmt.Println(tick.Tick)
	// Output: 0
}

func ExampleTermination() {
	app := app.New(1024)

	term := ecs.GetResource[resource.Termination](&app.World)

	fmt.Println(term.Terminate)
	// Output: false
}

func ExampleSelectedEntity() {
	app := app.New(1024)
	ecs.AddResource(&app.World, &resource.SelectedEntity{})

	sel := ecs.GetResource[resource.SelectedEntity](&app.World)

	fmt.Println(sel.Selected.IsZero())
	// Output: true
}
