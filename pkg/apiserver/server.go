package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fristonio/xene/pkg/defaults"
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
	router *gin.Engine
	server *http.Server

	shuttingDown bool

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
		scheme: schemeHTTP,
		host:   host,
		port:   port,

		shuttingDown: false,
	}
}

// NewUnixSocketServer returns the ApiServer configured for running
// HTTP server on a unix domain socket specified by the file.
func NewUnixSocketServer(file string) *APIServer {
	return &APIServer{
		scheme:         schemeUnix,
		unixDomainPath: file,

		shuttingDown: false,
	}
}

// NewHTTPSServer returns a new ApiServer capable of running HTTPs server
// on the specified port.
func NewHTTPSServer(host string, port uint32, keyFile, certFile string) *APIServer {
	return &APIServer{
		scheme:   schemeHTTPS,
		host:     host,
		port:     port,
		certFile: certFile,
		keyFile:  keyFile,

		shuttingDown: false,
	}
}

// RunServer runs the server configured for the API.
func (s *APIServer) RunServer() error {
	log.Info(defaults.XeneBanner)
	s.router = NewAPIServerRouter(true)

	switch s.scheme {
	case schemeHTTP:
		hostPort := fmt.Sprintf("%s:%d", s.host, s.port)
		s.server = &http.Server{
			Addr:    hostPort,
			Handler: s.router,
		}
		log.Infof("Xene API server is listening on: %s", hostPort)
		return s.server.ListenAndServe()
	case schemeUnix, schemeHTTPS:
		log.Warnf("the api server scheme %s is not suppported yet", s.scheme)
	}

	return nil
}

// Shutdown shuts down the api server running.
func (s *APIServer) Shutdown() {
	s.shuttingDown = true

	ctx, cancel := context.WithTimeout(context.Background(), defaults.APIServerTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("error while trying to shut down the apiserver: %s", err)
		s.shuttingDown = false
		return
	}

	s.shuttingDown = false
	log.Info("server shutdown successful.")
}
