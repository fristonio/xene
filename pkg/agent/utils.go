package agent

import (
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/proto"
)

func getAgentSecretsList(agentName string) []*proto.AgentSecretInfo {
	return []*proto.AgentSecretInfo{
		{
			Name: fmt.Sprintf("%s-client-key", agentName),
			Type: defaults.SecretTypeCertificates,
		},
		{
			Name: fmt.Sprintf("%s-client-cert", agentName),
			Type: defaults.SecretTypeCertificates,
		},
		{
			Name: fmt.Sprintf("%s-root-ca-chain", agentName),
			Type: defaults.SecretTypeCertificates,
		},
	}
}
