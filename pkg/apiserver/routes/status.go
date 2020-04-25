package routes

import (
	"github.com/gin-gonic/gin"
)

// StatusGroupRoutes sets up gin router for status api group
// The base path of this router gropup is `/api/v1/status/`
func StatusGroupRoutes(r *gin.RouterGroup) {
	workflowStatusGroup := r.Group("workflow")
	workflowStatusGroupRoutes(workflowStatusGroup)
}

// workflowGroupRoutes sets up the routes for workflow registry group.
func workflowStatusGroupRoutes(r *gin.RouterGroup) {
	r.GET("", workflowStatusGetHandler)
	r.GET("/:name", workflowStatusGetByNameHandler)
	r.POST("", workflowStatusCreateHandler)
	r.PATCH("", workflowStatusPatchHandler)
	r.DELETE("/:name", workflowStatusDeleteHandler)
}
