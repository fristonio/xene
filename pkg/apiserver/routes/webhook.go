package routes

import "github.com/gin-gonic/gin"

// WebhookGroupRoutes contains all the routes related to xene info.
func WebhookGroupRoutes(r *gin.RouterGroup) {
	r.GET("/trigger/workflow/:workflow/pipeline/:pipeline", webhookTriggerHandler)
}

// @Summary Webhook endpoints to trigger the pipelines of a particular workflow.
// @Tags webhook
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Param workflow path string true "Name of the workflow."
// @Param pipeline path string true "Name of the pipeline to be triggered."
// @Failure 500 {object} response.HTTPError
// @Failure 400 {object} response.HTTPError
// @Router /api/v1/webhook/trigger/{workflow}/pipeline/{pipeline} [get]
func webhookTriggerHandler(ctx *gin.Context) {
}
