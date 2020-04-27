package agent

import (
	"context"

	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/proto"
	log "github.com/sirupsen/logrus"
)

type agentServer struct {
	name string
}

func newAgentServer() *agentServer {
	return &agentServer{
		name: option.Config.Agent.Name,
	}
}

// Status returns the status of the agent running.
func (a *agentServer) Status(ctx context.Context, opts *proto.StatusOpts) (*proto.AgentStatus, error) {
	log.Debugf("rpc to check status of the agent running")
	return &proto.AgentStatus{
		Healthy: true,
	}, nil
}

// SchedulePipeline is the RPC to schedule a pipeline on to the agent
func (a *agentServer) SchedulePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	log.Debugf("rpc to schedule pipeline on the agent: %s", pipeline.Spec)
	return &proto.PipelineStatus{
		Status:   "Not Implemented",
		Executor: a.name,
	}, nil
}

// UpdatePipeline is the RPC to update a pipeline on to the agent
func (a *agentServer) UpdatePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	log.Debugf("rpc to update pipeline on the agent: %s", pipeline.Spec)
	return &proto.PipelineStatus{
		Status:   "Not Implemented",
		Executor: a.name,
	}, nil
}

// RemovePipeline is the RPC to remove a pipeline from the agent.
func (a *agentServer) RemovePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	log.Debugf("rpc to remove pipeline from the agent: %s", pipeline.Spec)
	return &proto.PipelineStatus{
		Status:   "Not Implemented",
		Executor: a.name,
	}, nil
}
