package observer_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/stretchr/testify/assert"
)

func TestLayersToGrid(t *testing.T) {
	app := app.New(1024)

	var mat1 observer.Matrix = &matObs{}
	var mat2 observer.Matrix = &matObs{}
	var mat3 observer.Matrix = &matObs{}

	var layers observer.MatrixLayers = observer.MatrixToLayers(mat1, mat2, mat3)

	var grid observer.GridLayers = observer.LayersToLayers(layers, &[2]float64{0, 0}, &[2]float64{1, 1})

	grid.Initialize(&app.World)
	grid.Update(&app.World)

	assert.Equal(t, 3, grid.Layers())

	w, h := grid.Dims()

	assert.Equal(t, 30, w)
	assert.Equal(t, 20, h)

	assert.Equal(t, 1.0, grid.X(1))
	assert.Equal(t, 1.0, grid.Y(1))

	data := grid.Values(&app.World)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, 20*30, len(data[0]))
}
