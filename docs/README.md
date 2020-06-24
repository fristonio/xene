# Xene

> A distributed workflow runner with a focus on performance and simplicity.

[![Netlify Status](https://api.netlify.com/api/v1/badges/f3adc406-ad04-4059-ad21-6a54f4be6771/deploy-status)](https://app.netlify.com/sites/sad-thompson-bcaa9a/deploys) ![Travis Status](https://travis-ci.com/fristonio/xene.svg?token=xvk2YsyqhEExfPszH3rV&branch=master) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


Xene is a high performance, distributed workflow execution platform. It lets you create custom workflows in the form of a Directed Acyclic Graph, which can then be executed based on a trigger configured with the workflow.

Xene is built with Golang with a focus on performance, resiliency, and simplicity. Xene internally runs in a way similar to Kubernetes by getting triggered on a level and then reconciling the desired state with the current state. Currently, xene supports [Badger](https://github.com/dgraph-io/badger) as the storage engines, which is pluggable and can be replaced with a similar key-value store.

## Documentation

- [Getting Started](/GettingStarted.md)
- [API Docs](https://xene-api-docs.netlify.app/apidocs.html)
- [Go Doc](https://pkg.go.dev/github.com/fristonio/xene)

## License

Xene is licensed under [MIT License](https://github.com/fristonio/xene/blob/master/LICENSE.md).
