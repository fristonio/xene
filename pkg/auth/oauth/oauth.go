package oauth

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var setup sync.Once

// User represents the user details required from the oauth provider.
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ProviderType for different OAuth Providers supported in xene
type ProviderType string

// Provider is an interface that represents any provider
// through which we can authenticate a user on the application.
type Provider interface {
	// Type returns the type of provider.
	Type() ProviderType

	// Configured checks if the provider is configured or not.
	Configured() bool

	// Configure helps us to initialize any variables or setup requirements.
	Configure()

	// GetLoginURL returns the URL which redirects user to the providers login page.
	GetLoginURL() string

	// GetUser gets the user after requesting the OAuth provider.
	// Returns the user, status code and error if any.
	GetUser(*gin.Context) (*User, int, error)
}

var (
	// ProvidersList is the list of all the providers configured for the API
	// server, this is populated by ConfigureProviders function.
	ProvidersList []Provider
)

// ConfigureProviders configures all the provider available and registers them
// in the ProvidersList.
func ConfigureProviders() {
	// Configure google provider
	Google.Configure()
	if Google.Configured() {
		ProvidersList = append(ProvidersList, Google)
	}
}
