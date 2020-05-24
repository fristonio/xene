package v1alpha1

import (
	"fmt"
)

var (
	// Version is the version string for the types registered for xene.
	Version string = "v1alpha1"

	// RegistryPath contains the path of the registry key
	RegistryPath string = fmt.Sprintf("%s/registry", Version)

	// WorkflowKeyPrefix is the key prefix for all the workflow objects.
	WorkflowKeyPrefix string = fmt.Sprintf("%s/workflow", RegistryPath)

	// WorkflowStatusKeyPrefix is the key prefix for all the workflow status objects.
	WorkflowStatusKeyPrefix string = fmt.Sprintf("%s/status/workflow", RegistryPath)

	// PipelineStatusKeyPrefix is the key prefix for all the pipeline status objects.
	PipelineStatusKeyPrefix string = fmt.Sprintf("%s/status/pipeline", RegistryPath)

	// PipelineKeyPrefix is the key prefix for all the pipeline objects.
	PipelineKeyPrefix string = fmt.Sprintf("%s/pipeline", RegistryPath)

	// TriggerStatusKeyPrefix is the key prefix for all the trigger status objects.
	TriggerStatusKeyPrefix string = fmt.Sprintf("%s/status/trigger", RegistryPath)

	// TriggerKeyPrefix is the key prefix for all the trigger objects.
	TriggerKeyPrefix string = fmt.Sprintf("%s/trigger", RegistryPath)

	// WorkflowKind is the Kind when creating a Workflow
	WorkflowKind string = "Workflow"

	// WorkflowStatusKind is the kind corresponding to WorkflowStatus
	WorkflowStatusKind string = "WorkflowStatus"

	// AgentKeyPrefix is the key prefix for all the agent objects in xene.
	AgentKeyPrefix string = fmt.Sprintf("%s/agent", RegistryPath)

	// AgentKind is the kind when creating Agent object
	AgentKind string = "Agent"

	// SecretKeyPrefix is the key prefix for all the workflow objects.
	SecretKeyPrefix string = fmt.Sprintf("%s/secret", RegistryPath)

	// SecretKind is the kind when creating a secret object in xene.
	SecretKind string = "Secret"

	// DefaultTriggerType contains the name of the default trigger type.
	DefaultTriggerType TriggerType = "default"

	// CronTriggerType contains the name of the cron job trigger type.
	CronTriggerType TriggerType = "cron"

	// DockerExecutor is the name of the docker executor used by agent.
	DockerExecutor Executor = "docker"

	// ContainerExecutor is the name of the container executor used by agent.
	ContainerExecutor Executor = "container"

	// RegisteredExecutors contains a list of executors ContainerExecutor with xene.
	RegisteredExecutors []Executor = []Executor{ContainerExecutor}

	// RegisteredTriggerTypes contains a list of TriggerTypes allowed with xene.
	RegisteredTriggerTypes []TriggerType = []TriggerType{DefaultTriggerType, CronTriggerType}

	// RegistryItemWorkflow contains the workflow registry item name
	RegistryItemWorkflow RegistryItem = "workflow"

	// RegistryItemAgent contains the Agent registry item name
	RegistryItemAgent RegistryItem = "agent"

	// RegistryItemSecret contains the secret registry item name
	RegistryItemSecret RegistryItem = "secret"

	// RegistryItems contains a list of all the available registry items
	RegistryItems []RegistryItem = []RegistryItem{RegistryItemAgent, RegistryItemSecret, RegistryItemWorkflow}

	// StatusError contains the error status
	StatusError string = "Error"

	// StatusSuccess contains the success status
	StatusSuccess string = "Success"

	// StatusRunning contains the running  status
	StatusRunning string = "Running"

	// StatusNotExecuted contais the status when a component is not executed.
	StatusNotExecuted string = "NotExecuted"
)
