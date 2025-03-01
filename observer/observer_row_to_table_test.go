package observer_test

import (
	"github.com/mlange-42/ark-tools/observer"
)

func ExampleRowToTable() {
	// A Row observer
	var row observer.Row = &RowObserver{}

	// A RowToTable observer, wrapping the Row observer
	var _ observer.Table = observer.RowToTable(row)
	// Output:
}
