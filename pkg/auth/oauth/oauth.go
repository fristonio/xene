package oauth

import "github.com/gin-gonic/gin"

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

	// Setup enables us to initialize any variables or setup requirements.
	Setup() error

	// GetLoginURL returns the URL which redirects user to the providers login page.
	GetLoginURL() string

	// GetUser gets the user after requesting the OAuth provider.
	// Returns the user, status code and error if any.
	GetUser(*gin.Context) (*User, int, error)
}
