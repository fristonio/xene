package routes

import (
	"github.com/gin-gonic/gin"
)

// AuthGroupRouter sets up the routes for auth group api endpoints.
// This includes the following endpoints:
// * /auth/refresh
// * /auth/{provider}
// * /auth/{provider}/refresh
func AuthGroupRouter(r *gin.RouterGroup) {
	return
}

// refreshToken handles the referesh route for oauth login handlers.
// @Summary Handle authentication token refresh for the oauth provider.
// @Description Handles authentication token refresh
// @Tags auth
// @Accept  json
// @Produce json
// @Router /auth/refresh/ [post]
func refreshToken(c *gin.Context) {
	return
}
