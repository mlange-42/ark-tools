package system_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func TestFixedTermination(t *testing.T) {
	app := app.New(1024)

	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	time := ecs.GetResource[resource.Tick](&app.World)
	assert.Equal(t, 100, int(time.Tick))
}

func ExampleFixedTermination() {
	app := app.New(1024)

	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()
	// Output:
}
