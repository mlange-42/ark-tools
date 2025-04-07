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
	tickRes  ecs.Resource[resource.Tick]
	termRes  ecs.Resource[resource.Termination]
}

// Initialize the system
func (s *CallbackTermination) Initialize(w *ecs.World) {
	s.tickRes = ecs.NewResource[resource.Tick](w)
	s.termRes = ecs.NewResource[resource.Termination](w)
}

// Update the system
func (s *CallbackTermination) Update(w *ecs.World) {
	tick := s.tickRes.Get().Tick

	if s.Callback(tick) {
		term := s.termRes.Get()
		term.Terminate = true
	}
}

// Finalize the system
func (s *CallbackTermination) Finalize(w *ecs.World) {}
