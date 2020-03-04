package defaults

import "time"

const (
	// XeneBanner is the ascii art acting as a banner for xene.
	XeneBanner string = `
$$\   $$\  $$$$$$\  $$$$$$$\   $$$$$$\
\$$\ $$  |$$  __$$\ $$  __$$\ $$  __$$\
 \$$$$  / $$$$$$$$ |$$ |  $$ |$$$$$$$$ |
 $$  $$<  $$   ____|$$ |  $$ |$$   ____|
$$  /\$$\ \$$$$$$$\ $$ |  $$ |\$$$$$$$\
\__/  \__| \_______|\__|  \__| \_______|
`

	// XeneConfigFile is the file name for xene config.
	XeneConfigFile string = "/etc/xene/xene.yaml"

	// APIServerHost is the default host to bind the API Server to.
	APIServerHost string = "0.0.0.0"

	// APIServerPort is the default port that the apiserver listens to, this can be
	// overridden when initializing the server.
	APIServerPort uint32 = 6060

	// APIServerScheme is the default server scheme to use for the API server.
	APIServerScheme string = "http"

	// APIServerTimeout is the default value of timeout for server shutdown.
	APIServerTimeout time.Duration = 5 * time.Second

	// APIServerUnixSocketPath is the default path where the unix domain socket is created
	// for the xene server.
	APIServerUnixSocketPath string = "/var/run/xene/xene.sock"

	// TimeFormat is the default format of time we use within xene.
	TimeFormat string = "01-Jan-1970:15:04:05 -0700"

	// JWTExpireInterval is the default value of the duration for the issued
	// JWT token to expire.
	// The default for this is 48 Hours ie: 2 days
	JWTExpireInterval time.Duration = 48 * time.Hour

	// JWTAuthType is the default type of authentication we support.
	JWTAuthType string = "Bearer"

	// ContextBucketKey is the key which is set in the context after JWT has been
	// verified from the Authentication header.
	ContextBucketKey string = "user"
)
