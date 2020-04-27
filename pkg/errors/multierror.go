package errors

import (
	"fmt"
	"strings"
)

// MultiError is a type to handle multiples errors simultaneously.
type MultiError struct {
	// errors contains the list of the errors
	errors []Error
}

// NewMultiError returns a new instance of MultiError type.
func NewMultiError() *MultiError {
	return &MultiError{
		errors: make([]Error, 0),
	}
}

// GetError returns the multi error itself if the severity of any one of
// the message is error, else it returns nil.
func (m *MultiError) GetError() error {
	if m.HasErrors() {
		return m
	}

	return nil
}

// Append appends the error with the multierror using the default
// serverity of Error
func (m MultiError) Append(err error) {
	m.errors = append(m.errors, &errorWithSeverity{
		err,
		SeverityTypeError,
	})
}

// AppendWithSeverity appends the error with the multierror using the provided severity
func (m MultiError) AppendWithSeverity(err error, severity Severity) {
	m.errors = append(m.errors, &errorWithSeverity{
		err,
		severity,
	})
}

// String returns the string representation of multiple errors
func (m *MultiError) String() string {
	var strs []string
	for _, err := range m.errors {
		strs = append(strs, fmt.Sprintf("%s: %s", err.Severity(), err.Error()))
	}

	return strings.Join(strs, "\n")
}

// Error is implemented to be categorized as an error.
func (m *MultiError) Error() string {
	return m.String()
}

// Severity returns the sevirity of the multierror type by selecting the highest
// severity among the given.
func (m *MultiError) Severity() Severity {
	var s Severity = SeverityTypeInfo
	for _, err := range m.errors {
		if err.Severity() == SeverityTypeWarn {
			s = SeverityTypeWarn
		} else if err.Severity() == SeverityTypeError {
			return SeverityTypeError
		}
	}

	return s
}

// HasErrors returns if any of the errors has severity error
func (m *MultiError) HasErrors() bool {
	return m.Severity() == SeverityTypeError
}
