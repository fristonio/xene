package main

import (
	"fmt"
	"os"

	"github.com/fristonio/xene/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the current build of xene.",
	Long: `Displays the version of the current build of xene, this information
include Version, Revision, Git-Branch, BuildUser, BuildDate, go-version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(version.VersionStr,
			version.Info["version"],
			version.Info["revision"],
			version.Info["branch"],
			version.Info["buildUser"],
			version.Info["buildDate"],
			version.Info["goVersion"])
		os.Exit(0)
	},
}
