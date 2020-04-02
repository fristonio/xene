package executor

// Executor is the standard interface which must be implemented
// by any plugin acting as an executor for xene.
type Executor interface {
	// Configure configures the executor for any task associated with running
	// the workflow on xene.
	Configure() error
}
