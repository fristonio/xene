package oauth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleProviderType     = "google"
)

type googleProvider struct {
	config *oauth2.Config
	state  string
}

var setup sync.Once

// Google is the google auth Provider.
var Google = &googleProvider{}

func (g *googleProvider) Type() oauth.ProviderType {
	return ProviderType(googleProviderType)
}

func (g *googleProvider) Setup() {
	setup.Once(func() {
		g.config = &oauth2.Config{
			ClientID:     conf.ClientID,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  conf.RedirectURL,
			Scopes:       conf.Scopes,
			Endpoint:     google.Endpoint,
		}
	})
}

func (g *googleProvider) GetLoginURL() string {
	g.state = randToken(32)
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
