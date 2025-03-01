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
	termRes ecs.Resource[resource.Termination]
	step    int64
}

// Initialize the system
func (s *FixedTermination) Initialize(w *ecs.World) {
	s.termRes = ecs.NewResource[resource.Termination](w)
	s.step = 0
}

// Update the system
func (s *FixedTermination) Update(w *ecs.World) {
	term := s.termRes.Get()

	if s.step+1 >= s.Steps {
		term.Terminate = true
	}
	s.step++
}

// Finalize the system
func (s *FixedTermination) Finalize(w *ecs.World) {}
