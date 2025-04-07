package system

import (
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

// FixedTermination system.
//
// Terminates a run after a fixed number of ticks.
//
// Expects a resource of type [app.Termination].
type FixedTermination struct {
	Steps   int64 // Number of simulation ticks to run.
	tickRes ecs.Resource[resource.Tick]
	termRes ecs.Resource[resource.Termination]
}

// Initialize the system
func (s *FixedTermination) Initialize(w *ecs.World) {
	s.tickRes = ecs.NewResource[resource.Tick](w)
	s.termRes = ecs.NewResource[resource.Termination](w)
}

// Update the system
func (s *FixedTermination) Update(w *ecs.World) {
	tick := s.tickRes.Get().Tick

	if tick+1 >= s.Steps {
		term := s.termRes.Get()
		term.Terminate = true
	}
}

// Finalize the system
func (s *FixedTermination) Finalize(w *ecs.World) {}
