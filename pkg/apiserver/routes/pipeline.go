package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/controller/agent"
	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/gin-gonic/gin"
)

// @Summary Returns verbose information about a workflow.
// @Tags info
// @Accept  json
// @Produce json
// @Param name path string true "Name of the workflow to get information about."
// @Param pipeline path string true "Name of the pipeline to return the info about."
// @Success 200 {object} response.PipelineInfo
// @Failure 400 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/info/workflow/{name}/pipeline/{pipeline} [get]
func pipelineInfoHandler(ctx *gin.Context) {
	wfName := ctx.Param("name")
	pipeline := ctx.Param("pipeline")
	if wfName == "" || pipeline == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "name and pipeline are required parameters",
		})
		return
	}

	val, err := store.KVStore.Get(context.TODO(), fmt.Sprintf("%s/%s", v1alpha1.WorkflowKeyPrefix, wfName))
	if err != nil {
		if store.KVStore.KeyDoesNotExistError(err) {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: fmt.Sprintf("requested workflow: %s does not exist", wfName),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while getting the workflow: %s: %s", wfName, err),
		})
		return
	}

	wf := v1alpha1.Workflow{}
	err = json.Unmarshal(val.Data, &wf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshalling the workflow: %s: %s", wfName, err),
		})
		return
	}

	if _, ok := wf.Spec.Pipelines[pipeline]; !ok {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("the pipeline %s does not exist in workflow manifest", pipeline),
		})
		return
	}

	spec, err := json.Marshal(wf.Spec.Pipelines[pipeline])
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while marshalling the pipeline: %s: %s", pipeline, err),
		})
		return
	}
	pipelineInfo := response.PipelineInfo{
		Workflow: wfName,
		Name:     pipeline,
		Spec:     string(spec),
	}

	v, err := store.KVStore.Get(context.TODO(), fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wfName))
	if err != nil {
		if store.KVStore.KeyDoesNotExistError(err) {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: fmt.Sprintf("requested workflow status: %s does not exist", wfName),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while getting the workflow status: %s: %s", wfName, err),
		})
		return
	}

	wfStatus := v1alpha1.WorkflowStatus{}
	err = json.Unmarshal(v.Data, &wfStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshalling the workflow status: %s: %s", wfName, err),
		})
		return
	}

	pStatus, ok := wfStatus.Pipelines[pipeline]
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while getting pipeline from workflow status"),
		})
		return
	}

	agents := pStatus.GetAllExecutors()
	jobs := make(chan string, len(agents))
	results := make(chan *proto.PipelinesRunInfo, len(agents))
	infos := []response.PipelineRunInfo{}

	for w := 1; w <= 4; w++ {
		go getPipelineInfoFromAgent(w, jobs, results, pipeline, wfName)
	}

	for j := 0; j < len(agents); j++ {
		jobs <- agents[j]
	}
	close(jobs)

	// Collect results
	for a := 1; a <= len(agents); a++ {
		info := <-results
		for _, i := range info.Pipelines {
			infos = append(infos, response.PipelineRunInfo{
				RunID:     i.RunID,
				Agent:     i.Agent,
				Status:    i.Status,
				StartTime: i.StartTime,
				EndTime:   i.EndTime,
			})
		}
	}
	pipelineInfo.Runs = infos
	ctx.JSON(http.StatusOK, pipelineInfo)
}

func getPipelineInfoFromAgent(id int, jobs <-chan string, results chan<- *proto.PipelinesRunInfo, pName, wfName string) {
	for ag := range jobs {
		conn := agent.AgentCtrl.AgentConnection(ag)
		if conn == nil {
			results <- &proto.PipelinesRunInfo{
				Name:         pName,
				ErrorMessage: fmt.Sprintf("cannot retrive agent connection for: %s", ag),
			}
			continue
		}

		client := proto.NewAgentServiceClient(conn)
		info, err := client.GetPipelinesRunInfo(context.TODO(), &proto.PipelineInfoOpts{
			Name:     pName,
			Workflow: wfName,
		})
		if err != nil {
			results <- &proto.PipelinesRunInfo{
				Name:         pName,
				ErrorMessage: fmt.Sprintf("error while getting pipelines info: %s", err),
			}
			continue
		}

		results <- info
	}
}
