package option

// AgentOpts contains the desired configuration to run.
type AgentOpts struct {
	// Host is the host to bind the agent on.
	Host string `yaml:"host"`

	// Port is the port to run the agent on.
	Port uint32 `yaml:"port"`

	// APIServer is the address of the api server to join.
	APIServer string `yaml:"apiServer"`

	// APIAuthToken is the authentication token to use when joining the API server.
	APIAuthToken string `yaml:"apiAuthToken"`

	// Address is the network address of the agent itself, this is used for communication between
	// APIserver and agent.
	Address string `yaml:"address"`

	// Insecure denotes if the Agent should run in insecure mode.
	Insecure bool `yaml:"insecure"`

	// KeyFile is the key file to use for agent's secure GRPC scheme.
	KeyFile string `yaml:"keyFile"`

	// CertFile contains the certificate for agent's secure GRPC scheme.
	CertFile string `yaml:"certFile"`

	// JWTSecret contains the jwt secret for authorizing clients for GRPC
	// server
	JWTSecret string `yaml:"jwtSecret"`
}
