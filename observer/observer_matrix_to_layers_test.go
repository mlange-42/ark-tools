package observer_test

import (
	"github.com/mlange-42/ark-tools/observer"
)

func ExampleMatrixToLayers() {
	// Multiple Matrix observers
	var matrix1 observer.Matrix = &MatrixObserver{}
	var matrix2 observer.Matrix = &MatrixObserver{}
	var matrix3 observer.Matrix = &MatrixObserver{}

	// A MatrixToGrid observer, wrapping the Matrix observers
	var _ observer.MatrixLayers = observer.MatrixToLayers(matrix1, matrix2, matrix3)
	// Output:
}
