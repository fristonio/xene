# Xene

> A distributed workflow runner with a focus on performance and simplicity.

[![Netlify Status](https://api.netlify.com/api/v1/badges/f3adc406-ad04-4059-ad21-6a54f4be6771/deploy-status)](https://app.netlify.com/sites/sad-thompson-bcaa9a/deploys) ![Travis Status](https://travis-ci.com/fristonio/xene.svg?token=xvk2YsyqhEExfPszH3rV&branch=master) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


Xene is a high performance, distributed workflow execution platform. It lets you create custom workflows in the form of a Directed Acyclic Graph, which can then be executed based on a trigger configured with the workflow.

Xene is built with Golang with a focus on performance, resiliency, and simplicity. Xene internally runs in a way similar to Kubernetes by getting triggered on a level and then reconciling the desired state with the current state. Currently, xene supports [Badger](https://github.com/dgraph-io/badger) as the storage engines, which is pluggable and can be replaced with a similar key-value store.

## Features

Xene aims to tackle two main problems in any workflow execution system:

1. Low resource footprint
2. Flexibility and Simplicity

Xene exploits a lot of small enhancement to reduce the network and memory footprint of the executors and APIServer.
APIServer acts as the brains of the whole ecosystem by directing the decisions of scheduling workloads and managing coordination
between all the moving components. All the communication between APIServer and Agent is encrypted using mTLS and xene doesn't require
Agents to communicate with each other reducing network footprint.

Xene derives a lot of its design decisions from Kubernetes, like:

- Like Kubernetes, xene is also level triggered and does the reconciliation of state described by the user with the actual state.
- It also uses a KeyValue store to store all the states and objects. To reduce the memory footprint, xene supports the use of [Badger](https://github.com/dgraph-io/badger) as a key-value store.
- Objects are reconciled using controllers configured to watch changes on keys in the KVStore.

Some of the main features of Xene are:

- Locally running workflow from provided specs.
- Command line tool for xene(xenectl) to interact with apiserver.
- mTLS based authentication between agent and apiserver.
- JWT based authentication for API server.
- Google OAuth based authentication integration.
- [Beta UI](https://github.com/fristonio/xene-ui) based on React and Typescript
- Workflow creation and running on top of Docker based container runtime executor on the agent.
- Secrets management for agents.
- Workflow pipelines run live log streaming and status feed.

## Contributing

If you'd like to contribute to this project, refer to the [contributing guide](Contributing.md).

You can start with setting up Xene on your system and try solving a few bugs listed here: https://github.com/fristonio/xene/issues

## Documentation

- [API Docs](https://xene-api-docs.netlify.app/apidocs.html)
- [Roadmap](/ROADMAP)
- [Xene Documentation](https://xene-api-docs.netlify.app/)
- [Go Doc](https://pkg.go.dev/github.com/fristonio/xene)
- [Getting Started](/docs/GettingStarted.md)

## License

Xene is licensed under [MIT License](https://github.com/fristonio/xene/blob/master/LICENSE.md).

## Contact

If you have any queries regarding the project or just want to say hello, feel free to drop a mail at deepshpathak@gmail.com.
