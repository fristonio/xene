package option

// APIServer is the type used for configuration option of xene apiserver.
type APIServerOpts struct {
	Host           string
	Port           uint32
	Scheme         string
	KeyFile        string
	CertFile       string
	ConfigFile     string
	UnixSocketPath string
}

// APIServer is the configuration option for running the xene api server.
var APIServer = &APIServerOpts{}
