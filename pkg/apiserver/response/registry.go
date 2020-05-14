package response

// RegistryItemsFromPrefix is the response of list prefix query on registry items.
type RegistryItemsFromPrefix struct {
	Count int `json:"count" example:"2"`

	// Items contains the Serialized kvstore items
	Items []string `json:"items"`
}

// RegistryItem is the reponse of registry item get on the apiserver.
type RegistryItem struct {
	// Items contains the Serialized kvstore item
	Item string `json:"item"`
}

// AgentInfo contains a list of registry items keys
type AgentInfo struct {
	Name string `json:"name"`

	Available bool `json:"available"`

	Address string `json:"address"`

	Secure bool `json:"secure"`
}

// AgentTriggerInfo is the info about trigger on the agent.
type AgentTriggerInfo struct {
	Name string `json:"name"`

	Pipelines []string `json:"pipelines"`
}

// AgentWorkflowInfo contains the info about the workflow on the agent.
type AgentWorkflowInfo struct {
	Name string `json:"name"`

	Triggers []AgentTriggerInfo `json:"triggers"`
}

// AgentVerboseInfo contains verbose information about the agent.
type AgentVerboseInfo struct {
	Name string `json:"name"`

	Healthy bool `json:"healthy"`

	Address string `json:"address"`

	Workflows []AgentWorkflowInfo `json:"workflows"`
}

// SecretInfo contains information about a secret
type SecretInfo struct {
	Name string `json:"name"`

	Type string `json:"type"`

	Restricted bool `json:"restricted"`
}

// WorkflowInfo contains information about a xene workflow
type WorkflowInfo struct {
	Name string `json:"name"`

	Pipelines []string `json:"pipelines"`

	Triggers []string `json:"triggers"`

	Agents []string `json:"agents"`
}

// WorkflowVerboseInfo contains verbose information about a particular workflow.
type WorkflowVerboseInfo struct {
}
