package v1alpha1

import (
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
)

// GetWorkflowPrefixedName returns the name prefixed with the workflow
// name.
func GetWorkflowPrefixedName(wf, name string) string {
	return fmt.Sprintf("%s%s%s", wf, defaults.Seperator, name)
}
