package trigger

import (
	"context"
	"fmt"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// Trigger contains details corresponding to the provided
// trigger.
type Trigger struct {
	// RunningPipelines is set to true if the Trigger is currently running
	// any pipeline.
	RunningPipelines int

	// UpdateAvailable indidcates wheather an update is available for the trigger
	// configured, if a signal is recieved it will first try and let all the running pipeline
	// execute and then update the manifest scheduling further pipelines run
	// later on.
	UpdateAvailable bool

	// TriggerSpecWithName contains the specification manifest of the trigger.
	*v1alpha1.TriggerSpecWithName
}

// SetupController sets up the controller for the trigger.
func (t *Trigger) SetupController(manager *controller.Manager) error {
	log.Infof("setting up controller for trigger: %s", t.Name)

	switch v1alpha1.TriggerType(t.Type) {
	case v1alpha1.CronTriggerType:
		triggerName := v1alpha1.GetWorkflowPrefixedName(t.Workflow, t.Name)
		ctrlFunc, err := controller.NewControllerFunction(t.RunPipelines, manager)
		if err != nil {
			return fmt.Errorf("error while creating controller function: %s", err)
		}

		i, err := controller.NewControllerInternalWithCron(t.CronConfig, ctrlFunc)
		if err != nil {
			return fmt.Errorf("error while parsing cron spec: %s", err)
		}

		err = manager.UpdateController(triggerName, "trigger-controller", i)
		if err != nil {
			return fmt.Errorf("error while updating trigger controller for : %s", t.Name)
		}
	default:
		return fmt.Errorf("not a valid trigger type: %s", t.Type)
	}

	return nil
}

// StopController stops the controller for the trigger.
func (t *Trigger) StopController(manager *controller.Manager) error {
	log.Infof("stopping controller for trigger: %s", t.Name)
	triggerName := v1alpha1.GetWorkflowPrefixedName(t.Workflow, t.Name)

	err := manager.RemoveControllerAndWait(triggerName)
	if err != nil {
		return fmt.Errorf("error while removing controller: %s", err)
	}

	return nil
}

// RunPipelines runs the pipelines associated with the trigger.
func (t *Trigger) RunPipelines(ctx context.Context, manager *controller.Manager) error {
	if t.UpdateAvailable {
		return fmt.Errorf("looking to update the pipeline specification, the Pipeline will be run next time")
	}

	// set running pipelines to true, as the controller is running the pipelines.
	t.RunningPipelines++

	log.Infof("starting to run pipelines: %v", t.Pipelines)

	t.RunningPipelines--

	return nil
}
