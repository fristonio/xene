package response

import (
	types "github.com/fristonio/xene/pkg/types/v1alpha1"
)

// WorkflowsFromPrefix is the response of list prefix query on workflows.
type WorkflowsFromPrefix struct {
	Count int `json:"count" example:"2"`

	Workflows []types.KVPairStruct `json:"workflows"`
}

// Workflow is the reponse of single workflow get on the apiserver.
type Workflow struct {
	Workflow types.KVPairStruct `json:"workflow" example:"Workflow Document"`
}
