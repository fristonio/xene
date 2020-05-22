package main

import (
	"fmt"

	"github.com/fristonio/xene/pkg/apiserver/client/registry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Subcommand for managing xene configured agents",

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var agentGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the agent from the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if agentName == "" {
			log.Fatalf("agent name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Registry.GetAPIV1RegistryAgentName(
			registry.NewGetAPIV1RegistryAgentNameParams().WithName(agentName),
			auth)
		if err != nil {
			log.Errorf("error while getting agent document: %s", err)
			return
		}

		if res.Payload.Item.Value == "" {
			log.Infof("the requested agent is not found")
			return
		}

		prettyPrintJSON(res.Payload.Item.Value)
	},
}

var agentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the agent with the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if agentName == "" {
			log.Fatalf("agent name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Registry.DeleteAPIV1RegistryAgentName(
			registry.NewDeleteAPIV1RegistryAgentNameParams().WithName(agentName),
			auth)
		if err != nil {
			log.Errorf("error while deleting agent: %s", err)
			return
		}

		fmt.Println(res.Payload.Message)
	},
}

var (
	agentName string
)

func init() {
	agentGetCmd.Flags().StringVarP(&agentName, "name", "n",
		"", "name of the agent definition to get.")
	agentDeleteCmd.Flags().StringVarP(&agentName, "name", "n",
		"", "name of the agent definition to delete.")

	agentCmd.AddCommand(agentGetCmd)
	agentCmd.AddCommand(agentDeleteCmd)
}
