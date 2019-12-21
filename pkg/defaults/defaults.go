package defaults

const (
	// APIServerHost is the default host to bind the API Server to.
	APIServerHost string = "0.0.0.0"

	// APIServerPort is the default port that the apiserver listens to, this can be
	// overridden when initializing the server.
	APIServerPort uint32 = 6060

	// APIServerUnixSocketPath is the default path where the unix domain socket is created
	// for the xene server.
	APIServerUnixSocketPath string = "/var/run/xene/xene.sock"

	// TimeFormat is the default format of time we use within xene.
	TimeFormat string = "01-Jan-1970:15:04:05 -0700"
)
