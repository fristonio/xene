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

// secretGetHandler retrivies a given secret object from the store.
// @Summary Returns the specified secret object from the store.
// @Description If a name is provided return the corresponding secret object, if prefix is set to some value
// then serach based on that prefix, if only name is provided retrive the associated
// secret object
// @Tags registry
// @Accept  json
// @Produce json
// @Param prefix query string false "Prefix based get for secret."
// @Param name query string false "name of the secret to get."
// @Success 200 {object} response.RegistryItemsFromPrefix
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/secret [get]
func secretGetHandler(ctx *gin.Context) {
	registryGetHandler(ctx, v1alpha1.SecretKeyPrefix)
}

// @Summary Returns the specified secret object from the store with the name in params.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "name of the secret to get."
// @Success 200 {object} response.RegistryItem
// @Failure 500 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/secret/{name} [get]
func secretGetByNameHandler(ctx *gin.Context) {
	registryGetByNameHandler(ctx, v1alpha1.SecretKeyPrefix)
}

// @Summary Creates a new secret in the store.
// @Tags registry
// @Accept  application/x-www-form-urlencoded
// @Produce json
// @Param secret formData string true "secret manifest to be created."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/registry/secret [post]
func secretRegisterHandler(ctx *gin.Context) {
	newSecret := ctx.PostForm("secret")
	if newSecret == "" {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: "secret is a required parameter",
		})
		return
	}

	var secret v1alpha1.Secret
	err := json.Unmarshal([]byte(newSecret), &secret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("error while unmarshling: %s", err),
		})
		return
	}

	err = secret.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: fmt.Sprintf("secret is not valid: %s", err),
		})
		return
	}

	err = store.KVStore.Set(context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.SecretKeyPrefix, secret.Metadata.ObjectMeta.Name),
		[]byte(newSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: fmt.Sprintf("error while setting key: %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPMessage{
		Message: fmt.Sprintf("%s secret created/updated", secret.Metadata.ObjectMeta.Name),
	})
}

// @Summary Patches the specified secret from the store.
// @Tags registry
// @Accept json
// @Produce json
// @Param name path string true "Name of the secret to be patched."
// @Failure 400 {object} response.HTTPError
// @Security ApiKeyAuth
// @Router /api/v1/registry/secret/{name} [patch]
func secretPatchHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, response.HTTPError{
		Error: "route not implemented yet",
	})
}

// @Summary Deletes the specified secret from the store.
// @Description Deletes the secret specified by the name parameter, if the secret is not
// present then an error is thrown.
// @Tags registry
// @Accept  json
// @Produce json
// @Param name path string true "Name of the secret to be deleted."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/registry/secret/{name} [delete]
func secretRemoveHandler(ctx *gin.Context) {
	registryDeleteHandler(ctx, v1alpha1.SecretKeyPrefix)
}
