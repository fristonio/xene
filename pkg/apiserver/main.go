package apiserver

import (
	"context"
	"fmt"
	"net/http"

	// Swagger docs
	_ "github.com/fristonio/xene/pkg/apiserver/docs"

	"github.com/fristonio/xene/pkg/auth/jwt"
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
	disableAuth  bool
	verboseLogs  bool

	scheme         serverScheme
	unixDomainPath string
	host           string
	port           uint32

	certFile string
	keyFile  string

	authProvider *jwt.AuthProvider
}

// NewHTTPServer returns the ApiServer configured for running
// HTTP API server.
func NewHTTPServer(host string, port uint32, disableAuth, verboseLogs bool, jwtSecret string) *APIServer {
	return &APIServer{
		scheme: schemeHTTP,
		host:   host,
		port:   port,

		shuttingDown: false,
		disableAuth:  disableAuth,
		verboseLogs:  verboseLogs,

		authProvider: jwt.NewJWTAuthProvider(jwtSecret),
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
func NewHTTPSServer(host string, port uint32, keyFile, certFile, jwtSecret string,
	disableAuth, verboseLogs bool) *APIServer {
	return &APIServer{
		scheme:   schemeHTTPS,
		host:     host,
		port:     port,
		certFile: certFile,
		keyFile:  keyFile,

		shuttingDown: false,
		disableAuth:  disableAuth,
		verboseLogs:  verboseLogs,

		authProvider: jwt.NewJWTAuthProvider(jwtSecret),
	}
}

// RunServer runs the server configured for the API.
// @title Xene API server
// @version 0.1.0
// @description Xene is the workflow creator and manager tool
//
// @contact.name Deepesh Pathak
// @contact.url https://dpathak.co
// @contact.email deepeshpathak09@gmail.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:6060
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (s *APIServer) RunServer() error {
	fmt.Println(defaults.XeneBanner)
	s.router = s.NewAPIServerRouter(true)

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
