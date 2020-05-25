package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/fristonio/xene/pkg/agent/controller"
	"github.com/fristonio/xene/pkg/apiserver"
	"github.com/fristonio/xene/pkg/dag"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/errors"
	"github.com/fristonio/xene/pkg/executor/cre/docker"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/proto"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

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

	if option.Config.Agent.LocalLogServer {
		err := s.runLocalLogServer()
		if err != nil {
			return err
		}
	}

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

// Run local log server for xene.
func (s *Server) runLocalLogServer() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(gin.Recovery())

	r.Use(apiserver.NewXeneLoggerMiddleware(log.New(), false))

	hostPort := fmt.Sprintf("%s:%d", option.Config.Agent.Host, option.Config.Agent.LogServerPort)
	server := &http.Server{
		Addr:    hostPort,
		Handler: r,
	}
	log.Infof("Xene log server is listening on: %s", hostPort)

	r.GET("/workflow/:workflow/pipeline/:pipeline/runs/:runID/logs", func(ctx *gin.Context) {
		workflow := ctx.Param("workflow")
		pipeline := ctx.Param("pipeline")
		runID := ctx.Param("runID")

		val, err := store.KVStore.Get(context.TODO(),
			fmt.Sprintf(
				"%s/%s",
				v1alpha1.PipelineKeyPrefix,
				v1alpha1.GetWorkflowPrefixedName(workflow, pipeline)))

		if err != nil {
			log.Errorf("error getting pipeline spec from store: %s", err)
			ctx.Abort()
			return
		}

		pipelineSpec := v1alpha1.PipelineSpecWithName{}
		err = json.Unmarshal(val.Data, &pipelineSpec)
		if err != nil {
			log.Errorf("error while unmarshalling pipeline spec on agent: %s", err)
			ctx.Abort()
			return
		}

		v, err := store.KVStore.Get(context.TODO(),
			fmt.Sprintf(
				"%s/%s/%s",
				v1alpha1.PipelineStatusKeyPrefix,
				v1alpha1.GetWorkflowPrefixedName(workflow, pipeline),
				runID))

		if err != nil {
			log.Errorf("error getting pipeline run status from store: %s", err)
			ctx.Abort()
			return
		}

		pipelineStatus := v1alpha1.PipelineRunStatus{}
		err = json.Unmarshal(v.Data, &pipelineStatus)
		if err != nil {
			log.Errorf("error while unmarshalling pipeline spec on agent: %s", err)
			ctx.Abort()
			return
		}

		err = pipelineSpec.Resolve(pipeline)
		if err != nil {
			log.Errorf("error while resolving pipeline spec: %s", err)
		}

		logFiles := []struct {
			path, task, step string
		}{}
		// Walk each of task in the pipeline in the required order.
		merrs := pipelineSpec.Dag.Walk(func(v dag.Vertex) *errors.MultiError {
			errs := errors.NewMultiError()
			task, ok := v.(*v1alpha1.TaskSpec)
			if !ok {
				// If any error is getting the task spec then append it to the list of errors
				errs.Append(fmt.Errorf("not a valid vertex to visit, must confirm to type *TaskSpec"))
				return errs
			}

			taskStatus := pipelineStatus.Tasks[task.Name()]
			for stepName, step := range taskStatus.Steps {
				if step.LogFile != "" {
					logFiles = append(logFiles, struct {
						path, task, step string
					}{
						path: step.LogFile,
						task: task.Name(),
						step: stepName,
					})
				}
			}

			return errs
		})

		log.Errorf("error during dag walk: %v", merrs)

		ctx.Header("Content-Type", "text/plain")
		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s-%s-%s.log"`, workflow, pipeline, runID))

		var readers []io.Reader
		var contentLength int64

		for _, logFile := range logFiles {
			cr := strings.NewReader(fmt.Sprintf(`
    ***************************************************************************
    Workflow: %s
    Pipeline: %s
    RunID: %s
    Task: %s
    Step: %s
    ***************************************************************************
`, workflow, pipeline, runID, logFile.task, logFile.step))
			f, err := os.Open(logFile.path)
			if err != nil {
				r := strings.NewReader(fmt.Sprintf("\nError while opening logfile path: %s\n", err))
				readers = append(readers, r)
				contentLength += r.Size()
				continue
			}

			info, err := os.Stat(logFile.path)
			if err != nil {
				r := strings.NewReader(fmt.Sprintf("\nError while stat on logfile path: %s\n", err))
				readers = append(readers, r)
				contentLength += r.Size()
				continue
			}

			readers = append(readers, cr)
			readers = append(readers, f)
			contentLength += cr.Size() + info.Size()
		}

		reader := io.MultiReader(readers...)
		ctx.DataFromReader(http.StatusOK, contentLength, "text/plain", reader, nil)
	})

	r.StaticFS("/logs/workflow/", http.Dir(defaults.AgentLogsDir))

	go server.ListenAndServe()
	return nil
}
