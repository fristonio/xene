package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/fristonio/xene/pkg/agent/controller"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/executor/cre/docker"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/store"

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

	certFile   string
	keyFile    string
	rootCACert string

	insecureMode bool

	server *grpc.Server
}

// NewServer returns a new server for agent.
func NewServer(host string, port uint32, address, certFile, keyFile, rootCACert,
	jwtSecret string, insecureMode bool) *Server {
	return &Server{
		shuttingDown: false,
		address:      address,
		host:         host,
		port:         port,
		certFile:     certFile,
		keyFile:      keyFile,
		rootCACert:   rootCACert,
		insecureMode: insecureMode,
	}
}

// RunServer starts running the GRPC server for xene agent.
func (s *Server) RunServer() error {
	fmt.Println(defaults.XeneBanner)

	// Try connecting to docker container runtime.
	docker.ConnectToDockerOrDie()

	hostPort := fmt.Sprintf("%s:%d", s.host, s.port)

	lis, err := net.Listen("tcp", hostPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	log.Infof("Xene agent server is listening on: %s", hostPort)

	// TODO: Implement GRPC server with MTLs rather then using JWT tokens
	// for request authentication.
	var grpcServer *grpc.Server
	if !s.insecureMode {
		certificate, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
		if err != nil {
			log.Fatalf("error while loading key cert pair for server: %s", err)
		}

		certPool := x509.NewCertPool()
		data, err := ioutil.ReadFile(s.rootCACert)
		if err != nil {
			log.Fatalf("failed to read root ca cert: %s", err)
		}

		ok := certPool.AppendCertsFromPEM(data)
		if !ok {
			log.Fatal("failed to append client certs")
		}
		tlsConfig := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		}

		c := grpc.Creds(credentials.NewTLS(tlsConfig))
		grpcServer = grpc.NewServer(c)
	} else {
		grpcServer = grpc.NewServer()
	}

	s.server = grpcServer
	proto.RegisterAgentServiceServer(grpcServer, newAgentServer())

	// Initialize store for xene agent.
	log.Infof("setting up xene agent store client")
	err = store.Setup(option.Config.Agent.StorageDir)
	if err != nil {
		log.Fatalf("error while initializing xene store: %s", err)
	}

	err = controller.RunControllers()
	if err != nil {
		log.Fatalf("error when running agent controllers: %s", err)
	}

	return grpcServer.Serve(lis)
}

// Shutdown shuts down the agent GRPC server.
func (s *Server) Shutdown() {
	s.shuttingDown = true
	s.server.GracefulStop()
	s.shuttingDown = false
	log.Info("server shutdown successful.")
}
