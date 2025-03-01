package app

import (
	"testing"

	"github.com/mlange-42/ark-tools/resource"
	"github.com/mlange-42/ark-tools/system"
	"github.com/mlange-42/ark/ecs"
	"github.com/stretchr/testify/assert"
)

func TestSystems(t *testing.T) {
	app := New(1024)
	for i := 0; i < 3; i++ {
		app.Reset()

		app.Seed()
		app.Seed(123)

		termSys := system.FixedTermination{
			Steps: 1000,
		}
		uiSys := uiSystem{}
		dualSys := dualSystem{}

		app.AddSystem(&termSys)
		app.AddSystem(&system.FixedTermination{
			Steps: 10,
		})
		app.AddUISystem(&uiSys)

		app.locked = true
		assert.Panics(t, func() { app.removeSystem(&termSys) })
		assert.Panics(t, func() { app.removeUISystem(&uiSys) })
		app.locked = false

		assert.Panics(t, func() { app.AddSystem(&dualSys) })
		app.AddUISystem(&dualSys)

		app.AddSystem(&removerSystem{
			Remove:   []System{&termSys},
			RemoveUI: []UISystem{&uiSys},
		})

		assert.Equal(t, app.systems, app.Systems.Systems())
		assert.Equal(t, app.uiSystems, app.Systems.UISystems())

		assert.Panics(t, func() { app.RemoveSystem(&dualSys) })

		assert.Equal(t, 4, len(app.systems))
		assert.Equal(t, 2, len(app.uiSystems))
		assert.Equal(t, 0, len(app.toRemove))
		assert.Equal(t, 0, len(app.uiToRemove))

		app.Run()

		assert.Equal(t, 3, len(app.systems))
		assert.Equal(t, 1, len(app.uiSystems))
		assert.Equal(t, 0, len(app.toRemove))
		assert.Equal(t, 0, len(app.uiToRemove))

		assert.Panics(t, func() { app.initialize() })

		app.RemoveUISystem(&dualSys)

		assert.Equal(t, 2, len(app.systems))
		assert.Equal(t, 0, len(app.uiSystems))
		assert.Equal(t, 0, len(app.toRemove))
		assert.Equal(t, 0, len(app.uiToRemove))

		assert.Panics(t, func() { app.RemoveUISystem(&dualSys) })

		assert.Panics(t, func() { app.RemoveSystem(&termSys) })

		assert.Panics(t, func() { app.RemoveUISystem(&uiSys) })

		assert.Panics(t, func() { app.AddSystem(&termSys) })
		assert.Panics(t, func() { app.AddUISystem(&uiSys) })
	}
}

func TestSystemsInit(t *testing.T) {
	app := New(1024)
	app.TPS = 0
	app.FPS = 0

	app.AddSystem(&system.FixedTermination{Steps: 5})
	app.AddUISystem(&uiSystem{})
	app.Run()

	assert.Equal(t, 30.0, app.FPS)
	assert.Equal(t, 0.0, app.TPS)
	assert.Equal(t, 5, int(app.time.Tick))

	app = New(1024)
	app.TPS = 10

	app.AddSystem(&system.FixedTermination{Steps: 5})
	app.AddUISystem(&uiSystem{})

	app.Run()

	app = New(1024)
	app.TPS = 10
	app.FPS = 30

	app.AddSystem(&system.FixedTermination{Steps: 5})
	app.AddUISystem(&uiSystem{})

	app.Run()

	app = New(1024)
	app.TPS = 10
	app.FPS = -1

	app.AddSystem(&system.FixedTermination{Steps: 5})
	app.AddUISystem(&uiSystem{})

	app.Run()
}

func TestSystemsPaused(t *testing.T) {
	app := New(1024)
	app.TPS = 0
	app.FPS = 0

	app.AddSystem(&system.FixedTermination{Steps: 5})
	app.AddUISystem(&uiTerminationSystem{Steps: 100})

	app.Paused = true
	app.Run()

	assert.Equal(t, 0, int(app.time.Tick))
}

type uiSystem struct{}

func (s *uiSystem) InitializeUI(w *ecs.World) {}
func (s *uiSystem) UpdateUI(w *ecs.World)     {}
func (s *uiSystem) PostUpdateUI(w *ecs.World) {}
func (s *uiSystem) FinalizeUI(w *ecs.World)   {}

type uiTerminationSystem struct {
	Steps   int
	step    int
	termRes ecs.Resource[resource.Termination]
}

func (s *uiTerminationSystem) InitializeUI(w *ecs.World) {
	s.termRes = ecs.NewResource[resource.Termination](w)
	s.step = 0
}

func (s *uiTerminationSystem) UpdateUI(w *ecs.World) {
	if s.step >= s.Steps {
		term := s.termRes.Get()
		term.Terminate = true
	}
	s.step++
}
func (s *uiTerminationSystem) PostUpdateUI(w *ecs.World) {}
func (s *uiTerminationSystem) FinalizeUI(w *ecs.World)   {}

type dualSystem struct{}

func (s *dualSystem) Initialize(w *ecs.World)   {}
func (s *dualSystem) InitializeUI(w *ecs.World) {}
func (s *dualSystem) Update(w *ecs.World)       {}
func (s *dualSystem) UpdateUI(w *ecs.World)     {}
func (s *dualSystem) PostUpdateUI(w *ecs.World) {}
func (s *dualSystem) Finalize(w *ecs.World)     {}
func (s *dualSystem) FinalizeUI(w *ecs.World)   {}

type removerSystem struct {
	Remove   []System
	RemoveUI []UISystem
	step     int
}

func (s *removerSystem) Initialize(w *ecs.World) {}
func (s *removerSystem) Update(w *ecs.World) {
	if s.step == 3 {
		systems := ecs.GetResource[Systems](w)
		for _, sys := range s.Remove {
			systems.RemoveSystem(sys)
		}
		for _, sys := range s.RemoveUI {
			systems.RemoveUISystem(sys)
		}
	}
	s.step++
}
func (s *removerSystem) Finalize(w *ecs.World) {}
