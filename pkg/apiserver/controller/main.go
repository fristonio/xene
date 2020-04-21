package controller

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	// RegisteredControllers contains a list of all the registered controllers
	// for the API server to run.
	RegisteredControllers []Controller
)

// Controller is the standard interface which should be implemented by all the registered controllers.
type Controller interface {
	// Run starts running the controller.
	Run() error

	// Stop shuts down the controller running.
	Stop() error

	// Name returns the name of the controller
	Name() string
}

// RunControllers run all the controller registered by API server.
func RunControllers() error {
	log.Info("starting to run configured controllers for API server")
	errs := make(map[string]string)
	for _, controller := range RegisteredControllers {
		err := controller.Run()
		if err != nil {
			errs[controller.Name()] = err.Error()
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("error while running apiserver controllers: %v", errs)
	}

	return nil
}
