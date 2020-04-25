package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var cmdrefDir string

var cmdrefCmd = &cobra.Command{
	Use:   "cmdref",
	Short: "Run xenectl cmdref command to generate command line reference.",
	Long: "Run xenectl cmdref command to generate command line reference " +
		"documentaiton in the specified directory, by default docs/cmdref/xenectl/.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("generating command line reference documentation in directory %s", cmdrefDir)
		err := doc.GenMarkdownTree(rootCmd, cmdrefDir)
		if err != nil {
			log.Errorf("error while generating cmdref for xene: %s", err)
		}
	},
}

func init() {
	cmdrefCmd.Flags().StringVarP(&cmdrefDir, "directory", "d",
		"docs/cmdref/xenectl", "Directory to use for creating cmd reference docs")
}
