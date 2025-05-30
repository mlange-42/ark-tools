package app

import (
	"fmt"
	"time"

	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark/ecs"
)

// System is the interface for ECS systems.
//
// See also [UISystem] for systems with an independent graphics step.
type System interface {
	Initialize(w *ecs.World) // Initialize the system.
	Update(w *ecs.World)     // Update the system.
	Finalize(w *ecs.World)   // Finalize the system.
}

// UISystem is the interface for ECS systems that display UI in an independent graphics step.
//
// See also [System] for normal systems.
type UISystem interface {
	InitializeUI(w *ecs.World) // InitializeUI the system.
	UpdateUI(w *ecs.World)     // UpdateUI/update the system.
	PostUpdateUI(w *ecs.World) // PostUpdateUI does the final part of updating, e.g. update the GL window.
	FinalizeUI(w *ecs.World)   // FinalizeUI the system.
}

// Systems manages and schedules ECS [System] and [UISystem] instances.
//
// [System] instances are updated with a frequency given by TPS (ticks per second).
// [UISystem] instances are updated independently of normal systems, with a frequency given by FPS (frames per second).
//
// [Systems] is an embed in [App] and it's methods are usually only used through a [App] instance.
// By also being a resource of each [App], however, systems can access it and e.g. remove themselves from aan app.
type Systems struct {
	// Ticks per second for normal systems.
	// Values <= 0 (the default) mean as fast as possible.
	TPS float64
	// Frames per second for UI systems.
	// A zero/unset value defaults to 30 FPS. Values < 0 sync FPS with TPS.
	// With fast movement, a value of 60 may be required for fluent graphics.
	FPS float64
	// Whether the simulation is currently paused.
	// When paused, only UI updates but no normal updates are performed.
	Paused bool

	world      *ecs.World
	systems    []System
	uiSystems  []UISystem
	toRemove   []System
	uiToRemove []UISystem

	nextDraw   time.Time
	nextUpdate time.Time

	initialized bool
	locked      bool

	tickRes ecs.Resource[resource.Tick]
	termRes ecs.Resource[resource.Termination]
}

// Systems returns the normal/non-UI systems.
func (s *Systems) Systems() []System {
	return s.systems
}

// UISystems returns the UI systems.
func (s *Systems) UISystems() []UISystem {
	return s.uiSystems
}

// AddSystem adds a [System] to the app.
//
// Panics if the system is also a [UISystem].
// To add systems that implement both [System] and [UISystem], use [Systems.AddUISystem]
func (s *Systems) AddSystem(sys System) {
	if s.initialized {
		panic("adding systems after app initialization is not implemented yet")
	}
	if sys, ok := sys.(UISystem); ok {
		panic(fmt.Sprintf("System %T is also an UI system. Must be added via AddUISystem.", sys))
	}
	s.systems = append(s.systems, sys)
}

// AddUISystem adds an [UISystem] to the app.
//
// Adds the [UISystem] also as a normal [System] if it implements the interface.
func (s *Systems) AddUISystem(sys UISystem) {
	if s.initialized {
		panic("adding systems after app initialization is not implemented yet")
	}
	s.uiSystems = append(s.uiSystems, sys)
	if sys, ok := sys.(System); ok {
		s.systems = append(s.systems, sys)
	}
}

// RemoveSystem removes a system from the app.
//
// Systems can also be removed during a run.
// However, this will take effect only after the end of the full update step.
func (s *Systems) RemoveSystem(sys System) {
	if sys, ok := sys.(UISystem); ok {
		panic(fmt.Sprintf("System %T is also an UI system. Must be removed via RemoveUISystem.", sys))
	}
	s.toRemove = append(s.toRemove, sys)
	if !s.locked {
		s.removeSystems()
	}
}

// RemoveUISystem removes an UI system from the app.
//
// Systems can also be removed during a run.
// However, this will take effect only after the end of the full update step.
func (s *Systems) RemoveUISystem(sys UISystem) {
	s.uiToRemove = append(s.uiToRemove, sys)
	if !s.locked {
		s.removeSystems()
	}
}

// Removes systems that were removed during the update step.
func (s *Systems) removeSystems() {
	rem := s.toRemove
	remUI := s.uiToRemove

	s.toRemove = s.toRemove[:0]
	s.uiToRemove = s.uiToRemove[:0]

	for _, sys := range rem {
		s.removeSystem(sys)
	}
	for _, sys := range remUI {
		if sys, ok := sys.(System); ok {
			s.removeSystem(sys)
		}
		s.removeUISystem(sys)
	}
}

func (s *Systems) removeSystem(sys System) {
	if s.locked {
		panic("can't remove a system in locked state")
	}
	idx := -1
	for i := 0; i < len(s.systems); i++ {
		if sys == s.systems[i] {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic(fmt.Sprintf("can't remove system %T: not in the app", sys))
	}
	s.systems[idx].Finalize(s.world)
	s.systems = append(s.systems[:idx], s.systems[idx+1:]...)
}

func (s *Systems) removeUISystem(sys UISystem) {
	if s.locked {
		panic("can't remove a system in locked state")
	}
	idx := -1
	for i := 0; i < len(s.uiSystems); i++ {
		if sys == s.uiSystems[i] {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic(fmt.Sprintf("can't remove UI system %T: not in the app", sys))
	}
	s.uiSystems[idx].FinalizeUI(s.world)
	s.uiSystems = append(s.uiSystems[:idx], s.uiSystems[idx+1:]...)
}

// Initialize all systems.
func (s *Systems) initialize() {
	if s.initialized {
		panic("app is already initialized")
	}

	if s.FPS == 0 {
		s.FPS = 30
	}

	s.tickRes = ecs.NewResource[resource.Tick](s.world)
	s.termRes = ecs.NewResource[resource.Termination](s.world)

	s.locked = true
	for _, sys := range s.systems {
		sys.Initialize(s.world)
	}
	for _, sys := range s.uiSystems {
		sys.InitializeUI(s.world)
	}
	s.locked = false
	s.removeSystems()
	s.initialized = true

	s.nextDraw = time.Time{}
	s.nextUpdate = time.Time{}

	s.tickRes.Get().Tick = 0
}

// Update all systems.
func (s *Systems) update() bool {
	s.locked = true
	update := s.updateSystemsTimed()
	s.updateUISystemsTimed(update)
	s.locked = false

	s.removeSystems()

	if update {
		time := s.tickRes.Get()
		time.Tick++
	} else {
		s.wait()
	}

	return !s.termRes.Get().Terminate
}

// updateSystems updates all normal systems
func (s *Systems) updateSystems() bool {
	if !s.initialized {
		panic("the app is not initialized")
	}
	if s.Paused {
		return true
	}
	s.locked = true
	updated := s.updateSystemsSimple()
	s.locked = false

	s.removeSystems()

	if updated {
		time := s.tickRes.Get()
		time.Tick++
	}

	return !s.termRes.Get().Terminate
}

// updateUISystems updates all UI systems
func (s *Systems) updateUISystems() {
	if !s.initialized {
		panic("the app is not initialized")
	}
	s.locked = true
	s.updateUISystemsSimple()
	s.locked = false

	s.removeSystems()
}

// Calculates and waits the time until the next update of UI update.
func (s *Systems) wait() {
	nextUpdate := s.nextUpdate

	if (s.Paused || s.FPS > 0) && s.nextDraw.Before(nextUpdate) {
		nextUpdate = s.nextDraw
	}

	t := time.Now()
	wait := nextUpdate.Sub(t)

	if wait > 0 {
		time.Sleep(wait)
	}
}

// Update normal systems.
func (s *Systems) updateSystemsSimple() bool {
	for _, sys := range s.systems {
		sys.Update(s.world)
	}
	return true
}

// Update normal systems.
func (s *Systems) updateSystemsTimed() bool {
	update := false
	if s.Paused {
		update = !time.Now().Before(s.nextUpdate)
		if update {
			tps := s.limitedFps(s.TPS, 10)
			s.nextUpdate = nextTime(s.nextUpdate, tps)
		}
		return false
	}
	if s.TPS <= 0 {
		update = true
		s.updateSystemsSimple()
	} else {
		update = !time.Now().Before(s.nextUpdate)
		if update {
			s.nextUpdate = nextTime(s.nextUpdate, s.TPS)
			s.updateSystemsSimple()
		}
	}
	return update
}

// Update ui systems.
func (s *Systems) updateUISystemsSimple() {
	for _, sys := range s.uiSystems {
		sys.UpdateUI(s.world)
	}
	for _, sys := range s.uiSystems {
		sys.PostUpdateUI(s.world)
	}
}

// Update UI systems.
func (s *Systems) updateUISystemsTimed(updated bool) {
	if !s.Paused && s.FPS <= 0 {
		if updated {
			s.updateUISystemsSimple()
		}
	} else {
		if !time.Now().Before(s.nextDraw) {
			fps := s.FPS
			if s.Paused {
				fps = s.limitedFps(s.FPS, 30)
			}
			s.nextDraw = nextTime(s.nextDraw, fps)
			s.updateUISystemsSimple()
		}
	}
}

// Finalize all systems.
func (s *Systems) finalize() {
	s.locked = true
	for _, sys := range s.systems {
		sys.Finalize(s.world)
	}
	for _, sys := range s.uiSystems {
		sys.FinalizeUI(s.world)
	}
	s.locked = false
	s.removeSystems()
}

// Run the app.
func (s *Systems) run() {
	if !s.initialized {
		s.initialize()
	}

	for s.update() {
	}

	s.finalize()
}

// Removes all systems.
func (s *Systems) reset() {
	s.systems = []System{}
	s.uiSystems = []UISystem{}
	s.toRemove = s.toRemove[:0]
	s.uiToRemove = s.uiToRemove[:0]

	s.nextDraw = time.Time{}
	s.nextUpdate = time.Time{}

	s.initialized = false
	s.tickRes = ecs.Resource[resource.Tick]{}
}

// Calculates frame rate capped to target
func (s *Systems) limitedFps(actual, target float64) float64 {
	if actual > target || actual <= 0 {
		return target
	}
	return actual
}
