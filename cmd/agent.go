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
		if option.Config.Agent.JWTSecret == "" {
			option.Config.Agent.JWTSecret = utils.RandToken(40)
		}
		server := agent.NewServer(option.Config.Agent.Host,
			option.Config.Agent.Port,
			option.Config.Agent.Address,
			option.Config.Agent.CertFile,
			option.Config.Agent.KeyFile,
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

		err := server.RunServer()
		if err != nil {
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
	agentFlags.StringVarP(&option.Config.Agent.APIServer, "api-server", "s",
		fmt.Sprintf("%s:%s", defaults.APIServerHost, defaults.APIServerPort),
		"api server address to connect to")
	agentFlags.StringVarP(&option.Config.Agent.APIAuthToken, "api-auth-token", "a",
		"", "Authentication token to use while joining api server.")
	agentFlags.StringVarP(&option.Config.Agent.KeyFile, "key-file", "k",
		"", "Key to use when using secure mode of GRPC protocol.")
	agentFlags.StringVarP(&option.Config.Agent.CertFile, "cert-file", "l",
		"", "Certificate to use for the agent when running in secure GRPC mode.")
	agentFlags.BoolVarP(&option.Config.Agent.Insecure, "insecure", "i",
		true, "Run agent in insecure mode")
	agentFlags.StringVarP(&option.Config.Agent.Address, "address", "m",
		"", "Own address of the agent, for the API server to communmicate")
	agentFlags.StringVarP(&option.Config.Agent.JWTSecret, "jwt-secret", "j",
		"", "JWT secret to use for authentication purpose for GRPC server")
	agentCmd.MarkPersistentFlagRequired("address")
	agentCmd.MarkPersistentFlagRequired("api-server")

	viper.SetEnvPrefix("XENE_AGENT")
	viper.AutomaticEnv()
	_ = viper.BindPFlags(agentFlags)
}
