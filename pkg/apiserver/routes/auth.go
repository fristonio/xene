package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fristonio/xene/pkg/apiserver/response"
	"github.com/fristonio/xene/pkg/auth"
	"github.com/fristonio/xene/pkg/auth/jwt"
	"github.com/fristonio/xene/pkg/auth/oauth"
	"github.com/fristonio/xene/pkg/defaults"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AuthGroupRouter sets up the routes for auth group api endpoints.
// This includes the following endpoints:
// * /oauth/refresh
// * /oauth/{provider}
// * /oauth/{provider}/refresh
func AuthGroupRouter(r *gin.RouterGroup, ap *jwt.AuthProvider) {
	r.GET("/refresh", refreshToken(ap))

	oauth.ConfigureProviders()
	for _, provider := range oauth.ProvidersList {
		log.Infof("%s oauth provider configured", provider.Type())

		loginRoute := string(provider.Type())
		redirectRoute := fmt.Sprintf("%s/redirect", string(provider.Type()))

		r.GET(loginRoute, loginHandler(provider))
		r.GET(redirectRoute, redirectHandler(provider, ap))
	}
}

// refreshToken handles the refresh route for oauth login handlers.
// @Summary Handle authentication token refresh for the oauth provider.
// @Description Handles authentication token refresh
// @Tags auth
// @Accept  json
// @Produce json
// @Success 200 {object} response.JWTAuth
// @Failure 500 {object} response.HTTPError
// @Router /oauth/refresh/ [get]
func refreshToken(ap *jwt.AuthProvider) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token, err := auth.GetTokenFromCtx(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		claims, err := ap.GetUserClaimsFromToken(token)
		// Since we're refreshing token, we don't care if the previous token
		// claims are valid or invalid
		if err != nil && err != jwt.ErrTokenExpired {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}

		// We need to check if the time since previous expired token is
		// under the max refresh time
		if time.Now().Unix() > claims.ExpiresAt+int64(defaults.JWTRefreshInterval) {
			ctx.JSON(http.StatusUnauthorized, response.HTTPError{
				Error: "Time exceeds max refresh time, login again",
			})
			return
		}

		refreshedToken, err := ap.NewAuthToken(claims.Email, claims.Name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.HTTPError{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.JWTAuth{
			Token:     refreshedToken,
			ExpiresIn: defaults.JWTExpireInterval / time.Second,
			UserEmail: claims.Email,
			UserName:  claims.Name,
		})
	}
}

// loginHandler handles the login for a particular user.
// @Summary Handles login for xene
// @Description Log in to xene using the configured oauth providers that xene supports.
// @Tags auth
// @Accept  json
// @Produce json
// @Success 200 {object} response.OauthLogin
// @Param provider path string true "Provider for oauth login"
// @Router /oauth/{provider} [get]
func loginHandler(provider oauth.Provider) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, response.OauthLogin{
			LoginURL: provider.GetLoginURL(),
		})
	}
}

// redirectHandler handles the redirect from the Oauth provider after the authentication
// The request will contain if the user has been authenticated or not with the relevant scopes
// details.
// @Summary Handles redirect from the login oauth provider.
// @Description redirectHandler handles the redirect from the Oauth provider after the authentication process has
// completed
// @Tags auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider for the oauth redirect"
// @Success 200 {object} response.JWTAuth
// @Failure 500 {object} response.HTTPError
// @Router /oauth/{provider}/redirect [get]
func redirectHandler(provider oauth.Provider, ap *jwt.AuthProvider) func(*gin.Context) {
	return func(ctx *gin.Context) {
		u, code, err := provider.GetUser(ctx)
		if err != nil {
			ctx.JSON(code, response.HTTPError{
				Error: err.Error(),
			})
			return
		}

		// Commit the user here to the database so that we can refresh
		// the token later and have the user info in our database.
		tok, err := ap.NewAuthToken(u.Email, u.Profile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.HTTPError{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.JWTAuth{
			Token:     tok,
			ExpiresIn: defaults.JWTExpireInterval / time.Hour,
			UserEmail: u.Email,
			UserName:  u.Profile,
		})
	}
}
