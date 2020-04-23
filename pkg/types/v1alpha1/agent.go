package v1alpha1

import (
	"context"
	"fmt"
	"net/url"

	"github.com/fristonio/xene/pkg/proto"
	"google.golang.org/grpc"
)

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
func (a *Agent) CheckHealth() error {
	var opts []grpc.DialOption

	// Fix this to include secure connection using Mutual TLS.
	opts = append(opts, grpc.WithInsecure())

	addr, err := url.Parse(a.Spec.Address)
	if err != nil {
		return fmt.Errorf("error while parsing Agent address: %s", err)
	}
	conn, err := grpc.Dial(addr.Host, opts...)
	if err != nil {
		return fmt.Errorf("error failed to dial: %s", err)
	}
	defer conn.Close()

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

	// AuthToken contains the auth token to communicate with the running agent at the
	// provided address.
	AuthToken string `json:"authToken"`
}

// Validate validates the specification provided for the agent..
func (a *AgentSpec) Validate() error {
	_, err := url.Parse(a.Address)
	if err != nil {
		return fmt.Errorf("address is not valid: %s", err)
	}

	return nil
}
