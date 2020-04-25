package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
		func(kv *v1alpha1.KVPairStruct) error {
			log.Infof("workflow deleted: %s", kv.Key)
			return nil
		},

		// Update function for an updated workflow object in the store.
		func(kv *v1alpha1.KVPairStruct, version uint64) error {
			log.Infof("workflow updated: %s", kv.Key)
			return a.addWorkflow(kv)
		},
	)
}

func (a *Controller) addWorkflow(kv *v1alpha1.KVPairStruct) error {
	var wf v1alpha1.Workflow
	err := json.Unmarshal([]byte(kv.Value), &wf)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the workflow spec from data: %s", err)
	}

	wfName := wf.Metadata.GetName()

	var (
		wfStatus      v1alpha1.WorkflowStatus
		updatedWfSpec v1alpha1.Workflow
		errs          []string
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

		if curWorkflow.Resolve() != nil || wf.Resolve() != nil {
			return fmt.Errorf("error while resolving workflows")
		}

		// At this point, try best efforts for each pipeline.
		log.Infof("starting pipeline scheduling for workflow: %s", wf.Metadata.GetName())
		for name, pipeline := range wf.Spec.Pipelines {
			if p, ok := curWorkflow.Spec.Pipelines[name]; ok {
				if !pipeline.DeepEqual(&p) {
					err := a.Scheduler.UpdatePipeline(wfName, name, &pipeline, &wfStatus)
					if err != nil {
						errs = append(errs, err.Error())
					} else {
						curWorkflow.Spec.Pipelines[name] = pipeline
					}
				}
			} else {
				err := a.Scheduler.SchedulePipeline(wfName, name, &pipeline, &wfStatus)
				if err != nil {
					errs = append(errs, err.Error())
				} else {
					curWorkflow.Spec.Pipelines[name] = pipeline
				}
			}
		}

		// remove pipelines which does not exist anymore.
		for name, pipeline := range curWorkflow.Spec.Pipelines {
			if _, ok := wf.Spec.Pipelines[name]; !ok {
				err := a.Scheduler.RemovePipeline(wfName, name, &pipeline, &wfStatus)
				if err != nil {
					errs = append(errs, err.Error())
				} else {
					curWorkflow.Spec.Pipelines[name] = pipeline
				}
			}
		}

		updatedWfSpec = curWorkflow

	} else {
		log.Debugf("workflow status not found for: %s", wfName)
		log.Debugf("scheduling the pipelines for the workflow")

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
