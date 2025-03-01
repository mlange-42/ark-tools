package reporter_test

import (
	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

func ExampleSnapshotCSV() {
	// Create a new model.
	app := app.New(1024)

	// Add a SnapshotCSV reporter with an Observer.
	app.AddSystem(&reporter.SnapshotCSV{
		Observer:    &ExampleSnapshotObserver{},
		FilePattern: "../out/test-%06d.csv",
		Sep:         ";",
	})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 100})

	// Run the simulation.
	app.Run()
}

// ExampleSnapshotObserver to generate some simple tables.
type ExampleSnapshotObserver struct{}

func (o *ExampleSnapshotObserver) Initialize(w *ecs.World) {}
func (o *ExampleSnapshotObserver) Update(w *ecs.World)     {}
func (o *ExampleSnapshotObserver) Header() []string {
	return []string{"A", "B", "C"}
}
func (o *ExampleSnapshotObserver) Values(w *ecs.World) [][]float64 {
	return [][]float64{
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
	}
}
