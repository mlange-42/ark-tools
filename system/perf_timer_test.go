package system_test

import (
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/system"
)

func TestPerfTimer(t *testing.T) {
	app := app.New(1024)

	app.AddSystem(&system.PerfTimer{UpdateInterval: 10})
	app.AddSystem(&system.FixedTermination{Steps: 30})

	app.Run()
}

func ExamplePerfTimer() {
	app := app.New(1024)

	app.AddSystem(&system.PerfTimer{UpdateInterval: 10})
	app.AddSystem(&system.FixedTermination{Steps: 30})

	// Uncomment the next line.

	// m.Run()
	// Output:
}
