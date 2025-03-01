package system_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func TestCallbackTermination(t *testing.T) {
	app := app.New(1024)

	app.AddSystem(&system.CallbackTermination{
		Callback: func(t int64) bool {
			return t >= 99
		}},
	)

	app.Run()

	time := ecs.GetResource[resource.Tick](&app.World)
	assert.Equal(t, 100, int(time.Tick))
}

func ExampleCallbackTermination() {
	app := app.New(1024)

	app.AddSystem(&system.CallbackTermination{
		Callback: func(t int64) bool {
			return t >= 99
		}},
	)

	app.Run()
	// Output:
}
