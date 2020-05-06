package controller

import (
	"fmt"

	"github.com/fristonio/xene/pkg/agent/controller/trigger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry = logrus.WithFields(logrus.Fields{
		"package": "trigger-controller",
	})
)

var (
	// RegisteredControllers contains a list of all the registered controllers
	// for the agent to run
	RegisteredControllers []Controller
)

// Controller is the standard interface which should be implemented by all the registered controllers.
type Controller interface {
	// Configure sets up the controller
	Configure()

	// Run starts running the controller.
	Run() error

	// Stop shuts down the controller running.
	Stop() error

	// Name returns the name of the controller
	Name() string
}

// RunControllers run all the controller registered by the agent.
func RunControllers() error {
	log.Info("starting to run configured controllers for the agent")
	errs := make(map[string]string)
	for _, controller := range RegisteredControllers {
		controller.Configure()
		err := controller.Run()
		if err != nil {
			errs[controller.Name()] = err.Error()
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("error while running agent controllers: %v", errs)
	}

	return nil
}

// StopControllers stops running all the controllers managed by API server.
func StopControllers() error {
	log.Info("stopping running controllers for the agent.")
	errs := make(map[string]string)
	for _, controller := range RegisteredControllers {
		err := controller.Stop()
		if err != nil {
			errs[controller.Name()] = err.Error()
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("error while running the agent controllers: %v", errs)
	}

	return nil
}

func init() {
	RegisteredControllers = append(RegisteredControllers, trigger.TriggerCtrl)
}
