package response

// KVPair contains a key value pair.
type KVPair struct {
	Key string `json:"key"`

	Value string `json:"value"`
}

// RegistryItemsFromPrefix is the response of list prefix query on registry items.
type RegistryItemsFromPrefix struct {
	Count int `json:"count" example:"2"`

	// Items contains the Serialized kvstore items
	Items []KVPair `json:"items"`
}

// RegistryItem is the reponse of registry item get on the apiserver.
type RegistryItem struct {
	// Items contains the Serialized kvstore item
	Item KVPair `json:"item"`
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

// AgentSecretInfo contains the info about the secret for the agent.
type AgentSecretInfo struct {
	Name string `json:"name"`

	Type string `json:"type"`
}

// AgentVerboseInfo contains verbose information about the agent.
type AgentVerboseInfo struct {
	Name string `json:"name"`

	Healthy bool `json:"healthy"`

	Address string `json:"address"`

	Secure bool `json:"secure"`

	ServerName string `json:"serverName"`

	Workflows []AgentWorkflowInfo `json:"workflows"`

	Secrets []AgentSecretInfo `json:"secrets"`
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

// PipelineRunInfo contains information about a particular pipeline run.
type PipelineRunInfo struct {
	RunID string `json:"runID"`

	Status string `json:"status"`

	Agent string `json:"agent"`

	StartTime int64 `json:"startTime"`

	EndTime int64 `json:"endTime"`
}

// PipelineRunVerboseInfo contains verbose information about a particular pipeline run.
type PipelineRunVerboseInfo struct {
	RunID string `json:"runID"`

	Status string `json:"status"`

	RunInfo string `json:"runInfo"`

	BaseLogURL string `json:"baseLogURL"`
}

// PipelineInfo contains verbose information about a particular pipeline.
type PipelineInfo struct {
	Workflow string `json:"workflow"`

	Name string `json:"name"`

	Spec string `json:"spec"`

	Runs []PipelineRunInfo `json:"runs"`

	Warnings []string `json:"warnings"`
}
