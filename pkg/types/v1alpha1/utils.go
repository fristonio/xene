package v1alpha1

import (
	"fmt"
	"regexp"
	"time"

	"github.com/fristonio/xene/pkg/defaults"
)

// GetWorkflowPrefixedName returns the name prefixed with the workflow
// name.
func GetWorkflowPrefixedName(wf, name string) string {
	return fmt.Sprintf("%s%s%s", wf, defaults.Seperator, name)
}

// GetDummyPipelineRunStatus returns a dummy pipeline run status
func GetDummyPipelineRunStatus(p *PipelineSpecWithName) PipelineRunStatus {
	status := PipelineRunStatus{
		Name:      p.Name,
		RunID:     "",
		Status:    StatusNotExecuted,
		Agent:     "",
		StartTime: time.Now().Unix(),
		EndTime:   0,
		Tasks:     make(map[string]*TaskRunStatus),
	}

	for tName, task := range p.Tasks {
		tStatus := TaskRunStatus{
			Status:       StatusNotExecuted,
			Dependencies: task.Dependencies,
			Steps:        make(map[string]*StepRunStatus),
		}

		for _, step := range task.Steps {
			sStatus := &StepRunStatus{
				Status: StatusNotExecuted,
			}

			tStatus.Steps[step.Name] = sStatus
		}

		status.Tasks[tName] = &tStatus
	}

	return status
}

// NameRegex is the regex for any ID/Name used by xene.
// Matches any DNS Subdomain Names
var NameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+[a-zA-Z0-9-._]*[a-zA-Z0-9]+$`)

// IsValidDNSSubdomainName checks if the provided string is a valid DNS subdomain name
func IsValidDNSSubdomainName(name string) bool {
	return NameRegex.Match([]byte(name)) && len(name) < 254
}
