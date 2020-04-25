package protoutils

import (
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
)

// GRPCPipelineStatusToTypeSpec converts the grpc wire PipelineStatus type to the one
// corresponding to xene types.
func GRPCPipelineStatusToTypeSpec(p *proto.PipelineStatus) *v1alpha1.PipelineStatus {
	return &v1alpha1.PipelineStatus{
		LastRunTimestamp: p.LastRunTimestamp,
		Executor:         p.Executor,
		Status:           p.Status,
	}
}

// UpdatePipelineStatusSpecFromGRPCTransport updates the pipeline status spec from the GRPC
// proto pipeline status.
func UpdatePipelineStatusSpecFromGRPCTransport(ps *v1alpha1.PipelineStatus, p *proto.PipelineStatus) {
	ps.LastRunTimestamp = p.LastRunTimestamp
	ps.Executor = p.Executor
	ps.Status = p.Status
}
