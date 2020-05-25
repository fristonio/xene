package routes

import (
	"github.com/gin-gonic/gin"
)

// InfoGroupRoutes contains all the routes related to xene info.
func InfoGroupRoutes(r *gin.RouterGroup) {
	r.GET("/agent/:name", agentInfoHandler)
	r.GET("/workflow/:workflow/pipeline/:pipeline", pipelineInfoHandler)
	r.GET("/workflow/:workflow/pipeline/:pipeline/runs/:runID", pipelineRunInfoHandler)
}
