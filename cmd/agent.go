package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fristonio/xene/pkg/agent"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run xene agent.",
	Long:  "Run xene agent to deploy executor for running configured workflows.",
	Run: func(cmd *cobra.Command, args []string) {
		if option.Config.Agent.Address == "" || option.Config.Agent.APIServer == "" ||
			option.Config.Agent.Name == "" {
			log.Error("--address, --name and --api-server are required flags for agent.")
			os.Exit(1)
		}

		if option.Config.Agent.JWTSecret == "" {
			option.Config.Agent.JWTSecret = utils.RandToken(40)
		}

		if !option.Config.Agent.Insecure {
			if option.Config.Agent.KeyFile == "" ||
				option.Config.Agent.CertFile == "" ||
				option.Config.Agent.RootCACert == "" ||
				option.Config.Agent.ClientCertFile == "" ||
				option.Config.Agent.ClientKeyFile == "" ||
				option.Config.Agent.ServerName == "" {
				log.Fatal("In insecure mode all the key and certificates file should be specifed")
			}
		}
		server := agent.NewServer(option.Config.Agent.Host,
			option.Config.Agent.Port,
			option.Config.Agent.Address,
			option.Config.Agent.CertFile,
			option.Config.Agent.KeyFile,
			option.Config.Agent.RootCACert,
			option.Config.Agent.JWTSecret,
			option.Config.Agent.Insecure)

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go func() {
			<-sigc
			log.Info("Signal recieved, shutting down api server.")
			server.Shutdown()
			os.Exit(0)
		}()

		log.Infof("Registered agent name is: %s", option.Config.Agent.Name)

		// Join the agent pool in the API server
		err := server.JoinAPIServer(
			option.Config.Agent.APIServer,
			option.Config.Agent.Name,
			option.Config.Agent.Address,
			option.Config.Agent.APIAuthToken)
		if err != nil {
			log.Errorf("error while joining API server: %s", err)
			os.Exit(1)
		}

		err = server.RunServer()
		if err != nil {
			log.Errorf("error while running agent server: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	agentFlags := agentCmd.Flags()

	agentFlags.StringVarP(&option.Config.Agent.Host, "host", "b",
		defaults.AgentHost, "Host to bind the agent to.")
	agentFlags.Uint32VarP(&option.Config.Agent.Port, "port", "p",
		defaults.AgentPort, "Port to bind the xene agent grpc server on.")
	agentFlags.StringVarP(&option.Config.Agent.ServerName, "server-name", "",
		"", "server domain name for the certificates to be used by clients.")
	agentFlags.StringVarP(&option.Config.Agent.APIServer, "api-server", "s",
		fmt.Sprintf("%s:%d", defaults.APIServerHost, defaults.APIServerPort),
		"api server address to connect to")
	agentFlags.StringVarP(&option.Config.Agent.APIAuthToken, "api-auth-token", "a",
		"", "Authentication token to use while joining api server.")
	agentFlags.StringVarP(&option.Config.Agent.KeyFile, "key-file", "k",
		"", "Key to use when using secure mode of GRPC protocol.")
	agentFlags.StringVarP(&option.Config.Agent.CertFile, "cert-file", "l",
		"", "Certificate to use for the agent when running in secure GRPC mode.")
	agentFlags.StringVarP(&option.Config.Agent.RootCACert, "root-ca", "r",
		"", "Root CA certificate chain to use for mTLS.")
	agentFlags.StringVarP(&option.Config.Agent.ClientKeyFile, "client-key-file", "",
		"", "Key to use when agent's grpc server.")
	agentFlags.StringVarP(&option.Config.Agent.ClientCertFile, "client-cert-file", "",
		"", "Certificate to use for the client connecting to agent.")
	agentFlags.StringVarP(&option.Config.Agent.Name, "name", "n",
		"", "Name to run the agent with.")
	agentFlags.BoolVarP(&option.Config.Agent.Insecure, "insecure", "i",
		false, "Run agent in insecure mode")
	agentFlags.StringVarP(&option.Config.Agent.Address, "address", "m",
		"", "Own address of the agent, for the API server to communmicate")
	agentFlags.StringVarP(&option.Config.Agent.JWTSecret, "jwt-secret", "j",
		"", "JWT secret to use for authentication purpose for GRPC server")
	agentFlags.StringVarP(&option.Config.Agent.StorageDir, "storage-directory", "d",
		defaults.AgentStorageDir, "Storage directory to use for xene agent.")
	agentFlags.BoolVarP(&option.Config.Agent.LocalLogServer, "local-log-server", "",
		true, "Run xene agent with embedded log server for handling log files.")
	agentFlags.Uint32VarP(&option.Config.Agent.LogServerPort, "log-server-port", "",
		defaults.AgentLogServerPort, "Run xene agent with embedded log server for handling log files.")

	_ = agentCmd.MarkPersistentFlagRequired("name")
	_ = agentCmd.MarkPersistentFlagRequired("address")
	_ = agentCmd.MarkPersistentFlagRequired("api-server")

	viper.SetEnvPrefix("XENE_AGENT")
	viper.AutomaticEnv()
	_ = viper.BindPFlags(agentFlags)
}
