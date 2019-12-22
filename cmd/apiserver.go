package main

import (
	"fmt"
	"os"

	"github.com/fristonio/xene/pkg/apiserver"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiServerCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "Run xene apiserver.",
	Long:  "Run xene apiserver which can then be used to communicate to user facing interface of xene.",

	Run: func(cmd *cobra.Command, args []string) {
		server := apiserver.NewHTTPServer(option.APIServer.Host, option.APIServer.Port)

		err := server.RunServer()
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	apiServerFlags := apiServerCmd.Flags()

	apiServerFlags.StringVarP(&option.APIServer.ConfigFile, "config", "c", "", "Config file for API server.")
	apiServerFlags.StringVarP(&option.APIServer.Host, "host", "b", defaults.APIServerHost, "Host to bind the api server to.")
	apiServerFlags.Uint32VarP(&option.APIServer.Port, "port", "p", defaults.APIServerPort, "Port to bind the xene api server on.")
	apiServerFlags.StringVarP(&option.APIServer.Scheme, "scheme", "s", defaults.APIServerScheme, "Scheme to use for the api server.")
	apiServerFlags.StringVarP(&option.APIServer.UnixSocketPath, "unix-socket", "u", defaults.APIServerUnixSocketPath, "Default path for the unix domain socket, when using unix scheme")
	apiServerFlags.StringVarP(&option.APIServer.KeyFile, "key-file", "k", "", "Key to use when using HTTPS scheme for the server.")
	apiServerFlags.StringVarP(&option.APIServer.CertFile, "cert-file", "l", "", "Certificate to use for the API Server when running under HTTPS scheme.")

	viper.BindPFlags(apiServerFlags)
}

func initAPIServerConfig() {
	fmt.Println(option.APIServer.ConfigFile)
	if option.APIServer.ConfigFile != "" {
		viper.SetConfigFile(option.APIServer.ConfigFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Errorf("error while getting home directory: %s", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".xene.conf")
	}

	viper.SetEnvPrefix("XENE_APISERVER")
	viper.AutomaticEnv()

	var err error
	if err = viper.ReadInConfig(); err == nil {
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	}
	log.Debugf("Error while reading config file for viper: %s", err)
}
