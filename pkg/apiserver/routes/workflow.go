package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	types "github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/gin-gonic/gin"
)

// workflowGetHandler retrivies a given workflow object from the store.
// @Summary Returns the specified workflow object from the store.
// @Description If a name is provided return the corresponding workflow object, if prefix  is set to some value
// then serach based on that prefix, if only name is provided retrive the associated
// workflow object
// @Tags registry
// @Accept  json
// @Produce json
// @Param prefix query string false "Prefix based get for workflow."
// @Param name query string false "name of the workflow to get."
// @Success 200 {object} response.RegistryItemsFromPrefix
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/workflow [get]
func workflowGetHandler(ctx *gin.Context) {
	storeGetHandler(ctx, v1alpha1.WorkflowKeyPrefix)
}

// @Summary Returns the specified workflow object from the store with the name in params.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "name of the workflow to get."
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/workflow/{name} [get]
func workflowGetByNameHandler(ctx *gin.Context) {
	storeGetByNameHandler(ctx, v1alpha1.WorkflowKeyPrefix)
}

// workflowCreateHandler creates a new workflow.
// @Summary Creates a new workflow in the store.
// @Description This route creates a new workflow for xene to operate on, if the workflow already exists
// this will update it with the new workflow.
// @Tags registry
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Param workflow formData string true "Workflow manifest to be created."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/registry/workflow [post]
func workflowCreateHandler(ctx *gin.Context) {
	workflow := ctx.PostForm("workflow")
	if workflow == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "workflow is a required parameter",
		})
		return
	}

	var wf types.Workflow
	err := json.Unmarshal([]byte(workflow), &wf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshling: %s", err),
		})
		return
	}

	err = wf.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while validating workflow spec: %s", err),
		})
		return
	}

	wf.RemoveNonLinkedTriggers()
	wfData, err := json.Marshal(&wf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while marshalling wrokflow object: %s", err),
		})
		return
	}

	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.WorkflowKeyPrefix, wf.Metadata.ObjectMeta.Name),
		wfData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while setting key: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: fmt.Sprintf("%s workflow created/updated", wf.Metadata.ObjectMeta.Name),
	})
}

// @Summary Patches the specified workflow from the store.
// @Tags registry
// @Accept json
// @Produce json
// @Param name path string true "Name of the workflow to be patched."
// @Failure 400 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/workflow/{name} [patch]
func workflowPatchHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, response.HTTPError{
		Error: "route not implemented yet",
	})
}

// @Summary Deletes the specified workflow from the store.
// @Description Deletes the workflow specified by the name parameter, if the workflow is not
// present then an error is thrown.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "Name of the workflow to be deleted."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/registry/workflow/{name} [delete]
func workflowDeleteHandler(ctx *gin.Context) {
	storeDeleteHandler(ctx, v1alpha1.WorkflowKeyPrefix)
}

// workflowStatusGetHandler retrivies a given workflowStatus object from the store.
// @Summary Returns the specified workflowStatus object from the store.
// @Tags status
// @Accept  json
// @Produce json
// @Param prefix query string false "Prefix based get for workflow status documents."
// @Param name query string false "name of the workflow to get status object for."
// @Success 200 {object} response.RegistryItemsFromPrefix
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/status/workflow [get]
func workflowStatusGetHandler(ctx *gin.Context) {
	storeGetHandler(ctx, v1alpha1.WorkflowStatusKeyPrefix)
}

// @Summary Returns the specified workflow object from the store with the name in params.
// @Tags status
// @Accept  json
// @Produce json
// @Param name path string true "name of the workflow to get status document for."
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/status/workflow/{name} [get]
func workflowStatusGetByNameHandler(ctx *gin.Context) {
	storeGetByNameHandler(ctx, v1alpha1.WorkflowStatusKeyPrefix)
}

// workflowStatusCreateHandler creates a new workflow status object.
// @Summary Creates a new workflow status in the store.
// @Description This route creates a new workflow status for xene to operate on, if the workflow already exists
// this will update it with the new workflow.
// @Tags status
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Param workflow formData string true "WorkflowStatus manifest to be created."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/status/workflow [post]
func workflowStatusCreateHandler(ctx *gin.Context) {
	workflow := ctx.PostForm("workflow")
	if workflow == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "workflow is a required parameter",
		})
		return
	}

	var wf types.WorkflowStatus
	err := json.Unmarshal([]byte(workflow), &wf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshling: %s", err),
		})
		return
	}

	err = wf.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while validating workflow status spec: %s", err),
		})
		return
	}

	wfData, err := json.Marshal(&wf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while marshalling wrokflow status object: %s", err),
		})
		return
	}
	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wf.Metadata.ObjectMeta.Name),
		wfData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while setting key: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: fmt.Sprintf("%s workflow status created/updated", wf.Metadata.ObjectMeta.Name),
	})
}

// @Summary Patches the specified workflow status from the store.
// @Tags registry
// @Accept json
// @Produce json
// @Param name path string true "Name of the workflow status to be patched."
// @Failure 400 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/status/workflow/{name} [patch]
func workflowStatusPatchHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, response.HTTPError{
		Error: "route not implemented yet",
	})
}

// @Summary Deletes the specified workflow from the store.
// @Description Deletes the workflow status specified by the name parameter, if the workflow is not
// present then an error is thrown.
// @Tags status
// @Accept  json
// @Produce json
// @Param name path string true "Name of the workflow status to be deleted."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/status/workflow/{name} [delete]
func workflowStatusDeleteHandler(ctx *gin.Context) {
	storeDeleteHandler(ctx, v1alpha1.WorkflowStatusKeyPrefix)
}

// @Summary Gather information about all the workflow objects managed by xene.
// @Tags registry
// @Accept  json
// @Produce json
// @Success 200 {array} response.WorkflowInfo
// @Security ApiKeyAuth
// @Router /api/v1/registry/list/workflows [get]
func workflowsListHandler(ctx *gin.Context) {
	resp := []response.WorkflowInfo{}

	store.KVStore.PrefixScanWithFunction(
		context.TODO(),
		v1alpha1.WorkflowKeyPrefix,
		func(kv *v1alpha1.KVPairStruct) {
			var wf v1alpha1.Workflow
			if err := json.Unmarshal([]byte(kv.Value), &wf); err == nil {
				pipelines := []string{}
				for name := range wf.Spec.Pipelines {
					pipelines = append(pipelines, name)
				}

				triggers := []string{}
				for name := range wf.Spec.Triggers {
					triggers = append(triggers, name)
				}

				wfInfo := response.WorkflowInfo{
					Name:      wf.Metadata.GetName(),
					Pipelines: pipelines,
					Triggers:  triggers,
				}

				val, err := store.KVStore.Get(
					context.TODO(), fmt.Sprintf("%s/%s", v1alpha1.WorkflowStatusKeyPrefix, wf.Metadata.GetName()))
				if err == nil {
					var wfStatus v1alpha1.WorkflowStatus
					err = json.Unmarshal(val.Data, &wfStatus)
					if err == nil {
						executors := make(map[string]struct{})
						for _, status := range wfStatus.Pipelines {
							if _, ok := executors[status.Executor]; !ok {
								executors[status.Executor] = struct{}{}
							}
						}

						execs := []string{}
						for name := range executors {
							execs = append(execs, name)
						}

						wfInfo.Agents = execs
					}
				}

				resp = append(resp, wfInfo)
			}
		})

	ctx.JSON(http.StatusOK, resp)
}
