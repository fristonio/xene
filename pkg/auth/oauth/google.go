package oauth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleProviderType     = "google"
)

type googleProvider struct {
	config     *oauth2.Config
	state      string
	configured bool
}

// Google is the google auth Provider.
var Google = &googleProvider{}

// Type returns the type of provider represented by google provider.
func (g *googleProvider) Type() ProviderType {
	return ProviderType(googleProviderType)
}

// Configured returns true if the google oauth provider has been configured.
func (g *googleProvider) Configured() bool {
	return g.configured
}

// Configure sets up the google oauth provider for xene.
func (g *googleProvider) Configure() {
	conf, ok := option.Config.APIServer.Oauth["googleProviderType"]
	if !ok {
		log.Warnf("google provider not configured, no config provided")
		return
	}
	setup.Do(func() {
		g.config = &oauth2.Config{
			ClientID:     conf.ClientID,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  conf.RedirectURL,
			Scopes:       conf.Scopes,
			Endpoint:     google.Endpoint,
		}
		g.configured = true
	})
}

// GetLoginURL returns the login URL for a client.
func (g *googleProvider) GetLoginURL() string {
	g.state = utils.RandToken(32)
	return g.config.AuthCodeURL(g.state)
}

// GetUser returns the User from the provided gin context.
func (g *googleProvider) GetUser(ctx *gin.Context) (*User, int, error) {
	token, err := g.config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	client := g.config.Client(context.TODO(), token)
	info, err := client.Get(googleUserInfoEndpoint)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer info.Body.Close() //nolint:errcheck

	data, err := ioutil.ReadAll(info.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	u := new(User)
	if err := json.Unmarshal(data, u); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return u, http.StatusOK, nil
}
