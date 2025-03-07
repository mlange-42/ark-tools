package system

import (
	"fmt"
	"time"

	"github.com/mlange-42/ark/ecs"
)

// PerfTimer system for printing elapsed time per step, and optional world statistics.
type PerfTimer struct {
	UpdateInterval int // Update/print interval in ticks.
	start          time.Time
	startSim       time.Time
	step           int64
}

// Initialize the system
func (s *PerfTimer) Initialize(w *ecs.World) {
	s.step = 0
}

// Update the system
func (s *PerfTimer) Update(w *ecs.World) {
	t := time.Now()
	if s.step == 0 {
		s.start = t
		s.startSim = t
	}
	if s.step%int64(s.UpdateInterval) == 0 {
		if s.step > 0 {
			dur := t.Sub(s.start)
			usec := float64(dur.Microseconds()) / float64(s.UpdateInterval)
			fmt.Printf("%d updates, %0.2f us/update\n", s.UpdateInterval, usec)
		}
		s.start = t
	}
	s.step++
}

// Finalize the system
func (s *PerfTimer) Finalize(w *ecs.World) {
	t := time.Now()
	dur := t.Sub(s.startSim)
	usec := float64(dur.Microseconds()) / float64(s.step)
	fmt.Printf("Total: %d updates, %0.2f us/update\n", s.step, usec)
}
