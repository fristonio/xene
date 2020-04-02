package apiserver

import (
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/auth"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/gin-gonic/gin"
)

// JWTVerficationMiddleware is used to authenticate any request for user
// Sets the ctx 'user' value to the user email
func (s *APIServer) JWTVerficationMiddleware(ctx *gin.Context) {
	token, err := auth.GetTokenFromCtx(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	claims, err := s.authProvider.GetUserClaimsFromToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	ctx.Set(defaults.ContextBucketKey, claims)
	ctx.Next()
}
