package apiserver

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// NewAPIServerRouter returns a new gin router for xene api server.
// The parameter defines if the router should be initialized with a logger.
func NewAPIServerRouter(includeLogger bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	if includeLogger {
		r.Use(NewXeneLoggerMiddleware(log.New()))
	}
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return r
}
