package main

import (
	"os"

	"github.com/fristonio/xene/pkg/apiserver"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/spf13/cobra"
)

var apiServerCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "Run xene apiserver.",
	Long:  "Run xene apiserver which can then be used to communicate to user facing interface of xene.",

	Run: func(cmd *cobra.Command, args []string) {
		server := apiserver.NewHTTPServer(defaults.APIServerHost, defaults.APIServerPort)

		err := server.RunServer()
		if err != nil {
			os.Exit(1)
		}
	},
}
