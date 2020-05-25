package agent

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fristonio/xene/pkg/dag"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
)

func getAgentSecretsList(agentName string) []*proto.AgentSecretInfo {
	return []*proto.AgentSecretInfo{
		{
			Name: fmt.Sprintf("%s-client-key", agentName),
			Type: defaults.SecretTypeCertificates,
		},
		{
			Name: fmt.Sprintf("%s-client-cert", agentName),
			Type: defaults.SecretTypeCertificates,
		},
		{
			Name: fmt.Sprintf("%s-root-ca-chain", agentName),
			Type: defaults.SecretTypeCertificates,
		},
	}
}

// GetPipelineRunLogReader returns a reader to read the logs for a particular pipeline run
func GetPipelineRunLogReader(
	workflow, pipeline, runID string,
	pipelineSpec *v1alpha1.PipelineSpec, pipelineStatus *v1alpha1.PipelineRunStatus) (io.Reader, int64, error) {
	logFiles := []struct {
		path, task, step string
	}{}
	// Walk each of task in the pipeline in the required order.
	merrs := pipelineSpec.Dag.Walk(func(v dag.Vertex) *errors.MultiError {
		errs := errors.NewMultiError()
		task, ok := v.(*v1alpha1.TaskSpec)
		if !ok {
			// If any error is getting the task spec then append it to the list of errors
			errs.Append(fmt.Errorf("not a valid vertex to visit, must confirm to type *TaskSpec"))
			return errs
		}

		taskStatus := pipelineStatus.Tasks[task.Name()]
		if taskStatus.Steps == nil {
			return errs
		}
		for stepName, step := range taskStatus.Steps {
			if step.LogFile != "" {
				logFiles = append(logFiles, struct {
					path, task, step string
				}{
					path: step.LogFile,
					task: task.Name(),
					step: stepName,
				})
			}
		}

		return errs
	})

	if len(merrs) > 0 {
		return nil, 0, fmt.Errorf("%v", merrs)
	}

	var readers []io.Reader
	var contentLength int64

	for _, logFile := range logFiles {
		cr := strings.NewReader(fmt.Sprintf(`
***************************************************
Workflow: %s
Pipeline: %s
RunID: %s
Task: %s
Step: %s
***************************************************
`, workflow, pipeline, runID, logFile.task, logFile.step))
		path := filepath.Join(defaults.AgentLogsDir, logFile.path)
		f, err := os.Open(path)
		if err != nil {
			r := strings.NewReader(fmt.Sprintf("\nError while opening logfile path: %s\n", err))
			readers = append(readers, r)
			contentLength += r.Size()
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			r := strings.NewReader(fmt.Sprintf("\nError while stat on logfile path: %s\n", err))
			readers = append(readers, r)
			contentLength += r.Size()
			continue
		}

		readers = append(readers, cr)
		readers = append(readers, f)
		contentLength += cr.Size() + info.Size()
	}

	return io.MultiReader(readers...), contentLength, nil
}
