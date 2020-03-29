package apiserver

import (
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/apiserver/routes"
	"github.com/fristonio/xene/pkg/store"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

// NewAPIServerRouter returns a new gin router for xene api server.
// The parameter defines if the router should be initialized with a logger.
func (s *APIServer) NewAPIServerRouter(includeLogger bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Initialize store for xene apiserver.
	err := store.Setup()
	if err != nil {
		log.Fatalf("error while initializing xene store: %s", err)
	}

	if includeLogger {
		r.Use(NewXeneLoggerMiddleware(log.New(), s.verboseLogs))
	}

	authGroup := r.Group("/oauth")
	apiV1Group := r.Group("/api/v1")
	if !s.disableAuth {
		apiV1Group.Use(s.JWTVerficationMiddleware)
	}

	r.Use(gin.Recovery())

	r.GET("/docs/api/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.setupHealthRoute(r)

	routes.AuthGroupRouter(authGroup, s.authProvider)
	routes.APIGroupRouter(apiV1Group)

	return r
}

// @Summary Health route for Xene API server.
// @Description Returns the health status of the API server.
// @Tags health
// @Accept  json
// @Produce json
// @Success 200 {object} response.HTTPMessage
// @Router /health [get]
func (s *APIServer) setupHealthRoute(r *gin.Engine) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, response.HTTPMessage{
			Message: "Healthy",
		})
	})
}
