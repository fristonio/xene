package defaults

import (
	"time"
)

var (
	// ControllerType is the default controller type we are assigning to the created
	// controller.
	ControllerType string = "Default"

	// ControllerRetryInterval is the interval for controller function execution retry.
	ControllerRetryInterval time.Duration = 1 * time.Second

	// ControllerInvalidDuration is the placeholder time duration for the execution
	// of the function when there is some error in exeuction.
	ControllerInvalidDuration time.Duration = 1000 * time.Second

	// AgentHealthCheckRetriesLimit is the limit on consecutive error while checking healths
	// of the agent. If more then limit number of consecutive error occur during the health
	// check, the agent is blacklisted.
	AgentHealthCheckRetriesLimit int64 = 3

	// AgentHealthCheckInterval is the default value of time interval at which to execute
	// the health checks for configured agents.
	AgentHealthCheckInterval time.Duration = time.Second * 5
)
