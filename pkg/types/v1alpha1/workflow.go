package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/fristonio/xene/pkg/dag"
	"github.com/fristonio/xene/pkg/utils"
)

// Workflow is the type which contains workflow object definition.
type Workflow struct {
	// TypeMeta stores meta type information for the workflow object.
	TypeMeta `json:",inline"`

	// Metadata contains metadata about the Workflow object.
	Metadata Metadata `json:"metadata"`

	// Spec contains the spec of the workflow
	Spec WorkflowSpec `json:"spec"`
}

// Validate checks for any issues in the information about the workflow in the type.
func (w *Workflow) Validate() error {
	if err := w.TypeMeta.Validate(WorkflowKind); err != nil {
		return err
	}

	if err := w.Metadata.Validate(); err != nil {
		return err
	}

	if err := w.Resolve(); err != nil {
		return err
	}

	return w.Spec.Validate()
}

// Resolve resolves the workflow object by acting upon different relations
// in the spec.
func (w *Workflow) Resolve() error {
	for name, pipeline := range w.Spec.Pipelines {
		if trigger, ok := w.Spec.Triggers[pipeline.TriggerName]; ok {
			pipeline.Trigger = &trigger
		} else {
			return fmt.Errorf("not a valid trigger(%s) for pipeline(%s)", pipeline.TriggerName, name)
		}

		if err := pipeline.Resolve(name); err != nil {
			return err
		}
	}

	return nil
}

// WorkflowSpec contains the spec of the workflow.
type WorkflowSpec struct {
	// Triggers contains a list of trigger.
	Triggers map[string]TriggerSpec `json:"triggers"`

	// Pipelines contains a list of pipeline configured with workflow.
	Pipelines map[string]PipelineSpec `json:"pipelines"`
}

// Validate validates the specification provided for the workflow.
func (w *WorkflowSpec) Validate() error {
	if w.Triggers == nil || w.Pipelines == nil {
		return fmt.Errorf("at least one pipeline and one trigger must be mentioned")
	}
	for name, trigger := range w.Triggers {
		if err := trigger.Validate(name); err != nil {
			return err
		}
	}

	for name, pipeline := range w.Pipelines {
		if err := pipeline.Validate(name); err != nil {
			return err
		}
	}

	return nil
}

// TriggerSpec contains spec of a trigger for xene workflow.
type TriggerSpec struct {
	// Type contains the type of the trigger we are using
	Type string `json:"type"`
}

// Validate validates the specification provided for the a trigger.
func (t *TriggerSpec) Validate(name string) error {
	return nil
}

// DeepEqual checks if the two Trigger objects are equal or not.
func (t *TriggerSpec) DeepEqual(tz *TriggerSpec) bool {
	return t.Type == tz.Type
}

// PipelineSpec contains the spec of a pipeline associated with the workflow.
type PipelineSpec struct {
	// trigger contains the Trigger for the configured pipeline.
	Trigger *TriggerSpec `json:"-"`

	// Dag contains the dag corresponding to tasks in a pipeline.
	Dag *dag.AcyclicGraph `json:"-"`

	// TriggerName contains the name of the trigger to use for the pipeline.
	TriggerName string `json:"trigger"`

	// Description contains description about the pipeline.
	Description string `json:"description"`

	// Executor describes executor for the pipeline.
	// Should be one of the preconfigured list of available executors.
	Executor string `json:"executor"`

	// RootTask contains the root task for the pipeline.
	RootTask dag.Vertex `json:"-"`

	// Tasks contains the list of the tasks in the pipeline.
	Tasks map[string]TaskSpec `json:"tasks"`
}

// DeepEqual checks if the two pipeline objects are equal or not.
func (p *PipelineSpec) DeepEqual(pz *PipelineSpec) bool {
	if p.TriggerName != pz.TriggerName {
		return false
	}

	if p.Trigger != nil && pz.Trigger != nil {
		return p.Trigger.DeepEqual(pz.Trigger)
	}

	if p.Trigger != nil || pz.Trigger != nil {
		return false
	}

	if p.Executor != pz.Executor {
		return false
	}

	for name, task := range p.Tasks {
		pzTask, ok := pz.Tasks[name]
		if !ok {
			return false
		}

		if !task.DeepEqual(&pzTask) {
			return false
		}
	}
	return true
}

// Validate validates the specification provided for the pipeline.
func (p *PipelineSpec) Validate(name string) error {
	return nil
}

// Resolve resolves the specification provided for the pipeline.
func (p *PipelineSpec) Resolve(pipelineName string) error {
	p.Dag = dag.NewAcyclicGraph()
	for name, task := range p.Tasks {
		task.Resolve(name)
		p.Dag.Add(&task)
	}

	for name, task := range p.Tasks {
		for _, depName := range task.Dependencies {
			t, ok := p.Tasks[depName]
			if !ok {
				return fmt.Errorf("pipeline %s: not a valid task dependency %s for task  %s", pipelineName, depName, name)
			}

			task.DependsOn = append(task.DependsOn, &t)
			p.Dag.Connect(dag.BasicEdge(&task, &t))
		}
	}

	root, err := p.Dag.Root()
	if err != nil {
		return fmt.Errorf("pipeline %s: error while resolving root: %s", pipelineName, err)
	}

	p.RootTask = root
	return nil
}

// TaskSpec contains the spec corresponding to a single task in a pipeline.
type TaskSpec struct {
	name string

	// Description contains the description of the task
	Description string `json:"description"`

	// DependsOn contains a list of tasks this task depends on.
	// This is used to build the dag for the pipeline.
	DependsOn []*TaskSpec `json:"-"`

	// Dependencies is a list of task names which this particular
	// task depends on. It is used to resolve the DependsOn variable in
	// the struct.
	Dependencies []string `json:"dependencies"`

	// WorkingDirectory for the task.
	WorkingDirectory string `json:"workingDir"`

	// Step contains a list of steps to go through for the task.
	// These are executed linear.
	Steps []TaskStepSpec `json:"steps"`
}

// DeepEqual checks if the two TaskSpec objects are equal or not.
func (t *TaskSpec) DeepEqual(tz *TaskSpec) bool {
	if t.WorkingDirectory != tz.WorkingDirectory {
		return false
	}

	if utils.CheckStringSliceEqual(t.Dependencies, tz.Dependencies) {
		return false
	}

	if len(t.Steps) != len(tz.Steps) {
		return false
	}

	for i, step := range t.Steps {
		if !step.DeepEqual(&tz.Steps[i]) {
			return false
		}
	}

	return false
}

// Resolve resolves the task specification.
func (t *TaskSpec) Resolve(name string) {
	t.name = name
}

// Hashcode returns a string representation of the task, which is uniquely defined by
// its name.
func (t *TaskSpec) Hashcode() string {
	return t.name
}

// TaskStepSpec contains the specification of an individual TaskStep
type TaskStepSpec struct {
	// Name contains the name of the Step
	Name string `json:"name"`

	// Description contains the description of the step, this is
	// optional.
	Description string `json:"description"`

	// Type contains the type of the Step we are exeucting.
	Type string `json:"type"`

	// Cmd defines the command to execute for a step, when the
	// type of step is shell.
	Cmd string `json:"cmd"`
}

// DeepEqual checks if the two TaskSteppec objects are equal or not.
func (t *TaskStepSpec) DeepEqual(tz *TaskStepSpec) bool {
	if t.Type != tz.Type &&
		t.Cmd != tz.Cmd {
		return false
	}

	return true
}

// WorkflowStatus contains the status corresponding to a defined workflow.
type WorkflowStatus struct {
	// TypeMeta stores meta type information for the WorkflowStatus object.
	TypeMeta `json:",inline"`

	// Metadata contains metadata about the WorkflowStatus object.
	Metadata Metadata `json:"metadata"`

	// WorkflowSpec contains the workflow spec currently being executed
	// by xene.
	WorkflowSpec string `json:"workflowSpec"`

	// Pipelines contains the status of all the pipelines.
	Pipelines map[string]PipelineStatus `json:"pipelines"`
}

// Validate checks for any issues in the information about the workflow status in the type.
func (w *WorkflowStatus) Validate() error {
	if err := w.TypeMeta.Validate(WorkflowStatusKind); err != nil {
		return err
	}

	if err := w.Metadata.Validate(); err != nil {
		return err
	}

	return nil
}

// NewWorkflowStatus returns a new WorkflowStatus object using the workflow
// specification provided.
func NewWorkflowStatus(wf *Workflow) (WorkflowStatus, error) {
	data, err := json.Marshal(wf)
	if err != nil {
		return WorkflowStatus{}, fmt.Errorf("error while marshaling wf spec: %s", err)
	}

	return WorkflowStatus{
		TypeMeta: TypeMeta{
			Kind:       WorkflowStatusKind,
			APIVersion: "v1alpha1",
		},
		Metadata: Metadata{
			ObjectMeta: ObjectMeta{
				Name: wf.Metadata.GetName(),
			},
		},
		WorkflowSpec: string(data),
		Pipelines:    make(map[string]PipelineStatus),
	}, nil
}

// PipelineStatus contains the status of the pipeline in context
type PipelineStatus struct {
	// LastRunTimestamp contains the timestamp of the time when the pipeline
	// was last run.
	LastRunTimestamp int64 `json:"lastRunTimestamp"`

	// Executor is the name of the agent which ran this pipeline
	Executor string `json:"executor"`

	// Status contains the status information of the pipeline.
	Status string `json:"status"`
}

// TriggerType is the type to specify the type of trigger.
type TriggerType string
