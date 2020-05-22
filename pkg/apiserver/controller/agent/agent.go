package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fristonio/xene/pkg/controller"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

// Controller corresponds to the agent controller managed by apiserver.
type Controller struct {
	// Manager is the controller manager for all the running controller
	// against agent manager.
	Manager *controller.Manager

	// Nodes contains the agent associated with the controller.
	// Each of these agents are managed by the agent controller using a global
	// controller manager
	Nodes map[string]Node

	// blacklistedNodes contains a list of blacklisted agents which we were not reachable
	// in the past. We will try to check for the agent health again after a cooldown.
	blacklistedNodes map[string]Node

	// storeInformer is the store informer corresponding to the agent objects
	// in the datastore.
	storeInformer *store.Informer

	// name is the name of the controller, this is derived from the store controller.
	name string

	// agentHealthCheckLimit is the number of times we are going to try to check agent's
	// health before we say that agent is not recoverable.
	agentHealthCheckLimit int64

	// healthCheckInterval is the interval at which to check the health of each of the agent.
	healthCheckInterval time.Duration
}

// Node represents a single Agent object
type Node struct {
	Agent *v1alpha1.Agent

	Name string

	conn *grpc.ClientConn
}

// CheckHealth checks for the health status of the provided agent in context.
// If the agent is not healthy or there is some issue with the connectivity the
// function returns an error.
func (n *Node) CheckHealth() error {
	client := proto.NewAgentServiceClient(n.conn)
	status, err := client.Status(context.TODO(), &proto.StatusOpts{Verbose: true})
	if err != nil {
		return fmt.Errorf("error while fetching status: %s", err)
	}

	if !status.Healthy {
		return fmt.Errorf("agent is not healthy")
	}

	return nil
}

// AgentCtrl is the global instance of the agent controller running for the API server.
var AgentCtrl *Controller = NewController()

// NewController returns a new agent controller manager which manages all the
// agents running for the API server.
func NewController() *Controller {
	return &Controller{
		Manager:               controller.NewManager(),
		Nodes:                 make(map[string]Node),
		blacklistedNodes:      make(map[string]Node),
		agentHealthCheckLimit: defaults.AgentHealthCheckRetriesLimit,
		healthCheckInterval:   defaults.AgentHealthCheckInterval,
	}
}

// Type returns the type of controller configured, for agent this is agent.
func (a *Controller) Type() string {
	return "agent"
}

// Configure sets up the Agent controller and all its required components.
func (a *Controller) Configure() {
	a.storeInformer = a.newAgentStoreController()
	a.name = a.storeInformer.Name()

	ctrlFunc, err := controller.NewControllerFunction(a.blacklistedNodesController)
	if err != nil {
		log.Warnf("error while creating controller for balacklisted nodes: %s", err)
	}
	err = a.Manager.UpdateController("agent-blacklisted-nodes", "agent-controller", controller.Internal{
		DoFunc:      ctrlFunc,
		RunInterval: a.healthCheckInterval * 10,
	})
	if err != nil {
		log.Warnf("error while updating blacklisted nodes controller: %s", err)
	}
}

// Run starts running the agent controller.
func (a *Controller) Run() error {
	return a.storeInformer.Run()
}

// Stop shuts down the controller.
func (a *Controller) Stop() error {
	a.Manager.RemoveAllAndWait()
	for _, node := range a.Nodes {
		if node.conn != nil {
			node.conn.Close()
		}
	}

	a.Nodes = make(map[string]Node)
	a.blacklistedNodes = make(map[string]Node)
	return a.storeInformer.Stop()
}

// Name returns the name of the agent controller, it is completely defined by the name of
// the underlying store controller.
func (a *Controller) Name() string {
	return a.name
}

// GetAllActiveAgents returns all the agent which are not blacklisted.
func (a *Controller) GetAllActiveAgents() []string {
	agents := make([]string, 0)
	for name := range a.Nodes {
		agents = append(agents, name)
	}

	return agents
}

// GetAllAgentsWithInfo returns a list with info of all the agents
// including blacklisted.
func (a *Controller) GetAllAgentsWithInfo() map[string]*v1alpha1.Agent {
	agents := make(map[string]*v1alpha1.Agent)
	for name, node := range a.Nodes {
		agents[name] = node.Agent
	}

	for name, node := range a.blacklistedNodes {
		agents[name] = node.Agent
	}

	return agents
}

// GetActiveAgentsWithInfo returns a list with info of all active agents
// not including blacklisted.
func (a *Controller) GetActiveAgentsWithInfo() map[string]*v1alpha1.Agent {
	agents := make(map[string]*v1alpha1.Agent)
	for name, node := range a.Nodes {
		agents[name] = node.Agent
	}

	return agents
}

// GetBlacklistedAgentsWithInfo returns a list with info of all the blacklisted agents
func (a *Controller) GetBlacklistedAgentsWithInfo() map[string]*v1alpha1.Agent {
	agents := make(map[string]*v1alpha1.Agent)
	for name, node := range a.blacklistedNodes {
		agents[name] = node.Agent
	}

	return agents
}

// AgentConnection returns the grpc connection for the agent.
func (a *Controller) AgentConnection(agent string) *grpc.ClientConn {
	ag, ok := a.Nodes[agent]
	if !ok {
		return nil
	}

	if ag.conn == nil {
		conn, err := GetAgentConnection(ag.Agent)
		if err != nil {
			log.Errorf("error while establishing agent connection, blacklisting agent: %s", err)
			err = a.Manager.RemoveController(agent)
			if err != nil {
				log.Errorf("error while removing controller for %s: %s", agent, err)
				return nil
			}
			delete(a.Nodes, agent)
			a.blacklistedNodes[agent] = Node{
				Name:  agent,
				Agent: ag.Agent,
			}
		}

		ag.conn = conn
	}

	return ag.conn
}

// BlacklistedAgents returns a list of agent nodes which are blacklisted by the agent
// controller
func (a *Controller) BlacklistedAgents() []string {
	var agents []string
	for name := range a.blacklistedNodes {
		agents = append(agents, name)
	}

	return agents
}

func (a *Controller) addController(kv *v1alpha1.KVPairStruct) error {
	var ag v1alpha1.Agent
	err := json.Unmarshal([]byte(kv.Value), &ag)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the agent spec from data: %s", err)
	}

	conn, err := GetAgentConnection(&ag)
	if err != nil {
		// if the error occurs while creating the connection then we delete
		// the event from the cache and will try to set this up later.
		return fmt.Errorf("error while creating agent connection: %s", err)
	}

	ctrlFunc, err := controller.NewControllerFunction(a.ControllerFunc, &ag)
	if err != nil {
		return fmt.Errorf("error while creating controller function: %s", err)
	}
	err = a.Manager.UpdateController(ag.Metadata.GetName(), "agent-controller", controller.Internal{
		DoFunc:      ctrlFunc,
		RunInterval: a.healthCheckInterval,
	})
	if err != nil {
		return fmt.Errorf("error while updating agent controller for : %s", ag.Metadata.GetName())
	}
	a.Nodes[ag.Metadata.GetName()] = Node{
		Agent: &ag,
		Name:  ag.Metadata.GetName(),
		conn:  conn,
	}
	delete(a.blacklistedNodes, ag.Metadata.GetName())
	return nil
}

// newAgentStoreController returns the agent store controller for apiserver.
// This controller watches for agent object in the store and perform action based
// on the changes to the object.
func (a *Controller) newAgentStoreController() *store.Informer {
	return store.NewInformerWithSharedCache(
		fmt.Sprintf("%s/", v1alpha1.AgentKeyPrefix),
		// Add function for a new agent.
		func(kv *v1alpha1.KVPairStruct) error {
			log.Infof("agent added: %s", kv.Key)
			return a.addController(kv)
		},

		// Delete function for the new agent.
		func(key string) error {
			agentName := strings.TrimPrefix(key, v1alpha1.AgentKeyPrefix+"/")
			log.Infof("agent deleted: %s", agentName)

			err := a.Manager.RemoveController(agentName)
			if err != nil {
				return fmt.Errorf("error while removing controller for %s: %s", agentName, err)
			}

			delete(a.Nodes, agentName)
			delete(a.blacklistedNodes, agentName)
			return nil
		},

		// Update function for a new agent function.
		func(kv *v1alpha1.KVPairStruct, version uint64) error {
			log.Infof("Agent updated: %s", kv.Key)
			return a.addController(kv)
		},
	)
}

// ControllerFunc is the function which runs for each agent connected to the Api
// server and deals with it's management.
func (a *Controller) ControllerFunc(_ctx context.Context, ag *v1alpha1.Agent) error {
	log.Debugf("running agent controller function for agent: %s", ag.Metadata.GetName())

	// Here we simply check for the health status of the agent and
	// and remove the agent from the controller list if the health check fails
	// according to a certail criteria.
	status := a.Manager.Lookup(ag.Metadata.GetName()).Status().Status
	node, ok := a.Nodes[ag.Metadata.GetName()]
	if !ok {
		log.Errorf("error node not registered: %s", ag.Metadata.GetName())
	}

	if status.ConsecutiveFailureCount > a.agentHealthCheckLimit {
		err := a.Manager.RemoveController(ag.Metadata.GetName())
		if err != nil {
			return fmt.Errorf("error while removing agent controller: %s", err)
		}
		node.conn.Close()
		delete(a.Nodes, ag.Metadata.GetName())
		a.blacklistedNodes[ag.Metadata.GetName()] = Node{
			Name:  ag.Metadata.GetName(),
			Agent: ag,
		}

		return fmt.Errorf("agent health check failed too many time, keeping agent in blacklist")
	}

	err := node.CheckHealth()
	return err
}

func (a *Controller) blacklistedNodesController(_ctx context.Context) error {
	for name, node := range a.blacklistedNodes {
		conn, err := GetAgentConnection(node.Agent)
		if err != nil {
			continue
		}
		err = node.CheckHealth()
		if err == nil {
			log.Infof("node %s back online", name)
			node.conn = conn
			ctrlFunc, err := controller.NewControllerFunction(a.ControllerFunc, node.Agent)
			if err != nil {
				log.Errorf("error while creating controller function: %s", err)
				continue
			}
			err = a.Manager.UpdateController(name, "agent-controller", controller.Internal{
				DoFunc:      ctrlFunc,
				RunInterval: a.healthCheckInterval,
			})
			if err != nil {
				log.Errorf("error while updating agent controller for : %s", name)
				continue
			}
			a.Nodes[name] = node
			delete(a.blacklistedNodes, name)
		}
	}

	return nil
}
