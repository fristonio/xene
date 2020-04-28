package v1alpha1

import (
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
)

// GetWorkflowAppendedName returns the name prefixed with the workflow
// name.
func GetWorkflowAppendedName(wf, name string) string {
	return fmt.Sprintf("%s%s%s", wf, defaults.Seperator, name)
}
