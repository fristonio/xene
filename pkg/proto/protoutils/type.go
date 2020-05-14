package protoutils

import (
	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
)

// UpdatePipelineStatusSpecFromGRPCTransport updates the pipeline status spec from the GRPC
// proto pipeline status.
func UpdatePipelineStatusSpecFromGRPCTransport(ps *v1alpha1.PipelineStatus, p *proto.PipelineStatus) {
	ps.Status = p.Status
}

// GetAgentVerboseInfoFromProtoAgentInfo returns the info of the agent from the proto information.
func GetAgentVerboseInfoFromProtoAgentInfo(info *proto.AgentInfo) *response.AgentVerboseInfo {
	var res = response.AgentVerboseInfo{
		Name:    info.Name,
		Healthy: info.Healthy,
		Address: info.Address,
	}

	wf := []response.AgentWorkflowInfo{}

	for _, w := range info.Workflows {
		triggers := []response.AgentTriggerInfo{}
		for _, trigger := range w.Triggers {
			triggers = append(triggers, response.AgentTriggerInfo{
				Name:      trigger.Name,
				Pipelines: trigger.Pipelines,
			})
		}
		wf = append(wf, response.AgentWorkflowInfo{
			Name:     w.Name,
			Triggers: triggers,
		})
	}

	res.Workflows = wf
	return &res
}
