package errors

// Severity type describes the sevierty of any message.
type Severity string

var (
	// SeverityTypeWarn is for warnings
	SeverityTypeWarn Severity = "warn"

	// SeverityTypeError is for errors
	SeverityTypeError Severity = "error"

	// SeverityTypeInfo is for information
	SeverityTypeInfo Severity = "info"
)
