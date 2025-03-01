package reporter_test

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleCSV() {
	// Create a new model.
	app := app.New(1024)

	// Add a CSV reporter with an Observer.
	app.AddSystem(&reporter.CSV{
		Observer: &ExampleObserver{},
		File:     "../out/test.csv",
		Sep:      ";",
	})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 100})

	// Run the simulation.
	app.Run()
}

// ExampleObserver to generate some simple time series.
type ExampleObserver struct{}

func (o *ExampleObserver) Initialize(w *ecs.World) {}
func (o *ExampleObserver) Update(w *ecs.World)     {}
func (o *ExampleObserver) Header() []string {
	return []string{"A", "B", "C"}
}
func (o *ExampleObserver) Values(w *ecs.World) []float64 {
	return []float64{1, 2, 3}
}
