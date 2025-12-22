package observer_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/stretchr/testify/assert"
)

func TestGridToLayers(t *testing.T) {
	app := app.New(1024)

	var mat1 observer.Matrix = &matObs{}
	var mat2 observer.Matrix = &matObs{}
	var mat3 observer.Matrix = &matObs{}

	grid1 := observer.MatrixToGrid(mat1, nil, nil)
	grid2 := observer.MatrixToGrid(mat2, nil, nil)
	grid3 := observer.MatrixToGrid(mat3, nil, nil)

	layers := observer.GridToLayers(grid1, grid2, grid3)

	layers.Initialize(app.World)
	layers.Update(app.World)

	assert.Equal(t, 3, layers.Layers())

	w, h := layers.Dims()

	assert.Equal(t, 30, w)
	assert.Equal(t, 20, h)

	assert.Equal(t, 1.0, layers.X(1))
	assert.Equal(t, 1.0, layers.Y(1))

	data := layers.Values(app.World)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, 20*30, len(data[0]))
}

func TestGridToLayersFail(t *testing.T) {
	app := app.New(1024)

	var mat1 observer.Matrix = &matObs{}
	mat2 := &matObs{}
	mat2.Rows = 15

	grid1 := observer.MatrixToGrid(mat1, nil, nil)
	grid2 := observer.MatrixToGrid(mat2, nil, nil)

	layers := observer.GridToLayers(grid1, grid2)
	assert.Panics(t, func() { layers.Initialize(app.World) })

	assert.Panics(t, func() { observer.GridToLayers() })
}
