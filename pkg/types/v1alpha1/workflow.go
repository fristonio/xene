package v1alpha1

import (
	"encoding/json"
	"fmt"
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

// Resolve resolves the workflow object by acting upon different relations
// in the spec.
func (w *Workflow) Resolve() error {
	for name, pipeline := range w.Spec.Pipelines {
		if trigger, ok := w.Spec.Triggers[pipeline.TriggerName]; ok {
			pipeline.Trigger = &trigger
		} else {
			return fmt.Errorf("not a valid trigger(%s) for pipeline(%s)", pipeline.Trigger, name)
		}
	}

	return nil
}

// WorkflowSpec contains the spec of the workflow.
type WorkflowSpec struct {
	// Triggers contains a list of trigger.
	Triggers map[string]TriggerSpec `json:"trigger"`

	// Pipelines contains a list of pipeline configured with workflow.
	Pipelines map[string]PipelineSpec `json:"pipelines"`
}

// TriggerSpec contains spec of a trigger for xene workflow.
type TriggerSpec struct {
	// Type contains the type of the trigger we are using
	Type string `json:"type"`
}

// DeepEqual checks if the two Trigger objects are equal or not.
func (t *TriggerSpec) DeepEqual(tz *TriggerSpec) bool {
	return t.Type == tz.Type
}

// PipelineSpec contains the spec of a pipeline associated with the workflow.
type PipelineSpec struct {
	// trigger contains the Trigger for the configured pipeline.
	Trigger *TriggerSpec

	// Trigger contains the name of the trigger to use for the pipeline.
	TriggerName string `json:"trigger"`
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
