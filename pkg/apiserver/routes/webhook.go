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

// WebhookGroupRoutes contains all the routes related to xene info.
func WebhookGroupRoutes(r *gin.RouterGroup) {
	r.GET("/trigger/:workflow/:trigger/:pipeline", webhookTriggerHandler)
}

// @Summary Webhook endpoints to trigger the pipelines of a particular workflow.
// @Tags webhook
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param workflow path string true "Name of the workflow."
// @Param pipeline path string true "Name of the pipeline to be triggered."
// @Param trigger path string true "Name of the trigger associated with the pipeline"
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Failure 200 {object} response.HTTPMessage
// @Router /api/v1/webhook/trigger/{workflow}/{trigger}/{pipeline} [get]
func webhookTriggerHandler(ctx *gin.Context) {
	// TODO: Make this independent of Trigger. Use only workflow and pipeline
	// For this add a new field TriggerName to PipelineStatus on the apiserver.
	wfName := ctx.Param("workflow")
	pipeline := ctx.Param("pipeline")
	trig := ctx.Param("trigger")

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

	conn := agent.AgentCtrl.AgentConnection(pStatus.Executor)
	if conn == nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while getting agent connection"),
		})
		return
	}

	client := proto.NewAgentServiceClient(conn)
	info, err := client.InvokeTrigger(context.TODO(), &proto.TriggerOpts{
		Name:     trig,
		Workflow: wfName,
		Pipeline: pipeline,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error invoking pipeline associated with trigger: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: info.Status,
	})
}
