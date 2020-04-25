package v1alpha1

import "fmt"

var (
	// Version is the version string for the types registered for xene.
	Version string = "v1alpha1"

	// WorkflowKeyPrefix is the key prefix for all the workflow objects.
	WorkflowKeyPrefix string = fmt.Sprintf("%s/registry/workflow", Version)

	// WorkflowStatusKeyPrefix is the key prefix for all the workflow status objects.
	WorkflowStatusKeyPrefix string = fmt.Sprintf("%s/status/workflow", Version)

	// WorkflowKind is the Kind when creating a Workflow
	WorkflowKind string = "Workflow"

	// WorkflowStatusKind is the kind corresponding to WorkflowStatus
	WorkflowStatusKind string = "WorkflowStatus"

	// AgentKeyPrefix is the key prefix for all the agent objects in xene.
	AgentKeyPrefix string = fmt.Sprintf("%s/registry/agent", Version)

	// AgentKind is the kind when creating Agent object
	AgentKind string = "Agent"

	// SecretKeyPrefix is the key prefix for all the workflow objects.
	SecretKeyPrefix string = fmt.Sprintf("%s/registry/secret", Version)

	// SecretKind is the kind when creating a secret object in xene.
	SecretKind string = "Secret"
)
