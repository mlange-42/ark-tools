# Ark Tools

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/ark-tools/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/ark-tools/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/mlange-42/ark-tools/badge.svg?branch=main)](https://coveralls.io/github/mlange-42/ark-tools?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/ark-tools)](https://goreportcard.com/report/github.com/mlange-42/ark-tools)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/ark-tools.svg)](https://pkg.go.dev/github.com/mlange-42/ark-tools)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/ark-tools)
[![MIT license](https://img.shields.io/github/license/mlange-42/ark-tools)](https://github.com/mlange-42/ark-tools/blob/main/LICENSE)

Ark Tools provides a wrapper around the [Ark](https://github.com/mlange-42/ark) Entity Component System (ECS), and some common systems and resources.
It's purpose is to get started with prototyping and developing simulation models immediately, focussing on the model logic.

<div align="center">

<a href="https://github.com/mlange-42/ark">
<img src="https://github.com/user-attachments/assets/4bbe57c6-2e16-43be-ad5e-0cf26c220f21" alt="Ark (logo)" width="500px" />
</a>

</div>

## Features

* Scheduler for running logic and UI systems with independent update rates.
* Interfaces for ECS systems and observers.
* Ready-to-use systems for common tasks like writing CSV files or terminating a simulation.
* Common ECS resources, like central PRNG source or the current model tick.

## Installation

```
go get github.com/mlange-42/ark-tools
```

## Usage

See the [API docs](https://pkg.go.dev/github.com/mlange-42/ark-tools) for more details and examples.  
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/ark-tools.svg)](https://pkg.go.dev/github.com/mlange-42/ark-tools)

```go
package main

import (
	"github.com/mlange-42/ark-tools/model"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
)

// Position component
type Position struct {
	X float64
	Y float64
}

// Velocity component
type Velocity struct {
	X float64
	Y float64
}

func main() {
	// Create a new, seeded model.
	m := model.New().Seed(123)
	// Limit simulation speed
	m.TPS = 30

	// Add systems to the model.
	m.AddSystem(&VelocitySystem{EntityCount: 1000})
	// Add a termination system that ends the simulation.
	m.AddSystem(&system.FixedTermination{Steps: 100})

	// Run the model.
	m.Run()
}

// VelocitySystem is an example system adding velocity to position.
// For simplicity, it also creates entities during initialization.
type VelocitySystem struct {
	EntityCount int
	filter      ecs.Filter2[Position, Velocity]
}

// Initialize the system
func (s *VelocitySystem) Initialize(w *ecs.World) {
	s.filter = ecs.NewFilter2[Position, Velocity]()

	mapper := ecs.NewMap2[Position, Velocity](w)
	mapper.NewBatch(s.EntityCount)
}

// Update the system
func (s *VelocitySystem) Update(w *ecs.World) {
	query := s.filter.Query(w)

	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

// Finalize the system
func (s *VelocitySystem) Finalize(w *ecs.World) {}
```

## License

This project is distributed under the [MIT licence](./LICENSE).
