package apiserver

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type serverScheme string

const (
	schemeHTTP  serverScheme = "http"
	schemeHTTPS serverScheme = "https"
	schemeUnix  serverScheme = "unix"
)

// APIServer is the standard type which is used to manage the
// API exposed by xene.
type APIServer struct {
	shutdown  chan struct{}
	interrupt chan os.Signal

	shuttingDown bool
	interrupted  bool

	scheme         serverScheme
	unixDomainPath string
	host           string
	port           uint32

	certFile string
	keyFile  string
}

// NewHTTPServer returns the ApiServer configured for running
// HTTP API server.
func NewHTTPServer(host string, port uint32) *APIServer {
	return &APIServer{
		scheme:    schemeHTTP,
		host:      host,
		port:      port,
		shutdown:  make(chan struct{}),
		interrupt: make(chan os.Signal, 1),

		shuttingDown: false,
		interrupted:  false,
	}
}

// NewUnixSocketServer returns the ApiServer configured for running
// HTTP server on a unix domain socket specified by the file.
func NewUnixSocketServer(file string) *APIServer {
	return &APIServer{
		scheme:         schemeUnix,
		unixDomainPath: file,
		shutdown:       make(chan struct{}),
		interrupt:      make(chan os.Signal, 1),

		shuttingDown: false,
		interrupted:  false,
	}
}

// NewHTTPSServer returns a new ApiServer capable of running HTTPs server
// on the specified port.
func NewHTTPSServer(host string, port uint32, keyFile, certFile string) *APIServer {
	return &APIServer{
		scheme:    schemeHTTPS,
		host:      host,
		port:      port,
		certFile:  certFile,
		keyFile:   keyFile,
		shutdown:  make(chan struct{}),
		interrupt: make(chan os.Signal, 1),

		shuttingDown: false,
		interrupted:  false,
	}
}

// RunServer runs the server configured for the API.
func (s *APIServer) RunServer() error {
	r := gin.New()
	r.Use(NewXeneLoggerMiddleware(log.New()))
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	switch s.scheme {
	case schemeHTTP:
		hostPort := fmt.Sprintf("%s:%d", s.host, s.port)
		return r.Run(hostPort)
	case schemeUnix, schemeHTTPS:
		log.Warnf("the api server scheme %s is not suppported yet", s.scheme)
	}

	return nil
}
