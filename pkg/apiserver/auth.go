package apiserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/gin-gonic/gin"
)

// getTokenFromCtx returns the token string parsing it from the Authorization
// header.
func getTokenFromCtx(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization Header")
	}
	authTypeLen := len(defaults.JWTAuthType)
	if authHeader[:authTypeLen] != defaults.JWTAuthType {
		return "", fmt.Errorf("invalid authentication type, supported type is: %s", defaults.JWTAuthType)
	}

	return authHeader[authTypeLen+1:], nil
}

// JWTVerficationMiddleware is used to authenticate any request for user
// Sets the ctx 'user' value to the user email
func (s *APIServer) JWTVerficationMiddleware(ctx *gin.Context) {
	token, err := getTokenFromCtx(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	claims, err := s.authProvider.GetClaimsFromToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	ctx.Set(defaults.ContextBucketKey, claims)
	ctx.Next()
}
