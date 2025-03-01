package reporter_test

import (
	"os"
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotCSV(t *testing.T) {
	app := app.New(1024)

	app.AddSystem(&reporter.SnapshotCSV{
		Observer:    &ExampleSnapshotObserver{},
		FilePattern: "../out/test-%06d.csv",
	})
	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	_, err := os.Stat("../out/test-000000.csv")
	assert.Nil(t, err)
	_, err = os.Stat("../out/test-000090.csv")
	assert.Nil(t, err)
}
