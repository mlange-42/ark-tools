package reporter_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/ark-tools/app"
	"github.com/mlange-42/ark-tools/reporter"
	"github.com/mlange-42/ark-tools/system"
	"github.com/stretchr/testify/assert"
)

func ExampleRowCallback() {
	// Create a new model.
	app := app.New(1024)

	data := [][]float64{}

	// Add a Print reporter with an Observer.
	app.AddSystem(&reporter.RowCallback{
		Observer: &ExampleObserver{},
		Callback: func(step int, row []float64) {
			data = append(data, row)
		},
		HeaderCallback: func(header []string) {},
	})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 3})

	// Run the simulation.
	app.Run()

	fmt.Println(data)
	// Output:
	// [[1 2 3] [1 2 3] [1 2 3]]
}

func ExampleTableCallback() {
	// Create a new model.
	app := app.New(1024)

	data := [][]float64{}

	// Add a Print reporter with an Observer.
	app.AddSystem(&reporter.TableCallback{
		Observer: &ExampleSnapshotObserver{},
		Callback: func(step int, table [][]float64) {
			data = append(data, table...)
		},
		HeaderCallback: func(header []string) {},
	})

	// Add a termination system that ends the simulation.
	app.AddSystem(&system.FixedTermination{Steps: 3})

	// Run the simulation.
	app.Run()

	fmt.Println(data)
	// Output:
	// [[1 2 3] [1 2 3] [1 2 3] [1 2 3] [1 2 3] [1 2 3] [1 2 3] [1 2 3] [1 2 3]]
}

func TestRowCallbackFinal(t *testing.T) {
	app := app.New(1024)
	counter := 0

	app.AddSystem(&reporter.RowCallback{
		Observer: &ExampleObserver{},
		Callback: func(step int, row []float64) {
			counter++
		},
		HeaderCallback: func(header []string) {},
		Final:          true,
	})
	app.AddSystem(&system.FixedTermination{Steps: 3})
	app.Run()

	assert.Equal(t, 1, counter)
}

func TestTableCallbackFinal(t *testing.T) {
	app := app.New(1024)
	counter := 0

	app.AddSystem(&reporter.TableCallback{
		Observer: &ExampleSnapshotObserver{},
		Callback: func(step int, table [][]float64) {
			counter++
		},
		HeaderCallback: func(header []string) {},
		Final:          true,
	})
	app.AddSystem(&system.FixedTermination{Steps: 3})
	app.Run()

	assert.Equal(t, 1, counter)
}
