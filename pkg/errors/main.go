package errors

// Error is a custom error interface used which also includes
// the sevirity of the error.
type Error interface {
	Error() string

	Severity() Severity
}

type errorWithSeverity struct {
	err      error
	severity Severity
}

func (e *errorWithSeverity) Error() string {
	return e.err.Error()
}

func (e *errorWithSeverity) Severity() Severity {
	return e.severity
}
