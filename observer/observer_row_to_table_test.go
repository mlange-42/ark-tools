package observer_test

import (
	"github.com/mlange-42/ark-tools/observer"
)

func ExampleRowToTable() {
	// A Row observer
	var row observer.Row = &RowObserver{}

	// A RowToTable observer, wrapping the Row observer
	_ = observer.RowToTable(row)
	// Output:
}
