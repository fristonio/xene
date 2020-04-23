package v1alpha1

import "fmt"

var (
	// Version is the version string for the types registered for xene.
	Version string = "v1alpha1"

	// WorkflowKeyPrefix is the key prefix for all the workflow objects.
	WorkflowKeyPrefix string = fmt.Sprintf("%s/registry/workflow", Version)

	// WorkflowKind is the Kind when creating a Workflow
	WorkflowKind string = "workflow"

	// AgentKeyPrefix is the key prefix for all the agent objects in xene.
	AgentKeyPrefix string = fmt.Sprintf("%s/registry/agent", Version)

	// AgentKind is the kind when creating Agent object
	AgentKind string = "agent"

	// SecretKeyPrefix is the key prefix for all the workflow objects.
	SecretKeyPrefix string = fmt.Sprintf("%s/registry/secret", Version)

	// SecretKind is the kind when creating a secret object in xene.
	SecretKind string = "secret"
)
