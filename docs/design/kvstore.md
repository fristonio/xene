# KVStore

> xene internally uses a KVstore to handle storage of created resources.

There are three major resources that xene stores in its local key value store, these are:

* Workflow
    * A workflow consists of multiple pipeline under a single context. The workflow can be thought of a logical
    grouping of the pipelines.
* Pipeline
    * A pipeline is a single execution unit consisting of multiple activities. Each pipeline consists of activities
    connected to each other in some fashion using some sort of Connectors. Each pipeline also consist of a trigger which
    in a dependency graph of activities can be seen as the root node.
* Task
    * Task is the most basic unit of execution in xene. A single task consist of many steps which executes under
    the same context. Each invidual type of task has its own internal structure in terms of how it handles various
    steps and which executor it uses.

In the key value store these units are stored using a predefined format of key as described below:

* Workflow: `/registry/workflows/<namespace>/<workflow-id>`
* Pipeline: `/registry/pipelines/<namespace>/<pipeline-id>`
* Task: `/registry/activities/<namespace>/<task-id>`

Apart from these public keys some internally used keys are:

* Namespaces: `/registry/namespaces/<namespace-name>`
* Agents: `/registry/agents/<agent-name>`
* Executors: `/registry/executors/<executor-name>`

Each key in the store also have some metadata associated with it which includes:

