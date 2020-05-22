package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/store"
	"github.com/gin-gonic/gin"
)

// RegistryGroupRoutes sets up gin router for registry api group
// The base path of this router gropup is `/api/v1/registry/`
func RegistryGroupRoutes(r *gin.RouterGroup) {
	workflowGroup := r.Group("workflow")
	workflowGroupRoutes(workflowGroup)

	agentGroup := r.Group("agent")
	agentGroupRoutes(agentGroup)

	secretsGroup := r.Group("secret")
	secretsGroupRoutes(secretsGroup)

	listGroup := r.Group("list")
	listGroupRoutes(listGroup)
}

// workflowGroupRoutes sets up the routes for workflow registry group.
func workflowGroupRoutes(r *gin.RouterGroup) {
	r.GET("", workflowGetHandler)
	r.GET("/:name", workflowGetByNameHandler)
	r.POST("", workflowCreateHandler)
	r.PATCH("", workflowPatchHandler)
	r.DELETE("/:name", workflowDeleteHandler)
}

// agentsGroupRoutes sets up the routes for agent registry group.
func agentGroupRoutes(r *gin.RouterGroup) {
	r.GET("", agentGetHandler)
	r.GET("/:name", agentGetByNameHandler)
	r.POST("", agentRegisterHandler)
	r.PATCH("", agentPatchHandler)
	r.DELETE("/:name", agentRemoveHandler)
}

// secretsGroupRoutes sets up the routes for secrets registry group.
func secretsGroupRoutes(r *gin.RouterGroup) {
	r.GET("", secretGetHandler)
	r.GET("/:name", secretGetByNameHandler)
	r.POST("", secretRegisterHandler)
	r.PATCH("", secretPatchHandler)
	r.DELETE("/:name", secretRemoveHandler)
}

func listGroupRoutes(r *gin.RouterGroup) {
	r.GET("agents", agentsListHandler)
	r.GET("secrets", secretsListHandler)
	r.GET("workflows", workflowsListHandler)
}

func storeGetHandler(ctx *gin.Context, pre string) {
	prefix := ctx.Query("prefix")
	name := ctx.Query("name")

	if prefix == "" && name == "" {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: "One of prefix or name is required to get a registry item",
		})
		return
	}

	// prefix takes preference over name
	if prefix != "" {
		kvPairs, err := store.KVStore.ListPrefix(context.TODO(), pre)
		if store.KVStore.KeyDoesNotExistError(err) {
			ctx.JSON(http.StatusOK, response.RegistryItemsFromPrefix{
				Count: 0,
			})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.HTTPError{
				Error: fmt.Sprintf("error while listing prefixes for %s: %s", prefix, err),
			})
			return
		}

		resp := response.RegistryItemsFromPrefix{
			Count: len(kvPairs),
		}
		wf := make([]response.KVPair, 0)
		for key, value := range kvPairs {
			data := response.KVPair{
				Key:   key,
				Value: string(value.Data),
			}
			wf = append(wf, data)
		}
		resp.Items = wf
		ctx.JSON(http.StatusOK, resp)
		return
	}

	val, err := store.KVStore.Get(context.TODO(), fmt.Sprintf("%s/%s", pre, name))
	if err != nil {
		if store.KVStore.KeyDoesNotExistError(err) {
			ctx.JSON(http.StatusOK, response.RegistryItem{
				Item: response.KVPair{},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while retriving value from store: %s", err),
		})
		return
	}

	data := response.KVPair{
		Key:   name,
		Value: string(val.Data),
	}
	ctx.JSON(http.StatusOK, response.RegistryItem{
		Item: data,
	})
}

func storeGetByNameHandler(ctx *gin.Context, prefix string) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "name parameter is required",
		})
	}
	val, err := store.KVStore.Get(context.TODO(), fmt.Sprintf("%s/%s", prefix, name))
	if err != nil {
		if store.KVStore.KeyDoesNotExistError(err) {
			ctx.JSON(http.StatusOK, response.RegistryItem{
				Item: response.KVPair{},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while getting the key: %s: %s", name, err),
		})
		return
	}

	data := response.KVPair{
		Key:   name,
		Value: string(val.Data),
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshaling objects %s: %s", name, err),
		})
		return
	}
	resp := response.RegistryItem{
		Item: data,
	}

	ctx.JSON(http.StatusOK, resp)
}

func storeDeleteHandler(ctx *gin.Context, prefix string) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "registry item name is required in the route",
		})
		return
	}

	err := store.KVStore.Delete(context.TODO(), fmt.Sprintf("%s/%s", prefix, name))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while deleting store item: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: fmt.Sprintf("store item(%s) has been deleted", name),
	})
}
