package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/option"
	"github.com/fristonio/xene/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "xene",
	Short: "xene is an open source workflow builder and executor tool.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&option.ConfigFile, "config", "c",
		"", "Config file for xene to use.")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cmdrefCmd)
	rootCmd.AddCommand(apiServerCmd)
}

func initConfig() {
	if option.ConfigFile == "" {
		option.ConfigFile = defaults.XeneConfigFile
	}
	if utils.FileExists(defaults.XeneConfigFile) {
		y, err := ioutil.ReadFile(defaults.XeneConfigFile)
		if err != nil {
			log.Errorf("error while reading the file: %s", err)
			return
		}
		err = yaml.Unmarshal(y, option.Config)
		if err != nil {
			log.Errorf("error while unmarshling xene config file: %s", err)
		}
	} else {
		log.Infof("No xene config file found")
	}
}
