package v1alpha1

// Agent is the type which contains agent object definition.
type Agent struct {
	// TypeMeta stores meta type information for the agent object.
	TypeMeta `json:",inline"`

	// Metadata contains metadata about the Agent object.
	Metadata Metadata `json:"metadata"`

	// Spec contains the spec of the agent.
	Spec AgentSpec `json:"spec"`
}

// AgentSpec contains the spec of the workflow.
type AgentSpec struct {
	// Address contains the address of the agent.
	Address string `json:"address"`

	// AuthToken contains the auth token to communicate with the running agent at the
	// provided address.
	AuthToken string `json:"authToken"`
}
