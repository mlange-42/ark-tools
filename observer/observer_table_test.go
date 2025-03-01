package observer_test

import (
	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark/ecs"
)

// Example observer, reporting a nonsense table.
type TableObserver struct{}

func (o *TableObserver) Initialize(w *ecs.World) {}

func (o *TableObserver) Update(w *ecs.World) {}

func (o *TableObserver) Header() []string {
	return []string{"TotalEntities"}
}

func (o *TableObserver) Values(w *ecs.World) [][]float64 {
	return [][]float64{
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
	}
}

func ExampleTable() {
	var _ observer.Table = &TableObserver{}
	// Output:
}
