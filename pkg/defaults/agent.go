package defaults

var (
	// AgentHost is the default host to bind the Agent to.
	AgentHost string = "0.0.0.0"

	// AgentPort is the port to bind the agent GRPC server on.
	AgentPort uint32 = 6061

	// AgentStorageDir is the storage directory for xene agent
	AgentStorageDir string = "/var/run/xene/agent/store/"

	// AgentVersion contains the version information regarding the agent.
	AgentVersion string = "v1alpha1"

	// PipelineIDSize contains the size of the string representing the
	// Pipeline ID.
	PipelineIDSize uint32 = 20

	// AgentLogsDir is the path containing the logs for the agent.
	AgentLogsDir string = "/var/run/xene/agent/logs/"
)
