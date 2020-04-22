package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
)

// AgentController corresponds to the agent controller managed by apiserver.
type AgentController struct {
	// Manager is the controller manager for all the running controller
	// against agent manager.
	Manager *controller.Manager

	// Nodes contains the agent associated with the controller.
	Nodes map[string]*v1alpha1.Agent

	// blacklistedNodes contains a list of blacklisted agents which we were not reachable
	// in the past. We will try to check for the agent health again after a cooldown.
	blacklistedNodes map[string]*v1alpha1.Agent

	// storeCtrl is the store controller corresponding to the agent objects
	// in the datastore.
	storeCtrl *store.Controller

	// name is the name of the controller, this is derived from the store controller.
	name string

	// agentHealthCheckLimit is the number of times we are going to try to check agent's
	// health before we say that agent is not recoverable.
	agentHealthCheckLimit int64

	// healthCheckInterval is the interval at which to check the health of each of the agent.
	healthCheckInterval time.Duration
}

// AgentCtrl is the global instance of the agent controller running for the API server.
var AgentCtrl *AgentController

// NewAgentController returns a new agent controller manager which manages all the
// agents running for the API server.
func NewAgentController() *AgentController {
	return &AgentController{
		Manager:               controller.NewManager(),
		Nodes:                 make(map[string]*v1alpha1.Agent),
		blacklistedNodes:      make(map[string]*v1alpha1.Agent),
		agentHealthCheckLimit: 3,
		healthCheckInterval:   time.Second * 5,
	}
}

// Type returns the type of controller configured, for agent this is agent.
func (a *AgentController) Type() string {
	return "agent"
}

// Configure sets up the Agent controller and all its required components.
func (a *AgentController) Configure() {
	a.storeCtrl = a.newAgentStoreController()
	a.name = a.storeCtrl.Name()
}

// Run starts running the agent controller.
func (a *AgentController) Run() error {
	return a.storeCtrl.Run()
}

// Stop shuts down the controller.
func (a *AgentController) Stop() error {
	return a.storeCtrl.Stop()
}

// Name returns the name of the agent controller, it is completely defined by the name of
// the underlying store controller.
func (a *AgentController) Name() string {
	return a.name
}

// BlacklistedAgents returns a list of agent nodes which are blacklisted by the agent
// controller
func (a *AgentController) BlacklistedAgents() []string {
	var agents []string
	for name := range a.blacklistedNodes {
		agents = append(agents, name)
	}

	return agents
}

func (a *AgentController) addAgentController(kv *v1alpha1.KVPairStruct) {
	var ag v1alpha1.Agent
	err := json.Unmarshal([]byte(kv.Value), &ag)
	if err != nil {
		log.Errorf("error while unmarshaling the agent spec from data: %s", err)
		return
	}

	ctrlFunc, err := controller.NewControllerFunction(a.agentControllerFunc, &ag)
	if err != nil {
		log.Errorf("error while creating controller function: %s", err)
		return
	}
	err = a.Manager.UpdateController(ag.Metadata.GetName(), "agent-controller", controller.Internal{
		DoFunc:      ctrlFunc,
		RunInterval: a.healthCheckInterval,
	})
	if err != nil {
		log.Errorf("error while updating agent controller for : %s", ag.Metadata.GetName())
		return
	}
	a.Nodes[ag.Metadata.GetName()] = &ag
	delete(a.blacklistedNodes, ag.Metadata.GetName())
}

// newAgentStoreController returns the agent store controller for apiserver.
// This controller watches for agent object in the store and perform action based
// on the changes to the object.
func (a *AgentController) newAgentStoreController() *store.Controller {
	return store.NewControllerWithSharedCache(
		fmt.Sprintf("%s/", v1alpha1.AgentKeyPrefix),
		true,
		// Add function for a new agent.
		func(kv *v1alpha1.KVPairStruct) {
			log.Infof("Agent added: %s", kv.Key)
			a.addAgentController(kv)
		},

		// Delete function for the new agent.
		func(kv *v1alpha1.KVPairStruct) {
			log.Infof("Agent deleted: %s", kv.Key)
			var ag v1alpha1.Agent
			err := json.Unmarshal([]byte(kv.Value), &ag)
			if err != nil {
				log.Errorf("error while unmarshaling the agent spec from data: %s", err)
				return
			}

			err = a.Manager.RemoveController(ag.Metadata.GetName())
			if err != nil {
				log.Errorf("error while removing controller for %s: %s", ag.Metadata.GetName(), err)
			}

			delete(a.Nodes, ag.Metadata.GetName())
			delete(a.blacklistedNodes, ag.Metadata.GetName())
		},

		// Update function for a new agent function.
		func(kv *v1alpha1.KVPairStruct, version uint64) {
			log.Infof("Agent updated: %s", kv.Key)
			a.addAgentController(kv)
		},
	)
}

// agentControllerFunc is the function which runs for each agent connected to the Api
// server and deals with it's management.
func (a *AgentController) agentControllerFunc(_ctx context.Context, ag *v1alpha1.Agent) error {
	log.Debugf("running agent controller function for agent: %s", ag.Metadata.GetName())

	// Here we simply check for the health status of the agent and
	// and remove the agent from the controller list if the health check fails
	// according to a certail criteria.
	status := a.Manager.Lookup(ag.Metadata.GetName()).Status().Status

	if status.ConsecutiveFailureCount > a.agentHealthCheckLimit {
		err := a.Manager.RemoveController(ag.Metadata.GetName())
		if err != nil {
			return fmt.Errorf("error while removing agent controller: %s", err)
		}
		delete(a.Nodes, ag.Metadata.GetName())
		a.blacklistedNodes[ag.Metadata.GetName()] = ag

		return fmt.Errorf("agent health check failed too many time, keeping agent in blacklist")
	}

	err := ag.CheckHealth()
	return err
}

func init() {
	agCtrl := NewAgentController()
	RegisteredControllers = append(RegisteredControllers, agCtrl)

	// Register the global agent controller.
	AgentCtrl = agCtrl
}
