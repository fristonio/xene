package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fristonio/xene/pkg/apiserver/client/status"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var workflowStatusCmd = &cobra.Command{
	Use:   "workflowstatus",
	Short: "Subcommand for managing xene configured workflows status",

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var workflowStatusCreateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a workflow status object from the provided file",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowStatusFileName == "" {
			log.Fatalf("workflow file name(--file) is a required parameter and must be a valid file")
		}

		data, err := ioutil.ReadFile(workflowStatusFileName)
		if err != nil {
			log.Fatalf("error while reading file: %s", err)
		}

		client, auth := getClientAndAuth()
		res, err := client.Status.PostAPIV1StatusWorkflow(
			status.NewPostAPIV1StatusWorkflowParams().WithWorkflow(string(data)),
			auth)
		if err != nil {
			log.Errorf("error while getting workflow document: %s", err)
			return
		}

		log.Infof(res.Payload.Message)
	},
}

var workflowStatusGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the workflow status with the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowName == "" {
			log.Fatalf("workflow name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Status.GetAPIV1StatusWorkflowName(
			status.NewGetAPIV1StatusWorkflowNameParams().WithName(workflowName),
			auth)
		if err != nil {
			log.Errorf("error while getting workflow document: %s", err)
			return
		}

		var kv v1alpha1.KVPairStruct
		err = json.Unmarshal([]byte(res.Payload.Item), &kv)
		if err != nil {
			log.Errorf("error while unmarshalling json: %s", err)
		}
		prettyPrintJSON(string(kv.Value))
	},
}

var workflowStatusDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Get the workflow with the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowName == "" {
			log.Fatalf("workflow name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Status.DeleteAPIV1StatusWorkflowName(
			status.NewDeleteAPIV1StatusWorkflowNameParams().WithName(workflowName),
			auth)
		if err != nil {
			log.Errorf("error while getting workflow document: %s", err)
			return
		}

		fmt.Println(res.Payload.Message)
	},
}

var (
	workflowStatusFileName string
)

func init() {
	workflowStatusCreateCmd.Flags().StringVarP(&workflowStatusFileName, "file", "f",
		"", "File to use for workflow manfiest.")
	workflowStatusGetCmd.Flags().StringVarP(&workflowName, "name", "n",
		"", "name of the workflow definition to get.")
	workflowStatusDeleteCmd.Flags().StringVarP(&workflowName, "name", "n",
		"", "name of the workflow definition to delete.")

	workflowStatusCmd.AddCommand(workflowStatusCreateCmd)
	workflowStatusCmd.AddCommand(workflowStatusGetCmd)
	workflowStatusCmd.AddCommand(workflowStatusDeleteCmd)
}
