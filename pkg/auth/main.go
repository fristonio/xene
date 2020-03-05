package auth

import (
	"errors"
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/gin-gonic/gin"
)

// GetTokenFromCtx returns the token string parsing it from the Authorization
// header.
func GetTokenFromCtx(ctx *gin.Context) (string, error) {
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
