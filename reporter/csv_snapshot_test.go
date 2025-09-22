package reporter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotCSV(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	fileName := filepath.Join(tempDir, "test-%06d.csv")
	fmt.Println(fileName)

	app := app.New(1024)

	app.AddSystem(&reporter.SnapshotCSV{
		Observer:    &ExampleSnapshotObserver{},
		FilePattern: fileName,
	})
	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	_, err = os.Stat(filepath.Join(tempDir, "test-000000.csv"))
	assert.Nil(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "test-000099.csv"))
	assert.Nil(t, err)

	err = os.RemoveAll(tempDir)
	assert.Nil(t, err)
}

func TestSnapshotCSVFinal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	fileName := filepath.Join(tempDir, "test-%06d.csv")

	app := app.New(1024)

	app.AddSystem(&reporter.SnapshotCSV{
		Observer:    &ExampleSnapshotObserver{},
		FilePattern: fileName,
		Final:       true,
	})
	app.AddSystem(&system.FixedTermination{Steps: 100})

	app.Run()

	_, err = os.Stat(filepath.Join(tempDir, "test-000000.csv"))
	assert.NotNil(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "test-000100.csv"))
	assert.Nil(t, err)

	err = os.RemoveAll(tempDir)
	assert.Nil(t, err)
}
