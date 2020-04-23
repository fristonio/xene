package routes

import (
	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/gin-gonic/gin"
)

// APIGroupRouter sets up the router for /api group
func APIGroupRouter(r *gin.RouterGroup) {
	r.GET("authok", authOKHandler)

	registryGroup := r.Group("registry")
	RegistryGroupRoutes(registryGroup)
}

// @Summary Handles ping event for api routes.
// @Description auth OK handler handles the ping to api routes which are
// protected by authentication.
// @Tags api
// @Accept  json
// @Produce json
// @Success 200 {object} response.HTTPMessage
// @Security ApiKeyAuth
// @Router /api/v1/authok [get]
func authOKHandler(ctx *gin.Context) {
	ctx.JSON(200, response.HTTPMessage{
		Message: "OK",
	})
}
