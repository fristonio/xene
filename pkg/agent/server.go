package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fristonio/xene/pkg/defaults"
	xerr "github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"
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
	pName := v1alpha1.GetWorkflowPrefixedName(pipeline.Workflow, pipeline.Name)
	log.Debugf("rpc to schedule pipeline on the agent: %s", pName)

	// During scheduling if the pipeline already exists in the datastore
	// we consider this as an error.
	exist, err := store.KVStore.Exists(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.PipelineKeyPrefix, pName))
	if err != nil {
		return nil, fmt.Errorf("error while checking pipeline in datastore: %s", err)
	}
	if exist {
		return nil, fmt.Errorf("The pipeline %s already exist in agent, cannot schedule again", pName)
	}

	return a.updatePipeline(ctx, pipeline)
}

// UpdatePipeline is the RPC to update a pipeline on to the agent
func (a *agentServer) UpdatePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	pName := v1alpha1.GetWorkflowPrefixedName(pipeline.Workflow, pipeline.Name)
	log.Debugf("rpc to update pipeline on the agent: %s", pName)
	return a.updatePipeline(ctx, pipeline)
}

// RemovePipeline is the RPC to remove a pipeline from the agent.
func (a *agentServer) RemovePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	pName := v1alpha1.GetWorkflowPrefixedName(pipeline.Workflow, pipeline.Name)
	tName := v1alpha1.GetWorkflowPrefixedName(pipeline.Workflow, pipeline.TriggerName)
	log.Debugf("rpc to remove pipeline from the agent: %s", pName)

	var (
		trigger v1alpha1.TriggerSpecWithName
	)
	val, err := store.KVStore.Get(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.TriggerKeyPrefix, tName))
	if err != nil && store.KVStore.KeyDoesNotExistError(err) {
		log.Infof("no trigger exists associated with the pipeline")
	} else if err != nil {
		return nil, fmt.Errorf("error while checking trigger in datastore: %s", err)
	} else {
		err := json.Unmarshal(val.Data, &trigger)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling trigger(%s) spec from datastore: %s", tName, err)
		}

		trigger.RemovePipeline(pName)
		if len(trigger.Pipelines) == 0 {
			err := store.KVStore.Delete(context.TODO(),
				fmt.Sprintf("%s/%s", v1alpha1.TriggerKeyPrefix, tName))
			if err != nil {
				return nil, fmt.Errorf("error while removing trigger: %s", err)
			}
		} else {
			log.Infof("more pipelines are associated with the trigger(%s), updating spec", tName)
			tData, err := json.Marshal(&trigger)
			if err != nil {
				return nil, fmt.Errorf("error while marshaling trigger(%s) spec: %s", tName, err)
			}
			err = store.KVStore.Set(context.TODO(),
				fmt.Sprintf("%s/%s", v1alpha1.TriggerKeyPrefix, tName),
				tData)
			if err != nil {
				return nil, fmt.Errorf("error while setting trigger(%s): %s", tName, err)
			}
		}
	}

	err = store.KVStore.Delete(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.PipelineKeyPrefix, pName))
	if err != nil {
		// Error at this point means that we don't have any trigger to run
		// the pipeline.
		return &proto.PipelineStatus{
			Status:   "NoTrigger",
			Executor: a.name,
		}, fmt.Errorf("error while removing pipeline(%s): %s", pName, err)
	}

	log.Infof("pipeline has been removed from the agent")
	return &proto.PipelineStatus{
		Status:   "NotScheduled",
		Executor: a.name,
	}, nil
}

// updatePipeline updates the provided pipeline spec.
func (a *agentServer) updatePipeline(ctx context.Context, pipeline *proto.Pipeline) (*proto.PipelineStatus, error) {
	pName := v1alpha1.GetWorkflowPrefixedName(pipeline.Workflow, pipeline.Name)
	tName := v1alpha1.GetWorkflowPrefixedName(pipeline.Workflow, pipeline.TriggerName)

	var (
		pSpec v1alpha1.PipelineSpec
		tSpec v1alpha1.TriggerSpec
	)
	err := json.Unmarshal([]byte(pipeline.Spec), &pSpec)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling pipeline spec: %s", err)
	}

	err = json.Unmarshal([]byte(pipeline.TriggerSpec), &tSpec)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling trigger spec: %s", err)
	}

	var (
		trigger v1alpha1.TriggerSpecWithName
	)
	val, err := store.KVStore.Get(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.TriggerKeyPrefix, tName))
	if err != nil && store.KVStore.KeyDoesNotExistError(err) {
		log.Infof("creating trigger specification in the datastore: %s", tName)
		trigger = v1alpha1.TriggerSpecWithName{
			Name:        pipeline.TriggerName,
			Workflow:    pipeline.Workflow,
			Pipelines:   []string{pName},
			TriggerSpec: tSpec,
		}
	} else if err != nil {
		return nil, fmt.Errorf("error while checking trigger in datastore: %s", err)
	} else {
		log.Infof("updating trigger specification in the datastore: %s", tName)
		err := json.Unmarshal(val.Data, &trigger)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling trigger(%s) spec from datastore: %s", tName, err)
		}

		trigger.TriggerSpec = tSpec
		trigger.AddPipeline(pName)
	}

	tData, err := json.Marshal(&trigger)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling trigger(%s) spec: %s", tName, err)
	}
	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.TriggerKeyPrefix, tName),
		tData)
	if err != nil {
		return nil, fmt.Errorf("error while setting trigger(%s): %s", tName, err)
	}

	pData, err := json.Marshal(&v1alpha1.PipelineSpecWithName{
		Name:         pipeline.Name,
		Workflow:     pipeline.Workflow,
		PipelineSpec: pSpec,
	})
	if err != nil {
		return nil, fmt.Errorf("error while marshaling pipeline(%s) spec: %s", pName, err)
	}
	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.PipelineKeyPrefix, pName),
		pData)
	if err != nil {
		return nil, fmt.Errorf("error while setting trigger(%s): %s", tName, err)
	}

	log.Infof("pipeline has been scheduled/updated on the agent.")
	return &proto.PipelineStatus{
		Status:   "Scheduled",
		Executor: a.name,
	}, nil
}

// Info is the RPC to return the info about the agent.
func (a *agentServer) Info(ctx context.Context, opts *proto.AgentInfoOpts) (*proto.AgentInfo, error) {
	var info = proto.AgentInfo{
		Healthy:    true,
		Name:       option.Config.Agent.Name,
		Secure:     !option.Config.Agent.Insecure,
		Address:    option.Config.Agent.Address,
		ServerName: option.Config.Agent.ServerName,
		Workflows:  []*proto.AgentWorkflowInfo{},
		Secrets:    getAgentSecretsList(option.Config.Agent.Name),
	}

	var wfInfo = make(map[string]*proto.AgentWorkflowInfo)
	store.KVStore.PrefixScanWithFunction(ctx, v1alpha1.TriggerKeyPrefix, func(kv *v1alpha1.KVPairStruct) {
		tName := strings.TrimPrefix(kv.Key, fmt.Sprintf("%s/", v1alpha1.TriggerKeyPrefix))
		wfName := strings.Split(tName, defaults.Seperator)[0]
		trigger := strings.Split(tName, defaults.Seperator)[1]

		var spec v1alpha1.TriggerSpecWithName
		err := json.Unmarshal([]byte(kv.Value), &spec)
		if err != nil {
			return
		}

		pipelines := utils.TrimWorkflowPrefix(spec.Pipelines, fmt.Sprintf("%s%s", wfName, defaults.Seperator))
		if _, ok := wfInfo[wfName]; ok {
			wfInfo[wfName].Triggers = append(wfInfo[wfName].Triggers, &proto.AgentTriggerInfo{
				Name:      trigger,
				Pipelines: pipelines,
			})
		} else {
			wfInfo[wfName] = &proto.AgentWorkflowInfo{
				Name: wfName,
				Triggers: []*proto.AgentTriggerInfo{
					{
						Name:      trigger,
						Pipelines: pipelines,
					},
				},
			}
		}
	})

	for _, i := range wfInfo {
		info.Workflows = append(info.Workflows, i)
	}

	return &info, nil
}

func (a *agentServer) GetPipelinesRunInfo(ctx context.Context, opts *proto.PipelineInfoOpts) (*proto.PipelinesRunInfo, error) {
	if opts.Name == "" || opts.Workflow == "" {
		return nil, errors.New("pipeline name is a required option")
	}

	merr := xerr.NewMultiError()

	resp := proto.PipelinesRunInfo{
		Name: opts.Name,
	}

	infos := []*proto.PipelineRunInfo{}

	store.KVStore.PrefixScanWithFunction(
		context.TODO(),
		fmt.Sprintf("%s/%s/", v1alpha1.PipelineStatusKeyPrefix, v1alpha1.GetWorkflowPrefixedName(opts.Workflow, opts.Name)),
		func(kv *v1alpha1.KVPairStruct) {
			pipeline := v1alpha1.PipelineRunStatus{}
			err := json.Unmarshal([]byte(kv.Value), &pipeline)
			if err != nil {
				merr.Append(fmt.Errorf("error unmarshalling %s", kv.Key))
				return
			}

			infos = append(infos, &proto.PipelineRunInfo{
				RunID:  pipeline.RunID,
				Agent:  pipeline.Agent,
				Status: pipeline.Status,
			})
		})

	resp.Pipelines = infos
	resp.ErrorMessage = merr.String()
	return &resp, nil
}
