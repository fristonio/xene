# Installation

Xene can be deployed in multiple ways, the choice of the method depends on the use case of the user.

## Manual Installation

For manual installation make sure that you have the following packages already installed on the host
operating system:

1. [Docker](https://docs.docker.com/get-docker/)

Pick a release of your choice from the [release page](https://github.com/fristonio/xene/releases) and download the binaries.

```bash
# Download xene and xenectl binaries
$ curl -fsSLO https://github.com/fristonio/xene/releases/download/v0.1.0/xene-linux-amd64
$ curl -fsSLO https://github.com/fristonio/xene/releases/download/v0.1.0/xenectl-linux-amd64

$ mv xene xenectl /usr/local/bin/

$ xene version

$$\   $$\  $$$$$$\  $$$$$$$\   $$$$$$\
\$$\ $$  |$$  __$$\ $$  __$$\ $$  __$$\
 \$$$$  / $$$$$$$$ |$$ |  $$ |$$$$$$$$ |
 $$  $$<  $$   ____|$$ |  $$ |$$   ____|
$$  /\$$\ \$$$$$$$\ $$ |  $$ |\$$$$$$$\
\__/  \__| \_______|\__|  \__| \_______|


    Version    : 0.1
    Revision   : de336e3
    Branch     : master
    Build-User : fristonio@fristonio
    Build-Date : 20200614-14:14:00
    Go-Version : 1.13.7
```

**NOTE:** xene agent and apiserver can run on different machines similar as to what shown below for a single machine.

Before we start running xene components, let's set up the environment.

Create a file `/etc/xene/conf/xene.yaml` and add agent and APIServer configuration to it from the template given below.

- [APIServer Config](/_examples/sample.apiserver.config.yaml)
- [Agent Config](/_examples/sample.agent.config.yaml)

To be able to use `xenectl` for interacting with Xene APIServer one must provide API server address and the authentication token.
You can create a config file to automatically be used for doing interaction with API server.

Template for `xenectl` config is present [here](/_examples/sample.xenectl.yaml)

**NOTE** For a distributed steup of Xene you can deploy agents on different machines each with it's own config. In such cases make
sure that your agents are able to reach your API server for initial registration.

If you want to disable authentication on either of APIService or Agent you can configure it using the config file.

Jump to [Configure Services](###ConfigureServices) if you don't want to set up authentication between agent and APIServer.

For all the agents, generate certificates using the below commands

```bash
$ git clone https://github.com/fristonio/xene && cd xene/contrib/certs/

$ make certs && \
    mkdir -p /etc/xene/certs/ && \
    mv *.gen /etc/xene/certs/
```

Server certificates are sent to the API server at the time of registration. For this registration we need to make
request to the APIServer for which we need authentication token for agent. To generate this authentication token
use the JWTSecret provided in the xene agent config and sign a token.

### Configure Services

Create systemd services on your host system using the templates below:

- [APIServer Service](/contrib/systemd/xene-apiserver.service)
- [Agent Service](/contrib/systemd/xene-agent.service)

```bash
$ mv contrib/systemd/* /etc/systemd/system/

$ systemctl start xene-apiserver
$ systemctl start xene-agent
```

**NOTE:** Run commands on different nodes if you want to set up a distributed system.

Your APIServer should be up and running at this point ready to run workflows. If want a visual interface for Xene,
you can take a look at the [setting up Xene UI](UI.md).

You can also interact with the APIServer using our command line tool `xenectl`. A few xenectl command examples is shown
below:

```bash
$ xenectl workflow create --file _examples/workflow.json
INFO[0000] TestWorkflow workflow created/updated

$ xenectl workflow get --name TestWorkflow
{
    "kind": "Workflow",
    "apiVersion": "v1alpha1",
    "metadata": {
        "name": "TestWorkflow",
        "description": "A workflow definition to test xene workflows"
    }
    ...
    ...
}

$ xenectl workflow delete --name TestWorkflow
store item(TestWorkflow) has been deleted
```

## Docker Compose Install

For docker compose installation make sure that you have the following packages already installed on the host
operating system:

1. [Docker](https://docs.docker.com/get-docker/)
2. [Docker Compose](https://docs.docker.com/compose/install/)

```bash
$ git clone https://github.com/fristonio/xene && cd xene

$ make docker
$ make docker-compose-up
```

**NOTE:** Do not use this method in production, its for local usage only.

Your APIServer should be up and running at this point ready to run workflows. If want a visual interface for Xene,
you can take a look at the [setting up Xene UI](UI.md).
