package v1alpha1

import "fmt"

var (
	// Version is the version string for the types registered for xene.
	Version string = "v1alpha1"

	// WorkflowKeyPrefix is the key prefix for all the workflow objects.
	WorkflowKeyPrefix string = fmt.Sprintf("%s/registry/workflow", Version)

	// AgentKeyPrefix is the key prefix for all the agent objects in xene.
	AgentKeyPrefix string = fmt.Sprintf("%s/registry/agent", Version)
)
