package v1alpha1

import (
	"fmt"
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
