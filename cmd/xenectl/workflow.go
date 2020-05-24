package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fristonio/xene/pkg/apiserver/client/registry"
	"github.com/fristonio/xene/pkg/executor"
	"github.com/fristonio/xene/pkg/executor/cre/docker"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

		if res.Payload.Item.Value == "" {
			log.Infof("the requested workflow is not found")
			return
		}

		prettyPrintJSON(res.Payload.Item.Value)
	},
}

var workflowDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the workflow with the provided name.",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowName == "" {
			log.Fatalf("workflow name(--name) is a required parameter and must be a valid one")
		}

		client, auth := getClientAndAuth()
		res, err := client.Registry.DeleteAPIV1RegistryWorkflowName(
			registry.NewDeleteAPIV1RegistryWorkflowNameParams().WithName(workflowName),
			auth)
		if err != nil {
			log.Errorf("error while deleting workflow document: %s", err)
			return
		}

		fmt.Println(res.Payload.Message)
	},
}

var workflowRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run start running the workflow with the file name provided",

	Run: func(cmd *cobra.Command, args []string) {
		if workflowFileName == "" {
			log.Fatalf("workflow file name(--file) is a required parameter and must be a valid one")
		}

		if !runLocal {
			log.Fatalf("Currently only running workflow locally is supported")
		}

		// Try connecting to docker for the execution of pipelines
		// else Die
		docker.ConnectToDockerOrDie()

		data, err := ioutil.ReadFile(workflowFileName)
		if err != nil {
			log.Fatalf("error while reading workflow manifest file: %s", err)
		}

		var workflow v1alpha1.Workflow
		err = json.Unmarshal(data, &workflow)
		if err != nil {
			log.Fatalf("error while unmarshalling workflow manifest: %s", err)
		}

		err = workflow.Validate()
		if err != nil {
			log.Fatalf("workflow is not valid: %s", err)
		}
		err = workflow.Resolve()
		if err != nil {
			log.Fatalf("error while resolving workflow: %s", err)
		}

		for name, pipeline := range workflow.Spec.Pipelines {
			log.Infof("\nProcessing pipeline: %s\n", name)

			spec := &v1alpha1.PipelineSpecWithName{
				PipelineSpec: *pipeline,
				Name:         name,
				Workflow:     workflow.Metadata.GetName(),
			}

			id := fmt.Sprintf("%s-%s", name, utils.RandToken(10))
			exec := executor.NewPipelineExecutor(name, id, spec).WithoutStore()
			status := v1alpha1.PipelineRunStatus{
				Name:      name,
				RunID:     id,
				Status:    "Running",
				Agent:     "XENECTL",
				StartTime: time.Now().Unix(),
			}
			exec.Run(status)

			data, err := json.Marshal(exec.GetStatus())
			if err != nil {
				log.Fatalf("error while marshalling pipeline run status: %s", err)
			}

			log.Infof("\n\nPipeline processing finished")
			log.Infof("----------------- STATUS REPORT -----------------")
			prettyPrintJSON(string(data))
		}
	},
}

var (
	workflowFileName string
	workflowName     string
	runLocal         bool
)

func init() {
	workflowCreateCmd.Flags().StringVarP(&workflowFileName, "file", "f",
		"", "File to use for workflow manfiest.")
	workflowGetCmd.Flags().StringVarP(&workflowName, "name", "n",
		"", "name of the workflow definition to get.")
	workflowDeleteCmd.Flags().StringVarP(&workflowName, "name", "n",
		"", "name of the workflow definition to delete.")
	workflowRunCmd.Flags().StringVarP(&workflowFileName, "file", "f",
		"", "File to use for workflow manfiest.")
	workflowRunCmd.Flags().BoolVarP(&runLocal, "local", "l",
		true, "Run the workflow pipelines from the manifest locally")

	workflowCmd.AddCommand(workflowCreateCmd)
	workflowCmd.AddCommand(workflowGetCmd)
	workflowCmd.AddCommand(workflowDeleteCmd)
	workflowCmd.AddCommand(workflowRunCmd)
}
