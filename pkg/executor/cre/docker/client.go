package docker

import (
	"context"

	dockerapi "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{
	"executor": "docker-container",
})

// Get a *dockerapi.Client, either using
// DOCKER_HOST, DOCKER_TLS_VERIFY, and DOCKER_CERT path per their spec
func getDockerClient() (*dockerapi.Client, error) {
	return dockerapi.NewEnvClient()
}

// ConnectToDockerOrDie tries to connect to the docker and prints the version
func ConnectToDockerOrDie() {
	client, err := getDockerClient()
	if err != nil {
		log.Fatalf("Error while creating docker client: %s", err)
	}

	ver, err := client.ServerVersion(context.TODO())
	if err != nil {
		log.Fatalf("error while retrieving docker runtime information: %s", err)
	}

	log.Infof("DOCKER SERVER VERSION: \n%v", ver)
}
