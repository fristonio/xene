package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fristonio/xene/pkg/apiserver"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiServerCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "Run xene apiserver.",
	Long:  "Run xene apiserver which can then be used to communicate to user facing interface of xene.",

	Run: func(cmd *cobra.Command, args []string) {
		if !option.Config.APIServer.DisableAuth && option.Config.APIServer.JWTSecret == "" {
			log.Error("Either specify disable-auth flag or provide a JWT secret to use for authentication.")
			os.Exit(1)
		}
		server := apiserver.NewHTTPServer(option.Config.APIServer.Host,
			option.Config.APIServer.Port,
			option.Config.APIServer.DisableAuth,
			option.Config.APIServer.VerboseLogs,
			option.Config.APIServer.JWTSecret)

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
	apiServerFlags := apiServerCmd.Flags()

	apiServerFlags.StringVarP(&option.Config.APIServer.Host, "host", "b",
		defaults.APIServerHost, "Host to bind the api server to.")
	apiServerFlags.Uint32VarP(&option.Config.APIServer.Port, "port", "p",
		defaults.APIServerPort, "Port to bind the xene api server on.")
	apiServerFlags.StringVarP(&option.Config.APIServer.Scheme, "scheme", "s",
		defaults.APIServerScheme, "Scheme to use for the api server.")
	apiServerFlags.BoolVarP(&option.Config.APIServer.DisableAuth, "disable-auth", "n",
		false, "If the authentication should be disabled for the API server.")
	apiServerFlags.StringVarP(&option.Config.APIServer.UnixSocketPath, "unix-socket", "u",
		defaults.APIServerUnixSocketPath, "Default path for the unix domain socket, when using unix scheme")
	apiServerFlags.StringVarP(&option.Config.APIServer.KeyFile, "key-file", "k",
		"", "Key to use when using HTTPS scheme for the server.")
	apiServerFlags.StringVarP(&option.Config.APIServer.CertFile, "cert-file", "l",
		"", "Certificate to use for the API Server when running under HTTPS scheme.")
	apiServerFlags.StringVarP(&option.Config.APIServer.JWTSecret, "jwt-secret", "j",
		"", "JWT secret for authentication purposes, make sure it is secure and non bruteforcable.")
	apiServerFlags.BoolVarP(&option.Config.APIServer.VerboseLogs, "verbose-logs", "v",
		false, "Print verbose APIServer request logs.")
	apiServerFlags.StringVarP(&option.Config.Store.Engine, "storage-engine", "e",
		defaults.StorageEngineBadger, "Storage engine to use for the API server")
	apiServerFlags.StringVarP(&option.Config.Store.StorageDirectory, "storage-directory", "d",
		defaults.StorageDir, "Storage directory to use for xene apiserver.")

	viper.SetEnvPrefix("XENE_APISERVER")
	viper.AutomaticEnv()
	_ = viper.BindPFlags(apiServerFlags)
}
