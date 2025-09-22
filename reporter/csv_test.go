package reporter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func TestCSV(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	fileName := filepath.Join(tempDir, "test.csv")

	app := app.New(1024)

	app.AddSystem(&reporter.CSV{
		Observer: &ExampleObserver{},
		File:     fileName,
	})
	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	_, err = os.Stat(fileName)
	assert.Nil(t, err)

	err = os.RemoveAll(tempDir)
	assert.Nil(t, err)
}

func TestCSVFinal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	fileName := filepath.Join(tempDir, "test.csv")

	app := app.New(1024)

	app.AddSystem(&reporter.CSV{
		Observer: &ExampleObserver{},
		File:     fileName,
		Final:    true,
	})
	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	_, err = os.Stat(fileName)
	assert.Nil(t, err)

	err = os.RemoveAll(tempDir)
	assert.Nil(t, err)
}
