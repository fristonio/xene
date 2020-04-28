package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// Controller represents a type corresponding to workflow manager and controller.
type Controller struct {
	// storeCtrl is the store controller corresponding to the workflow objects
	// in the datastore.
	storeCtrl *store.Controller

	// name contains the name of the controller
	name string

	Scheduler *Scheduler
}

// WorkflowCtrl is the global instance of the workflow controller running for the API server.
var WorkflowCtrl *Controller = NewController()

// NewController returns a new agent controller manager which manages all the
// agents running for the API server.
func NewController() *Controller {
	return &Controller{
		Scheduler: NewSchedulerWithDefaultAgentCtrl(),
	}
}

// Type returns the type of controller configured, for agent this is agent.
func (a *Controller) Type() string {
	return "workflow"
}

// Configure sets up the Agent controller and all its required components.
func (a *Controller) Configure() {
	a.storeCtrl = a.newWorkflowStoreController()
	a.name = a.storeCtrl.Name()
}

// Run starts running the workflow controller.
func (a *Controller) Run() error {
	return a.storeCtrl.Run()
}

// Stop shuts down the controller.
func (a *Controller) Stop() error {
	return a.storeCtrl.Stop()
}

// Name returns the name of the agent controller, it is completely defined by the name of
// the underlying store controller.
func (a *Controller) Name() string {
	return a.name
}

// newWorkflowStoreController returns the workflow store controller for apiserver.
// This controller watches for workflow object in the store and perform action based
// on the changes to the object.
func (a *Controller) newWorkflowStoreController() *store.Controller {
	return store.NewControllerWithSharedCache(
		fmt.Sprintf("%s/", v1alpha1.WorkflowKeyPrefix),
		// Add function for a new workflow.
		func(kv *v1alpha1.KVPairStruct) error {
			log.Infof("workflow added: %s", kv.Key)
			return a.addWorkflow(kv)
		},

		// Delete function for the workflow object.
		func(key string) error {
			log.Infof("workflow deleted: %s", key)
			return a.deleteWorkflow(key)
		},

		// Update function for an updated workflow object in the store.
		func(kv *v1alpha1.KVPairStruct, version uint64) error {
			log.Infof("workflow updated: %s", kv.Key)
			return a.addWorkflow(kv)
		},
	)
}

func (a *Controller) deleteWorkflow(key string) error {
	wfName := strings.TrimPrefix(key, v1alpha1.WorkflowKeyPrefix+"/")
	var (
		wfStatus v1alpha1.WorkflowStatus
		errs     = errors.NewMultiError()
	)

	val, err := store.KVStore.Get(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wfName))
	if err == nil {
		err := json.Unmarshal(val.Data, &wfStatus)
		if err != nil {
			return fmt.Errorf("error while unmarshaling the workflow status spec from data: %s", err)
		}

		var curWorkflow v1alpha1.Workflow
		err = json.Unmarshal([]byte(wfStatus.WorkflowSpec), &curWorkflow)
		if err != nil {
			return fmt.Errorf("error while unmarshaling the workflow spec from workflow status data: %s", err)
		}

		err = curWorkflow.Resolve()
		if err != nil {
			return fmt.Errorf("error while resolving workflow from status: %s", err)
		}

		for name, pipeline := range curWorkflow.Spec.Pipelines {
			err = a.Scheduler.RemovePipeline(wfName, name, &pipeline, &wfStatus)
			if err != nil {
				errs.Append(err)
			}
		}
	} else if store.KVStore.KeyDoesNotExistError(err) {
		log.Infof("WorkflowStatus not found for the workflow: %s, assuming already processed.", wfName)
		return nil
	} else {
		return fmt.Errorf("error while getting workflow status: %s", err)
	}

	// Finally remove the status associated with the workflow from the db.
	err = store.KVStore.Delete(context.TODO(), fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wfName))
	if err != nil {
		return fmt.Errorf("error while deleting workflow status object: %s", err)
	}

	return errs.GetError()
}

func (a *Controller) addWorkflow(kv *v1alpha1.KVPairStruct) error {
	var wf v1alpha1.Workflow
	err := json.Unmarshal([]byte(kv.Value), &wf)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the workflow spec from data: %s", err)
	}

	err = wf.Resolve()
	if err != nil {
		return fmt.Errorf("error while resolving workflow object: %s", err)
	}

	wfName := wf.Metadata.GetName()

	var (
		wfStatus      v1alpha1.WorkflowStatus
		updatedWfSpec v1alpha1.Workflow
		errs          []string
	)

	val, err := store.KVStore.Get(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wfName))
	if err == nil && !val.DeletedOrExpired {
		err := json.Unmarshal(val.Data, &wfStatus)
		if err != nil {
			return fmt.Errorf("error while unmarshaling the workflow status spec from data: %s", err)
		}

		var curWorkflow v1alpha1.Workflow
		err = json.Unmarshal([]byte(wfStatus.WorkflowSpec), &curWorkflow)
		if err != nil {
			return fmt.Errorf("error while unmarshaling the workflow spec from workflow status data: %s", err)
		}

		err = curWorkflow.Resolve()
		if err != nil {
			return fmt.Errorf("error while resolving workflow from status: %s", err)
		}

		// At this point, try best efforts for each pipeline.
		if wfStatus.Pipelines == nil {
			wfStatus.Pipelines = make(map[string]v1alpha1.PipelineStatus)
		}
		log.Infof("starting pipeline scheduling for workflow: %s", wf.Metadata.GetName())

		// First loop through all the pipelines which are supposed to be configured
		// using the workflow status manifest and remove them if they don't exist in the
		// new workflow.
		// Current running source of truth for the workflow state is from the curWorkflow.
		for name, pipeline := range curWorkflow.Spec.Pipelines {
			if _, ok := wfStatus.Pipelines[name]; !ok {
				wfStatus.Pipelines[name] = v1alpha1.PipelineStatus{}
			}

			var deleteTrigger bool
			if _, ok := wf.Spec.Pipelines[name]; !ok {
				ps := curWorkflow.GetTriggerAsssociatedPipelines(pipeline.TriggerName)
				deleteTrigger = true
				if len(ps) > 1 {
					deleteTrigger = false
				}
				err := a.Scheduler.RemovePipeline(wfName, name, &pipeline, &wfStatus)
				if err != nil {
					errs = append(errs, err.Error())
				} else {
					if deleteTrigger {
						delete(curWorkflow.Spec.Triggers, pipeline.TriggerName)
					}
					delete(curWorkflow.Spec.Pipelines, name)
				}
			}
		}

		// Process pipelines in the Workflow document.
		//
		// First range through all the pipelines configured in the new workflow manifest
		// then check which of these pipelines exists in the workflow registered
		// in the status.
		// If the pipeline exists in the workflow manifest in the status then check if the
		// two pipelines differ, if they do then do an update on the pipeline.
		// If the pipeline does not exist in the workflow manifest described in the statsus
		// then schedule the pipeline and update the status.
		for name, pipeline := range wf.Spec.Pipelines {
			if _, ok := wfStatus.Pipelines[name]; !ok {
				wfStatus.Pipelines[name] = v1alpha1.PipelineStatus{}
			}

			var err error
			if p, ok := curWorkflow.Spec.Pipelines[name]; ok {
				// This DeepEquals check also checks if the two pipelines
				// have the same trigger configured.
				if !pipeline.DeepEqual(&p) {
					err := a.Scheduler.UpdatePipeline(wfName, name, &pipeline, &p, &wfStatus)
					if err != nil {
						errs = append(errs, err.Error())
					} else {
						curWorkflow.Spec.Pipelines[name] = pipeline
						curWorkflow.Spec.Triggers[pipeline.TriggerName] = *pipeline.Trigger
					}
				}
			} else {
				err = a.Scheduler.SchedulePipeline(wfName, name, &pipeline, &wfStatus)
				if err != nil {
					errs = append(errs, err.Error())
				} else {
					curWorkflow.Spec.Pipelines[name] = pipeline
					curWorkflow.Spec.Triggers[pipeline.TriggerName] = *pipeline.Trigger
				}
			}
		}

		// Remove all the triggers which are not required.
		curWorkflow.RemoveNonLinkedTriggers()
		updatedWfSpec = curWorkflow

	} else if store.KVStore.KeyDoesNotExistError(err) {
		log.Infof("workflow status not found for: %s", wfName)
		log.Infof("scheduling the pipelines for the workflow")

		wfStatus, err := v1alpha1.NewWorkflowStatus(&wf)
		if err != nil {
			return err
		}

		for name, pipeline := range wf.Spec.Pipelines {
			err = a.Scheduler.SchedulePipeline(wfName, name, &pipeline, &wfStatus)
			if err != nil {
				errs = append(errs, err.Error())
				delete(wf.Spec.Pipelines, name)
			}
		}

		updatedWfSpec = wf
	} else {
		return err
	}

	// Save the workflow status.
	spec, err := json.Marshal(&updatedWfSpec)
	if err != nil {
		log.Errorf("error while marshaling workflow spec for status")
	} else {
		wfStatus.WorkflowSpec = string(spec)
	}

	statusVal, err := json.Marshal(&wfStatus)
	if err != nil {
		return fmt.Errorf("error while marshaling workflow status: %s", err)
	}

	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wfName),
		statusVal)
	if err != nil {
		return fmt.Errorf("error while updating workflow status: %s", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("error while operating pipelines: %v", strings.Join(errs, " :::: "))
	}

	return nil
}
