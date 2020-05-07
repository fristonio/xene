package executor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fristonio/xene/pkg/dag"
	"github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/executor/cre"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/sirupsen/logrus"
)

// PipelineExecutor is the type for executing pipelines.
type PipelineExecutor struct {
	// name contains the name of the pipeline prefixed with the
	// workflow name the pipeline is associated with.
	name string

	// id contains the ID of the pipeline we are currently executing.
	id string

	// Spec contains the specification of the pipeline.
	Spec *v1alpha1.PipelineSpecWithName

	// log contains the logger for the pipeline executor
	log *logrus.Entry

	// re contains the runtime executor for the pipeline.
	re RuntimeExecutor
}

// NewPipelineExecutor returns a new instance of PipelineExecutor for the provided
// specification.
func NewPipelineExecutor(name, id string, spec *v1alpha1.PipelineSpecWithName) *PipelineExecutor {
	var re RuntimeExecutor

	switch v1alpha1.Executor(spec.Executor.Type) {
	case v1alpha1.ContainerExecutor:
		re = cre.NewCRExecutor(string(v1alpha1.DockerExecutor), id, name, spec)
	default:
		logrus.Warnf("not a valid executor: %s, using default", spec.Executor.Type)
		re = cre.NewCRExecutor(string(v1alpha1.DockerExecutor), id, name, spec)
	}

	return &PipelineExecutor{
		Spec: spec,
		log: logrus.WithFields(logrus.Fields{
			"pipeline": name,
		}),
		name: name,
		id:   id,
		re:   re,
	}
}

// Run starts running the pipeline.
func (p *PipelineExecutor) Run() {
	p.log.Debugf("running PipelineExecutor")

	status := v1alpha1.NewPipelineRunStatus()
	status.Name = p.name
	status.RunID = p.id

	// Configure the pipeline runtime executor.
	err := p.re.Configure()
	if err != nil {
		p.log.Errorf("error while setting up runtime executor: %s", err)
		return
	}

	// This also transitively reduces the DAG we have made for the tasks
	err = p.Spec.Resolve(p.name)
	if err != nil {
		p.log.Errorf("error while resolving pipeline: %s", err)
	}

	// Walk each of task in the pipeline in the required order.
	walkErrors := p.Spec.Dag.Walk(func(v dag.Vertex) *errors.MultiError {
		errs := errors.NewMultiError()
		task, ok := v.(*v1alpha1.TaskSpec)
		if !ok {
			// If any error is getting the task spec then append it to the list of errors
			errs.Append(fmt.Errorf("not a valid vertex to visit, must confirm to type *TaskSpec"))
			return errs
		}

		// Run the task using the runtime executor.
		tStatus, err := p.re.RunTask(task.Name(), task)
		if err != nil {
			errs.Append(fmt.Errorf("error while running the task: %s", err))
			p.log.Errorf("error running task(%s): %s", task.Name(), err)
			tStatus.Status = v1alpha1.StatusError
		} else {
			tStatus.Status = v1alpha1.StatusSuccess
		}

		status.Tasks[task.Name()] = tStatus
		return errs
	})

	// Put dummy information in the status for any task which might not have been
	// executed
	for name, task := range p.Spec.Tasks {
		if _, ok := status.Tasks[name]; !ok {

			stepStatus := make(map[string]*v1alpha1.StepRunStatus)
			for _, step := range task.Steps {
				stepStatus[step.Name] = &v1alpha1.StepRunStatus{
					Status: v1alpha1.StatusNotExecuted,
				}
			}

			status.Tasks[name] = &v1alpha1.TaskRunStatus{
				Status: v1alpha1.StatusNotExecuted,
				Steps:  stepStatus,
			}
		}
	}

	if len(walkErrors) > 0 {
		status.Status = v1alpha1.StatusError
		p.log.Errorf("error while walking task graph: \n%v", walkErrors)
	} else {
		status.Status = v1alpha1.StatusSuccess
	}

	err = p.saveStatusToStore(&status)
	if err != nil {
		p.log.Error(err)
	}

	err = p.re.Shutdown()
	if err != nil {
		p.log.Errorf("error while shutting down runtime executor: %s", err)
		return
	}
}

func (p *PipelineExecutor) saveStatusToStore(status *v1alpha1.PipelineRunStatus) error {
	statusKey := fmt.Sprintf("%s/%s/%s", v1alpha1.PipelineStatusKeyPrefix, p.name, p.id)

	val, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("error while marshalling status: %s", err)
	}

	err = store.KVStore.Set(context.TODO(), statusKey, val)
	if err != nil {
		return fmt.Errorf("error while setting status key in kvstore: %s", err)
	}

	return nil
}
