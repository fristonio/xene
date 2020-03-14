package v1alpha1

// Workflow is the type which contains workflow object definition.
type Workflow struct {
	// TypeMeta stores meta type information for the workflow object.
	TypeMeta `json:",inline"`

	// Metadata contains metadata about the Workflow object.
	Metadata Metadata `json:"metadata"`

	// Spec contains the spec of the workflow
	Spec WorkflowSpec `json:"spec"`
}

// WorkflowSpec contains the spec of the workflow.
type WorkflowSpec struct {
	// Triggers contains a list of trigger.
	Triggers []Trigger `json:"trigger"`

	// Pipelines contains a list of pipeline configured with workflow.
	Pipelines []Pipeline `json:"pipelines"`
}

// Trigger contains spec of a trigger for xene workflow.
type Trigger struct {
}

// Pipeline contains the spec of a pipeline associated with the workflow.
type Pipeline struct {
}
