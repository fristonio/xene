package executor

import "github.com/fristonio/xene/pkg/types/v1alpha1"

// RuntimeExecutor is the interface implemented by each of the runtime executor
// provided for xene.
type RuntimeExecutor interface {
	// Configure sets up the environment for the RuntimeExecutor
	// to execute the tasks.
	Configure() error

	// RunTask runs the task provided in the specification
	RunTask(string, *v1alpha1.TaskSpec) error

	// Shutdown shuts down the RuntimeExecutor instance and environment.
	Shutdown() error
}
