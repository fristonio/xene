package trigger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/executor"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"
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

	errs := errors.NewMultiError()
	// set running pipelines to true, as the controller is running the pipelines.
	t.RunningPipelines++

	log.Infof("starting to run pipelines: %v", t.Pipelines)
	for _, name := range t.Pipelines {
		val, err := store.KVStore.Get(context.TODO(), fmt.Sprintf("%s/%s", v1alpha1.PipelineKeyPrefix, name))
		if err != nil {
			errs.Append(fmt.Errorf("error while getting pipeline(%s) from kvstore: %s", name, err))
			continue
		}

		var pipeline v1alpha1.PipelineSpecWithName
		err = json.Unmarshal(val.Data, &pipeline)
		if err != nil {
			errs.Append(fmt.Errorf("error while unmarshaling pipeline(%s): %s", name, err))
		}

		pipelineID := utils.RandToken(defaults.PipelineIDSize)

		statusKey := fmt.Sprintf("%s/%s/%s", v1alpha1.PipelineStatusKeyPrefix, name, pipelineID)
		pipelineStatus := v1alpha1.GetDummyPipelineRunStatus(&pipeline)
		pipelineStatus.RunID = pipelineID
		pipelineStatus.Status = v1alpha1.StatusRunning
		pipelineStatus.Agent = option.Config.Agent.Name

		v, err := json.Marshal(&pipelineStatus)
		if err != nil {
			return fmt.Errorf("error while marshalling pipeline status: %s", err)
		}

		err = store.KVStore.Set(context.TODO(), statusKey, v)
		if err != nil {
			return fmt.Errorf("error while setting pipeline status key in kvstore: %s", err)
		}

		p := executor.NewPipelineExecutor(name, pipelineID, &pipeline)
		// Runs the pipeline
		go p.Run(pipelineStatus)
	}

	t.RunningPipelines--

	return nil
}
