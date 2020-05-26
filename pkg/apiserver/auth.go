package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/auth"
	"github.com/fristonio/xene/pkg/auth/rbac"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
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

	// Check if the roles assumed by the user accessing the API allows him to perform
	// the required action and then only let it pass.
	if !rbac.APIServerRBACMap.ValidateAccessI(claims.Roles, ctx.Request.Method, ctx.FullPath()) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPError{
			Error: "Role not valid for the provided API endpoint",
		})
		return
	}

	emailComps := strings.Split(claims.Email, "@")
	if len(emailComps) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPError{
			Error: "Invalid email address in token",
		})
		return
	}

	emailValid := false
	for _, domain := range option.Config.APIServer.AllowedDomains {
		if emailComps[1] == domain {
			emailValid = true
			break
		}
	}

	if !emailValid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPError{
			Error: fmt.Sprintf("Email must conform to a valid domain: %s", option.Config.APIServer.AllowedDomains),
		})
		return
	}

	ctx.Set(defaults.ContextBucketKey, claims)
	ctx.Next()
}
