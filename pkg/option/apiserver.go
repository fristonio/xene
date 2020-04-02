package option

// APIServerOpts is the type used for configuration option of xene apiserver.
type APIServerOpts struct {
	// Standalone is the boolean value which represents if xene apiserver is running
	// in a standalone mode and thus the agent is running locally.
	Standalone bool `json:"standalone"`

	// Host is the host to bind the apiserver on.
	Host string `json:"host"`

	// Port is the port to run the apiserver on.
	Port uint32 `json:"port"`

	// Scheme to user for APIServer, it can be one of http, https or unix
	Scheme string `json:"scheme"`

	// DisableAuth disables authentication for apiserver api routes.
	DisableAuth bool `json:"disableAuth"`

	// VerboseLogs enables verbose logging for the apiserver.
	VerboseLogs bool `json:"verboseLogs"`

	// KeyFile is the key file to use for apiserver HTTPS scheme.
	KeyFile string `json:"keyFile"`

	// CertFile contains the certificate for apiserver HTTPS scheme.
	CertFile string `json:"certFile"`

	// UnixSocketPath is the unix socket for apiserver unix socket scheme.
	UnixSocketPath string `json:"unixSocketPath"`

	// JWTSecret is the secret to use for issuing JWT tokens by the API server.
	JWTSecret string `json:"jwtSecret"`

	// Oauth contains oauth providers configuration for the API server, multiple
	// oauth providers can be configured at the same time.
	Oauth map[string]OauthProvider `json:"oauth"`
}

// OauthProvider is the authentication provider configuration for Oauth2 based
// authentication for xene.
type OauthProvider struct {
	// Name is the name of the oauth provider
	Name string `json:"name"`

	// ClientID of the oauth provider
	ClientID string `json:"clientID"`

	// ClientSecret is the secret for the oauth provider
	ClientSecret string `json:"clientSecret"`

	// RedirectURL is the redirect URL of the oauth provider.
	RedirectURL string `json:"redirectURL"`

	// Scopes are the scopes associated with the oauth provider to access.
	Scopes []string `json:"scopes"`
}
