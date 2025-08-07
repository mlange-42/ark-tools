package resource

import (
	"math/rand/v2"

	"github.com/mlange-42/ark/ecs"
)

// Rand is a PRNG resource to be used in [System] implementations.
//
// This resource is provided by [github.com/mlange-42/ark-tools/app.App] per default.
type Rand struct {
	rand.Source `json:"-"` // Source to use for PRNGs in [System] implementations.
}

// Tick is a resource holding the app's time step.
// The tick value should not be modified by user code, as it is managed by the scheduler.
//
// This resource is provided by [github.com/mlange-42/ark-tools/app.App] per default.
type Tick struct {
	Tick int64 // The current tick.
}

// Termination is a resource holding whether the simulation should terminate after the current step.
//
// This resource is provided by [github.com/mlange-42/ark-tools/app.App] per default.
type Termination struct {
	Terminate bool // Whether the simulation run is finished. Can be set by systems.
}

// SelectedEntity is a resource holding the currently selected entity.
//
// The primarily purpose is communication between UI systems, e.g. for entity inspection or manipulation by the user.
type SelectedEntity struct {
	Selected ecs.Entity
}
