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

	agentsLoadMap map[string]uint64
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
		agentsLoadMap:   make(map[string]uint64),
	}
}

func (s *Scheduler) assignNewAgent(status *v1alpha1.PipelineStatus) error {
	agents := s.AgentController.GetAllActiveAgents()
	for _, agent := range agents {
		if _, ok := s.agentsLoadMap[agent]; !ok {
			s.agentsLoadMap[agent] = 0
		}
	}

	if len(s.agentsLoadMap) == 0 {
		return fmt.Errorf("no agent configured, skipping pipeline schedule")
	}
	var (
		minWeight uint64 = 1<<64 - 1
		agent     string
	)
	for name, weight := range s.agentsLoadMap {
		if weight < minWeight {
			minWeight = weight
			agent = name
		}
	}

	prevExec := []string{}
	if status.Executor != "" {
		prevExec = append(prevExec, status.Executor)
	}
	status.Executor = agent
	for _, exec := range status.PreviousExecutors {
		if exec != agent {
			prevExec = append(prevExec, exec)
		}
	}

	status.PreviousExecutors = prevExec
	return nil
}

func (s *Scheduler) performPipelineAction(action Action, wfName, name string,
	pipeline *v1alpha1.PipelineSpec, status *v1alpha1.PipelineStatus) error {

	if status.Executor == "" {
		err := s.assignNewAgent(status)
		if err != nil {
			return err
		}
	}

	agentName := status.Executor
	conn := s.AgentController.AgentConnection(agentName)
	if conn == nil {
		return fmt.Errorf("error while retriving agent connection")
	}

	data, err := json.Marshal(pipeline)
	if err != nil {
		return fmt.Errorf("error while marshaling pipeline spec: %s", err)
	}

	var triggerData []byte
	if pipeline.Trigger != nil {
		triggerData, err = json.Marshal(pipeline.Trigger)
		if err != nil {
			return fmt.Errorf("error while marshaling trigger spec: %s", err)
		}
	}

	client := proto.NewAgentServiceClient(conn)
	var pStatus *proto.PipelineStatus

	protoPipeline := &proto.Pipeline{
		Name:        name,
		Spec:        string(data),
		Workflow:    wfName,
		TriggerName: pipeline.TriggerName,
		TriggerSpec: string(triggerData),
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
func (s *Scheduler) UpdatePipeline(wfName, name string, new *v1alpha1.PipelineSpec,
	status *v1alpha1.WorkflowStatus) error {
	log.Debugf("in scheduler update pipeline step")
	plStatus, ok := status.Pipelines[name]
	if !ok {
		return fmt.Errorf("no pipeline with name %s found in workflow status", name)
	}

	return s.performPipelineAction(ActionTypeUpdate, wfName, name, new, plStatus)
}

// RemovePipeline removes the scheduled pipeline from an agent.
func (s *Scheduler) RemovePipeline(wfName, name string, pipeline *v1alpha1.PipelineSpec,
	status *v1alpha1.WorkflowStatus) error {
	log.Debugf("in scheduler remove pipeline step")
	plStatus, ok := status.Pipelines[name]
	if !ok {
		return fmt.Errorf("no pipeline with name %s found in workflow status", name)
	}

	return s.performPipelineAction(ActionTypeRemove, wfName, name, pipeline, plStatus)
}

// SchedulePipeline finds out an appropriate agent for the pipeline
// and schedules the run.
func (s *Scheduler) SchedulePipeline(wfName, name string, pipeline *v1alpha1.PipelineSpec,
	status *v1alpha1.WorkflowStatus) error {
	log.Debugf("in scheduler schedule pipeline step")
	ps := v1alpha1.PipelineStatus{}
	status.Pipelines[name] = &ps

	return s.performPipelineAction(ActionTypeSchedule, wfName, name, pipeline, &ps)
}
