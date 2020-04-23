package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/gin-gonic/gin"
)

// agentGetHandler retrivies a given agent object from the store.
// @Summary Returns the specified agent object from the store.
// @Description If a name is provided return the corresponding agent object, if prefix is set to some value
// then serach based on that prefix, if only name is provided retrive the associated
// agent object
// @Tags registry
// @Accept  json
// @Produce json
// @Param prefix query string false "Prefix based get for agent."
// @Param name query string false "name of the agent to get."
// @Success 200 {object} response.RegistryItemsFromPrefix
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/agent [get]
func agentGetHandler(ctx *gin.Context) {
	registryGetHandler(ctx, v1alpha1.AgentKeyPrefix)
}

// @Summary Returns the specified agent object from the store with the name in params.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "name of the agent to get."
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/agent/{name} [get]
func agentGetByNameHandler(ctx *gin.Context) {
	registryGetByNameHandler(ctx, v1alpha1.AgentKeyPrefix)
}

// @Summary Creates a new agent in the store.
// @Tags registry
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Param agent formData string true "Agent manifest to be created."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/registry/agent [post]
func agentRegisterHandler(ctx *gin.Context) {
	newAgent := ctx.PostForm("agent")
	if newAgent == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "agent is a required parameter",
		})
		return
	}

	var agent v1alpha1.Agent
	err := json.Unmarshal([]byte(newAgent), &agent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshling: %s", err),
		})
		return
	}

	err = agent.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while validating agent spec: %s", err),
		})
		return
	}
	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.AgentKeyPrefix, agent.Metadata.ObjectMeta.Name),
		[]byte(newAgent))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while setting key: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: fmt.Sprintf("%s agent created/updated", agent.Metadata.ObjectMeta.Name),
	})
}

// @Summary Patches the specified agent from the store.
// @Tags registry
// @Accept json
// @Produce json
// @Param name path string true "Name of the agent to be patched."
// @Failure 400 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/agent/{name} [patch]
func agentPatchHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, response.HTTPError{
		Error: "route not implemented yet",
	})
}

// @Summary Deletes the specified agent from the store.
// @Description Deletes the agent specified by the name parameter, if the agent is not
// present then an error is thrown.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "Name of the agent to be deleted."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/registry/agent/{name} [delete]
func agentRemoveHandler(ctx *gin.Context) {
	registryDeleteHandler(ctx, v1alpha1.AgentKeyPrefix)
}
