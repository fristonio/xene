package controller

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Map is the map of a controller name with the underlying Controller.
type Map map[string]*Controller

// Manager manages a ControllerMap and perform actions on it.
type Manager struct {
	controllers Map

	terminate chan struct{}
	mutex     sync.RWMutex
}

// NewManager Creates a new manager instance for the controller map.
func NewManager() *Manager {
	return &Manager{
		controllers: Map{},

		terminate: make(chan struct{}),
	}
}

// NoopFunc is a nil function.
func NoopFunc(_ctx context.Context) (FunctionResult, error) {
	return nil, nil
}

// GetAllControllers Returns the name of all the controllers that are managed by this
// manager.
func (m *Manager) GetAllControllers() []string {
	var ctrls []string

	for key := range m.controllers {
		ctrls = append(ctrls, key)
	}

	return ctrls
}

// UpdateController installs or updates a controller in the manager. A
// controller is identified by its name. If a controller with the name already
// exists, the controller will be shut down and replaced with the provided
// controller. Updating a controller will cause the DoFunc to be run
// immediately regardless of any previous conditions. It will also cause any
// statistics to be reset.
func (m *Manager) UpdateController(name, cType string, internal Internal) error {
	_, err := m.updateController(name, cType, internal)

	return err
}

func (m *Manager) updateController(name, cType string, internal Internal) (*Controller, error) {
	start := time.Now()

	if internal.StopFunc == nil {
		internal.StopFunc, _ = NewControllerFunction(NoopFunc) //nolint:errcheck
	}

	m.mutex.Lock()

	if m.controllers == nil {
		m.controllers = Map{}
	}

	ctrl, exists := m.controllers[name]
	if exists {
		m.mutex.Unlock()

		ctrl.getLogger().Debug("Updating existing controller")
		ctrl.mutex.Lock()
		ctrl.updateController(internal, true)
		ctrl.mutex.Unlock()

		ctrl.getLogger().Debug("Controller update time: ", time.Since(start))
	} else {
		ctrl = &Controller{
			name:  name,
			cType: cType,

			stop:       make(chan struct{}),
			update:     make(chan struct{}, 1),
			terminated: make(chan struct{}),

			executionStatistics: make(map[time.Time]time.Duration),
		}
		ctrl.updateController(internal, false)
		ctrl.getLogger().Debug("Starting new controller")

		ctrl.ctxDoFunc, ctrl.cancelDoFunc = context.WithCancel(context.Background())
		m.controllers[ctrl.name] = ctrl
		m.mutex.Unlock()

		go ctrl.RunController()
	}

	return ctrl, nil
}

func (m *Manager) removeController(ctrl *Controller) {
	ctrl.stopController()
	delete(m.controllers, ctrl.name)

	ctrl.getLogger().Debug("Removed controller")
}

func (m *Manager) lookup(name string) *Controller {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if c, ok := m.controllers[name]; ok {
		return c
	}

	return nil
}

func (m *Manager) removeAndReturnController(name string) (*Controller, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.controllers == nil {
		return nil, fmt.Errorf("empty controller map")
	}

	oldCtrl, ok := m.controllers[name]
	if !ok {
		return nil, fmt.Errorf("unable to find controller %s", name)
	}

	m.removeController(oldCtrl)

	return oldCtrl, nil
}

// RemoveController stops and removes a controller from the manager. If DoFunc
// is currently running, DoFunc is allowed to complete in the background.
func (m *Manager) RemoveController(name string) error {
	_, err := m.removeAndReturnController(name)
	return err
}

// RemoveControllerAndWait stops and removes a controller using
// RemoveController() and then waits for it to run to completion.
func (m *Manager) RemoveControllerAndWait(name string) error {
	oldCtrl, err := m.removeAndReturnController(name)
	if err == nil {
		<-oldCtrl.terminated
	}

	return err
}

func (m *Manager) removeAll() []*Controller {
	ctrls := []*Controller{}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.controllers == nil {
		return ctrls
	}

	for _, ctrl := range m.controllers {
		m.removeController(ctrl)
		ctrls = append(ctrls, ctrl)
	}

	return ctrls
}

// RemoveAll stops and removes all controllers of the manager
func (m *Manager) RemoveAll() {
	m.removeAll()
}

// RemoveAllAndWait stops and removes all controllers of the manager and then
// waits for all controllers to exit
func (m *Manager) RemoveAllAndWait() {
	ctrls := m.removeAll()
	for _, ctrl := range ctrls {
		<-ctrl.terminated
	}
}

// Terminate terminates all the controllers managed by the manager.
func (m *Manager) Terminate() {
	m.RemoveAllAndWait()

	<-m.terminate
	close(m.terminate)
}

// Wait waits for the manager to terminate the controllers.
func (m *Manager) Wait() {
	<-m.terminate
}

// GetStats Return the entire stats for manager.
func (m *Manager) GetStats() []*Status {
	var stats []*Status

	for _, controller := range m.controllers {
		stats = append(stats, controller.Status())
	}

	return stats
}

// PullLatestControllerStatistics pulls the latest controller stats.
func (m *Manager) PullLatestControllerStatistics() []ExecutionStat {
	stat := []ExecutionStat{}

	for _, controller := range m.controllers {
		statMap := controller.ExtractExecutionStatistics()

		for t, duration := range statMap {
			stat = append(stat, ExecutionStat{
				Name: controller.Name(),
				Type: controller.Type(),

				StartTime: t,
				Duration:  duration,
			})
			break
		}
	}

	return stat
}
