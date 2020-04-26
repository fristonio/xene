package v1alpha1

import (
	"context"
	"fmt"
	"net/url"

	"github.com/fristonio/xene/pkg/proto"
	"google.golang.org/grpc"
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

// CheckHealth checks for the health status of the provided agent in context.
// If the agent is not healthy or there is some issue with the connectivity the
// function returns an error.
func (a *Agent) CheckHealth(conn *grpc.ClientConn) error {
	client := proto.NewAgentServiceClient(conn)
	status, err := client.Status(context.TODO(), &proto.StatusOpts{Verbose: true})
	if err != nil {
		return fmt.Errorf("error while fetching status: %s", err)
	}

	if !status.Healthy {
		return fmt.Errorf("agent is not healthy")
	}

	return nil
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

	return nil
}
