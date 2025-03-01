package observer_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/observer"
	"github.com/stretchr/testify/assert"
)

func TestMatrixToLayers(t *testing.T) {
	app := app.New(1024)

	var mat1 observer.Matrix = &matObs{}
	var mat2 observer.Matrix = &matObs{}
	var mat3 observer.Matrix = &matObs{}

	var layers observer.MatrixLayers = observer.MatrixToLayers(mat1, mat2, mat3)

	layers.Initialize(&app.World)
	layers.Update(&app.World)

	assert.Equal(t, 3, layers.Layers())

	w, h := layers.Dims()

	assert.Equal(t, 30, w)
	assert.Equal(t, 20, h)

	data := layers.Values(&app.World)
	assert.Equal(t, 3, len(data))
	assert.Equal(t, 20*30, len(data[0]))
}

func TestMatrixToLayersFail(t *testing.T) {
	app := app.New(1024)

	var mat1 observer.Matrix = &matObs{}
	var mat2 *matObs = &matObs{}
	mat2.Rows = 15

	var layers observer.MatrixLayers = observer.MatrixToLayers(mat1, mat2)
	assert.Panics(t, func() { layers.Initialize(&app.World) })

	assert.Panics(t, func() { observer.MatrixToLayers() })
}
