package v1alpha1

import (
	"fmt"
	"net/url"
)

// Executor is a type for representing an executor on the agent.
type Executor string

// Agent is the type which contains agent object definition.
type Agent struct {
	// TypeMeta stores meta type information for the agent object.
	TypeMeta `json:",inline"`

	// Metadata contains metadata about the Agent object.
	Metadata Metadata `json:"metadata"`

	// Spec contains the spec of the agent.
	Spec AgentSpec `json:"spec"`
}

// Validate checks for any issues in the information about the agent in the type.
func (a *Agent) Validate() error {
	if err := a.TypeMeta.Validate(AgentKind); err != nil {
		return err
	}

	if err := a.Metadata.Validate(); err != nil {
		return err
	}

	return a.Spec.Validate()
}

// DeepEquals checks if the object provided is equal to the agent object.
func (a *Agent) DeepEquals(ag *Agent) bool {
	if a.TypeMeta.DeepEquals(&ag.TypeMeta) &&
		a.Metadata.DeepEquals(&ag.Metadata.ObjectMeta) &&
		a.Spec.DeepEquals(&ag.Spec) {
		return true
	}

	return false
}

// AgentSpec contains the spec of the workflow.
type AgentSpec struct {
	// Address contains the address of the agent.
	Address string `json:"address"`

	// Insecure specifies if the agent is running insecure mode,
	// in that case we don't need certifictes for communication.
	Insecure bool `json:"insecure"`

	// ClientKeySecret is the xene client key registry secret object name to use
	// for the GRPC client for agent.
	ClientKeySecret string `json:"clientKeySecret"`

	// ClientCertSecret is the xene registry secret which contains the client certificate.
	ClientCertSecret string `json:"clientCertSecret"`

	// RootCASecret corresponds to the root certificate authority xene secret.
	RootCASecret string `json:"rootCA"`

	// ServerName is the hostname of the server, this is used in case
	// of mTLS authentication between apiserver and the agent.
	ServerName string `json:"serverName"`

	// LogServerEnabled specifies if the log server is enabled for the agent.
	LogServerEnabled bool `json:"logServerEnabled"`

	// LogServerPort contains the port of the http log file handler run with
	// the agent.
	LogServerPort uint32 `json:"logServerPort"`
}

// Validate validates the specification provided for the agent..
func (a *AgentSpec) Validate() error {
	_, err := url.Parse(a.Address)
	if err != nil {
		return fmt.Errorf("address is not valid: %s", err)
	}

	if !a.Insecure {
		if a.ClientCertSecret == "" || a.ClientKeySecret == "" || a.RootCASecret == "" || a.ServerName == "" {
			return fmt.Errorf("insecure mode, requires client and root CA certs to be specified")
		}
	}

	if a.LogServerEnabled {
		if a.LogServerPort == 0 {
			return fmt.Errorf("If the log server is enabled then port must be specified")
		}
	}

	return nil
}

// DeepEquals checks if the two AgentSpec objects are equal or not
func (a *AgentSpec) DeepEquals(az *AgentSpec) bool {
	if az.Address != a.Address ||
		a.Insecure != az.Insecure ||
		a.ClientKeySecret != az.ClientKeySecret ||
		a.ClientCertSecret != az.ClientCertSecret ||
		a.RootCASecret != az.RootCASecret ||
		a.ServerName != az.ServerName ||
		a.LogServerEnabled != az.LogServerEnabled ||
		a.LogServerPort != az.LogServerPort {
		return false
	}

	return true
}
