package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/store"
	types "github.com/fristonio/xene/pkg/types/v1alpha1"

	"github.com/gin-gonic/gin"
)

// RegistryGroupRoutes sets up gin router for registry api group
// The base path of this router gropup is `/api/v1/registry/`
func RegistryGroupRoutes(r *gin.RouterGroup) {
	workflowGroup := r.Group("workflow")
	workflowGroupRoutes(workflowGroup)
}

// workflowGroupRoutes sets up the routes for workflow registry group.
func workflowGroupRoutes(r *gin.RouterGroup) {
	r.GET("", workflowGetHandler)
	r.GET("/:name", workflowGetByNameHandler)
	r.POST("", workflowCreateHandler)
	r.PATCH("", workflowPatchHandler)
	r.DELETE("/:name", workflowDeleteHandler)
}

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
// @Success 200 {object} response.WorkflowsFromPrefix
// @Success 200 {object} response.Workflow
// @Failure 500 {object} response.HTTPError
// @Router /api/v1/registry/workflow [get]
func workflowGetHandler(ctx *gin.Context) {
	prefix := ctx.Query("prefix")
	name := ctx.Query("name")

	if prefix == "" && name == "" {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: "One of prefix or name is required to get a workflow",
		})
		return
	}

	// prefix takes preference over name
	if prefix != "" {
		kvPairs, err := store.KVStore.ListPrefix(context.TODO(), prefix)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.HTTPError{
				Error: fmt.Sprintf("error while listing prefixes for %s: %s", prefix, err),
			})
			return
		}

		resp := response.WorkflowsFromPrefix{
			Count: len(kvPairs),
		}
		wf := []types.KVPairStruct{}
		for key, value := range kvPairs {
			wf = append(wf, types.KVPairStruct{
				Key:   key,
				Value: string(value.Data),

				Version:          value.Version,
				ExpiresAt:        value.ExpiresAt,
				DeletedOrExpired: value.DeletedOrExpired,
			})
		}
		resp.Workflows = wf
		ctx.JSON(http.StatusOK, resp)
		return
	}

	val, err := store.KVStore.Get(context.TODO(), name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while retriving value from store: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Workflow{
		Workflow: types.KVPairStruct{
			Key:   name,
			Value: string(val.Data),

			Version:          val.Version,
			ExpiresAt:        val.ExpiresAt,
			DeletedOrExpired: val.DeletedOrExpired,
		},
	})
}

// @Summary Returns the specified workflow object from the store with the name in params.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "name of the workflow to get."
// @Success 200 {object} response.Workflow
// @Failure 500 {object} response.HTTPError
// @Router /api/v1/registry/workflow/{name} [get]
func workflowGetByNameHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "name parameter is required",
		})
	}
	val, err := store.KVStore.Get(context.TODO(), name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while getting the key: %s: %s", name, err),
		})
		return
	}

	resp := response.Workflow{
		Workflow: types.KVPairStruct{
			Key:   name,
			Value: string(val.Data),

			Version:          val.Version,
			ExpiresAt:        val.ExpiresAt,
			DeletedOrExpired: val.DeletedOrExpired,
		},
	}

	ctx.JSON(http.StatusOK, resp)
}

// workflowCreateHandler creates a new workflow.
// @Summary Creates a new workflow in the store.
// @Description This route creates a new workflow for xene to operate on, if the workflow already exists
// this will update it with the new workflow.
// @Tags registry
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Param workflow formData string true "Workflow manifest to be created."
// @Example {"kind": "workflow", "apiVersion": "v1", "spec": {"triggers": [], "pipelines": []}, "metadata": {"name": "abc"}
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
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

	if wf.Metadata.ObjectMeta.Name == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("name is a required field in workflow description"),
		})
		return
	}
	err = store.KVStore.Set(context.TODO(), wf.Metadata.ObjectMeta.Name, []byte(workflow))
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
// @Router /api/v1/registry/workflow/{name} [delete]
func workflowDeleteHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "workflow name is required in the route",
		})
		return
	}

	err := store.KVStore.Delete(context.TODO(), name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while deleting workflow: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: fmt.Sprintf("workflow(%s) has been deleted", name),
	})
}
