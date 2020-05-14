package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fristonio/xene/pkg/apiserver/controller/agent"
	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/proto/protoutils"
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
	storeGetHandler(ctx, v1alpha1.AgentKeyPrefix)
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
	storeGetByNameHandler(ctx, v1alpha1.AgentKeyPrefix)
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

	agentData, err := json.Marshal(&agent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while marshalling agent object: %s", err),
		})
		return
	}
	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.AgentKeyPrefix, agent.Metadata.ObjectMeta.Name),
		agentData)
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
	storeDeleteHandler(ctx, v1alpha1.AgentKeyPrefix)
}

// @Summary List all the keys of items in the registry of the provided type agent.
// @Tags registry
// @Accept  json
// @Produce json
// @Success 200 {array} response.AgentInfo
// @Security ApiKeyAuth
// @Router /api/v1/registry/list/agents [get]
func agentsListHandler(ctx *gin.Context) {
	resp := []response.AgentInfo{}
	operating := agent.AgentCtrl.GetActiveAgentsWithInfo()
	blacklisted := agent.AgentCtrl.GetBlacklistedAgentsWithInfo()

	store.KVStore.PrefixScanWithFunction(
		context.TODO(),
		v1alpha1.AgentKeyPrefix,
		func(kv *v1alpha1.KVPairStruct) {
			ag := strings.TrimPrefix(kv.Key, v1alpha1.AgentKeyPrefix+"/")
			if agInfo, ok := operating[ag]; ok {
				resp = append(resp, response.AgentInfo{
					Name:      ag,
					Address:   agInfo.Spec.Address,
					Secure:    !agInfo.Spec.Insecure,
					Available: true,
				})
			}

			if agInfo, ok := blacklisted[ag]; ok {
				resp = append(resp, response.AgentInfo{
					Name:      ag,
					Address:   agInfo.Spec.Address,
					Secure:    !agInfo.Spec.Insecure,
					Available: false,
				})
			}
		})

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Returns verbose information about the agent.
// @Tags info
// @Accept  json
// @Produce json
// @Param name path string true "Name of the agent to get information about."
// @Success 200 {object} response.AgentVerboseInfo
// @Failure 400 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/info/agent/{name} [get]
func agentInfoHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "agent name is required in the route",
		})
		return
	}

	conn := agent.AgentCtrl.AgentConnection(name)
	if conn == nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: "Agent connection not found",
		})
		return
	}

	client := proto.NewAgentServiceClient(conn)
	info, err := client.Info(context.TODO(), &proto.AgentInfoOpts{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: "Error getting agent info",
		})
		return
	}

	ctx.JSON(http.StatusOK, protoutils.GetAgentVerboseInfoFromProtoAgentInfo(info))
}
