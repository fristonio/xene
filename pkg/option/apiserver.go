package option

// APIServerOpts is the type used for configuration option of xene apiserver.
type APIServerOpts struct {
	// Standalone is the boolean value which represents if xene apiserver is running
	// in a standalone mode and thus the agent is running locally.
	Standalone bool `yaml:"standalone"`

	// Host is the host to bind the apiserver on.
	Host string `yaml:"host"`

	// Port is the port to run the apiserver on.
	Port uint32 `yaml:"port"`

	// Scheme to user for APIServer, it can be one of http, https or unix
	Scheme string `yaml:"scheme"`

	// DisableAuth disables authentication for apiserver api routes.
	DisableAuth bool `yaml:"disableAuth"`

	// VerboseLogs enables verbose logging for the apiserver.
	VerboseLogs bool `yaml:"verboseLogs"`

	// KeyFile is the key file to use for apiserver HTTPS scheme.
	KeyFile string `yaml:"keyFile"`

	// CertFile contains the certificate for apiserver HTTPS scheme.
	CertFile string `yaml:"certFile"`

	// UnixSocketPath is the unix socket for apiserver unix socket scheme.
	UnixSocketPath string `yaml:"unixSocketPath"`

	// JWTSecret is the secret to use for issuing JWT tokens by the API server.
	JWTSecret string `yaml:"jwtSecret"`

	// AllowedDomains contains a list of domain which are allowed to be logged in using the
	// OAUTH provider
	AllowedDomains []string `yaml:"allowedDomains"`

	// Oauth contains oauth providers configuration for the API server, multiple
	// oauth providers can be configured at the same time.
	Oauth map[string]OauthProvider `yaml:"oauth"`
}

// OauthProvider is the authentication provider configuration for Oauth2 based
// authentication for xene.
type OauthProvider struct {
	// Name is the name of the oauth provider
	Name string `yaml:"name"`

	// ClientID of the oauth provider
	ClientID string `yaml:"clientID"`

	// ClientSecret is the secret for the oauth provider
	ClientSecret string `yaml:"clientSecret"`

	// RedirectURL is the redirect URL of the oauth provider.
	RedirectURL string `yaml:"redirectURL"`

	// Scopes are the scopes associated with the oauth provider to access.
	Scopes []string `yaml:"scopes"`
}
