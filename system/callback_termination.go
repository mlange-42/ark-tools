package system

import (
	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

// CallbackTermination system.
//
// Terminates a run according to the return value of a callback function.
//
// Expects a resource of type [app.Termination].
type CallbackTermination struct {
	Callback func(t int64) bool // The callback. ends the simulation when it returns true.
	termRes  ecs.Resource[resource.Termination]
	step     int64
}

// Initialize the system
func (s *CallbackTermination) Initialize(w *ecs.World) {
	s.termRes = ecs.NewResource[resource.Termination](w)
	s.step = 0
}

// Update the system
func (s *CallbackTermination) Update(w *ecs.World) {
	term := s.termRes.Get()

	if s.Callback(s.step) {
		term.Terminate = true
	}
	s.step++
}

// Finalize the system
func (s *CallbackTermination) Finalize(w *ecs.World) {}
