package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fristonio/xene/pkg/apiserver/controller/agent"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/proto/protoutils"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// Scheduler is the scheduler used for managing workflows execution on different
// agent.
type Scheduler struct {
	// AgentController is for interacting with agent configured for xene.
	AgentController *agent.Controller
}

// Action is a type specifying action of an operation
type Action string

var (
	// ActionTypeUpdate is the action name corresponding to an update operation.
	ActionTypeUpdate Action = "update"

	// ActionTypeRemove is the action name corresponding to a remove operation.
	ActionTypeRemove Action = "remove"

	// ActionTypeSchedule is the action name corresponding to schedule operation.
	ActionTypeSchedule Action = "schedule"
)

// NewSchedulerWithDefaultAgentCtrl returns a new wofklow scheduler with the globally
// configured AgentController.
func NewSchedulerWithDefaultAgentCtrl() *Scheduler {
	return &Scheduler{
		AgentController: agent.AgentCtrl,
	}
}

func (s *Scheduler) performPipelineAction(action Action, wfName, name string,
	pipeline *v1alpha1.PipelineSpec, status *v1alpha1.PipelineStatus) error {

	agentName := status.Executor
	conn := s.AgentController.AgentConnection(agentName)
	if conn == nil {
		return fmt.Errorf("error while retriving agent connection")
	}

	data, err := json.Marshal(pipeline)
	if err != nil {
		return fmt.Errorf("error while marshaling pipeline spec: %s", err)
	}

	client := proto.NewAgentServiceClient(conn)
	var pStatus *proto.PipelineStatus

	protoPipeline := &proto.Pipeline{
		Name:     name,
		Spec:     string(data),
		Workflow: wfName,
	}

	switch action {
	case ActionTypeUpdate:
		log.Infof("performing update on pipline: %s", name)
		pStatus, err = client.UpdatePipeline(context.TODO(), protoPipeline)
	case ActionTypeSchedule:
		log.Infof("performing schedule on pipeline: %s", name)
		pStatus, err = client.SchedulePipeline(context.TODO(), protoPipeline)
	case ActionTypeRemove:
		log.Infof("performing delete on pipeline: %s", name)
		pStatus, err = client.RemovePipeline(context.TODO(), protoPipeline)
	}

	if err != nil {
		return fmt.Errorf("error while updating from agent: %s", err)
	}

	protoutils.UpdatePipelineStatusSpecFromGRPCTransport(status, pStatus)
	return nil
}

// UpdatePipeline updates the pipeline manifest scheduled on some agent.
func (s *Scheduler) UpdatePipeline(wfName, name string, new *v1alpha1.PipelineSpec, status *v1alpha1.WorkflowStatus) error {
	plStatus, ok := status.Pipelines[name]
	if !ok {
		return fmt.Errorf("no pipeline with name %s found in workflow status", name)
	}

	return s.performPipelineAction(ActionTypeUpdate, wfName, name, new, &plStatus)
}

// RemovePipeline removes the scheduled pipeline from an agent.
func (s *Scheduler) RemovePipeline(wfName, name string, pipeline *v1alpha1.PipelineSpec, status *v1alpha1.WorkflowStatus) error {
	plStatus, ok := status.Pipelines[name]
	if !ok {
		return fmt.Errorf("no pipeline with name %s found in workflow status", name)
	}

	return s.performPipelineAction(ActionTypeRemove, wfName, name, pipeline, &plStatus)
}

// SchedulePipeline finds out an appropriate agent for the pipeline
// and schedules the run.
func (s *Scheduler) SchedulePipeline(wfName, name string, pipeline *v1alpha1.PipelineSpec, status *v1alpha1.WorkflowStatus) error {
	ps := v1alpha1.PipelineStatus{}
	status.Pipelines[name] = ps

	return s.performPipelineAction(ActionTypeSchedule, wfName, name, pipeline, &ps)
}
