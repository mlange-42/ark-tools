package reporter_test

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
)

func ExamplePrint() {
	// Create a new model.
	app := app.New(1024)

	// Add a Print reporter with an Observer.
	app.AddSystem(&reporter.Print{
		Observer: &ExampleObserver{},
	})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 3})

	// Run the simulation.
	app.Run()
	// Output:
	// [A B C]
	// [1 2 3]
	// [A B C]
	// [1 2 3]
	// [A B C]
	// [1 2 3]
}
