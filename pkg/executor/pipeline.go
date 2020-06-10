package executor

import (
	"fmt"
	"time"

	"github.com/fristonio/xene/pkg/dag"
	"github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/executor/cre"
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

	// useStore specifies wheather to use kvstore interaction during the execution
	// of the pipeline.
	useStore bool

	// status contains the status of the pipeline execution
	// This is only set to a value if we are not using the store
	// for save run status.
	status *v1alpha1.PipelineRunStatus
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

// WithoutStore sets the boolean useStore to false which will make the executor
// not perform any KVStore interactions
func (p *PipelineExecutor) WithoutStore() *PipelineExecutor {
	p.re.WithoutStore()
	return p
}

// Run starts running the pipeline.
// Make sure that Pipeline status contains dummy information about all the tasks
// and step in the pipeline
// it should have a blueprint of the pipeline spec.
func (p *PipelineExecutor) Run(status v1alpha1.PipelineRunStatus) {
	p.log.Debugf("running PipelineExecutor")
	p.status = &status

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

	// Associate status with the runtime executor
	p.re.WithStatus(&status)

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
		err := p.re.RunTask(task.Name(), task)
		if err != nil {
			errs.Append(fmt.Errorf("error while running the task: %s", err))
			p.log.Errorf("error running task(%s): %s", task.Name(), err)
		}

		return errs
	})

	if len(walkErrors) > 0 {
		status.Status = v1alpha1.StatusError
		p.log.Errorf("error while walking task graph: \n%v", walkErrors)
	} else {
		status.Status = v1alpha1.StatusSuccess
	}

	status.EndTime = time.Now().Unix()
	p.re.SaveStatusToStore()

	err = p.re.Shutdown()
	if err != nil {
		p.log.Errorf("error while shutting down runtime executor: %s", err)
		return
	}
}

// GetStatus returns the status of the pipeline run.
func (p *PipelineExecutor) GetStatus() *v1alpha1.PipelineRunStatus {
	return p.status
}
