package reporter_test

import (
	"os"
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func TestCSV(t *testing.T) {
	app := app.New(1024)

	app.AddSystem(&reporter.CSV{
		Observer: &ExampleObserver{},
		File:     "../out/test.csv",
	})
	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	_, err := os.Stat("../out/test.csv")
	assert.Nil(t, err)
}
