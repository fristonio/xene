package defaults

var (
	// AgentHost is the default host to bind the Agent to.
	AgentHost string = "0.0.0.0"

	// AgentPort is the port to bind the agent GRPC server on.
	AgentPort uint32 = 6061

	// AgentStorageDir is the storage directory for xene agent
	AgentStorageDir string = "/var/run/xene/agent/store/"
)
