package protoutils

import (
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
)

// UpdatePipelineStatusSpecFromGRPCTransport updates the pipeline status spec from the GRPC
// proto pipeline status.
func UpdatePipelineStatusSpecFromGRPCTransport(ps *v1alpha1.PipelineStatus, p *proto.PipelineStatus) {
	ps.Status = p.Status
}
