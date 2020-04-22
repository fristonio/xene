package agent

import (
	"context"

	"github.com/fristonio/xene/pkg/proto"
	log "github.com/sirupsen/logrus"
)

type agentServer struct{}

func newAgentServer() *agentServer {
	return &agentServer{}
}

// Status returns the status of the agent running.
func (a *agentServer) Status(ctx context.Context, opts *proto.StatusOpts) (*proto.AgentStatus, error) {
	log.Debug("rpc to check status of the agent running.")
	return &proto.AgentStatus{
		Healthy: true,
	}, nil
}

// SchedulePipeline is the RPC to schedule a pipeline on to the agent
func (a *agentServer) SchedulePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	log.Debug("rpc to schedule pipeline on the agent")
	return &proto.PipelineStatus{
		Status: "Not Implemented",
	}, nil
}
