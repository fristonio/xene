package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/fristonio/xene/pkg/apiserver/client"
	"github.com/fristonio/xene/pkg/apiserver/client/registry"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getClientAndAuth() (*client.XeneAPIServer, runtime.ClientAuthInfoWriter) {
	addr, err := url.Parse(config.APIServerAddr)
	if err != nil {
		log.Fatalf("invalid api address URL: %s", err)
	}
	client := client.New(httptransport.New(addr.Host, "", nil), strfmt.Default)
	bearerTokenAuth := httptransport.BearerToken(config.AuthToken)

	return client, bearerTokenAuth
}

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Subcommand for managing xene configured workflows",

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var workflowCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workflow from the provided file",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowFileName == "" {
			log.Fatalf("workflow file name(--file) is a required parameter and must be a valid file")
		}

		data, err := ioutil.ReadFile(workflowFileName)
		if err != nil {
			log.Fatalf("error while reading file: %s", err)
		}

		client, auth := getClientAndAuth()
		res, err := client.Registry.PostAPIV1RegistryWorkflow(
			registry.NewPostAPIV1RegistryWorkflowParams().WithWorkflow(string(data)),
			auth)
		if err != nil {
			log.Errorf("error while getting workflow document: %s", err)
			return
		}

		log.Infof(res.Payload.Message)
	},
}

var workflowGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the workflow from the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowName == "" {
			log.Fatalf("workflow name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Registry.GetAPIV1RegistryWorkflowName(
			registry.NewGetAPIV1RegistryWorkflowNameParams().WithName(workflowName),
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

var workflowDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Get the workflow with the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowName == "" {
			log.Fatalf("workflow name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Registry.DeleteAPIV1RegistryWorkflowName(
			registry.NewDeleteAPIV1RegistryWorkflowNameParams().WithName(workflowName),
			auth)
		if err != nil {
			log.Errorf("error while getting workflow document: %s", err)
			return
		}

		fmt.Println(res.Payload.Message)
	},
}

var (
	workflowFileName string
	workflowName     string
)

func init() {
	workflowCreateCmd.Flags().StringVarP(&workflowFileName, "file", "f",
		"", "File to use for workflow manfiest.")
	workflowGetCmd.Flags().StringVarP(&workflowName, "name", "n",
		"", "name of the workflow definition to get.")
	workflowDeleteCmd.Flags().StringVarP(&workflowName, "name", "n",
		"", "name of the workflow definition to delete.")

	workflowCmd.AddCommand(workflowCreateCmd)
	workflowCmd.AddCommand(workflowGetCmd)
}
