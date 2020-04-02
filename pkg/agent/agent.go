package agent

import (
	"fmt"
	"net"

	"github.com/fristonio/xene/pkg/auth/jwt"
	"github.com/fristonio/xene/pkg/defaults"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server is the GRPC server for agent running.
type Server struct {
	shuttingDown bool
	address      string

	host string
	port uint32

	certFile string
	keyFile  string

	insecureMode bool

	server       *grpc.Server
	authProvider *jwt.AuthProvider
}

// NewServer returns a new server for agent.
func NewServer(host string, port uint32, address, certFile, keyFile, jwtSecret string, insecureMode bool) *Server {
	return &Server{
		shuttingDown: false,
		address:      address,
		host:         host,
		port:         port,
		certFile:     certFile,
		keyFile:      keyFile,
		insecureMode: insecureMode,
		authProvider: jwt.NewJWTAuthProvider(jwtSecret),
	}
}

// RunServer starts running the GRPC server for xene agent.
func (s *Server) RunServer() error {
	fmt.Println(defaults.XeneBanner)

	hostPort := fmt.Sprintf("%s:%d", s.host, s.port)

	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	log.Infof("Xene API server is listening on: %s", hostPort)

	// TODO: Implement GRPC server with MTLs rather then using JWT tokens
	// for request authentication.
	var grpcServer *grpc.Server
	if !s.insecureMode {
		creds, err := credentials.NewServerTLSFromFile(s.certFile, s.keyFile)
		if err != nil {
			return fmt.Errorf("failed to create credentials: %v", err)
		}

		c := grpc.Creds(creds)
		grpcServer = grpc.NewServer(c, grpc.UnaryInterceptor(s.JWTVerficationMiddleware))
	} else {
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(s.JWTVerficationMiddleware))
	}

	s.server = grpcServer
	return grpcServer.Serve(lis)
}

// Shutdown shuts down the agent GRPC server.
func (s *Server) Shutdown() {
	s.shuttingDown = true
	s.server.GracefulStop()
	s.shuttingDown = false
	log.Info("server shutdown successful.")
}
