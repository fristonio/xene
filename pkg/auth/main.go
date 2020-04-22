package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func getTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("missing Authorization Header")
	}
	authTypeLen := len(defaults.JWTAuthType)
	if authHeader[:authTypeLen] != defaults.JWTAuthType {
		return "", fmt.Errorf("invalid authentication type, supported type is: %s", defaults.JWTAuthType)
	}

	return authHeader[authTypeLen+1:], nil
}

// GetTokenFromCtx returns the token string parsing it from the Authorization
// header.
func GetTokenFromCtx(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	return getTokenFromHeader(authHeader)
}

// GetTokenFromGRPCContext returns the JWT token from the grpc context
func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata in the context")
	}

	if len(md["Authorization"]) > 0 {
		return getTokenFromHeader(md["Authorization"][0])
	}

	return "", fmt.Errorf("no Authorization header in grpc request")
}
