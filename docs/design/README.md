# Xene Design

> A distributed workflow runner with focus on performance and simplicity.

Xene derives a lot of its functionalities/architectural decisions from Kubernetes. This is to make the system
flexible but at the same with Xene we try to reduce any resource footprint not necessirily required.

Xene is Level triggered just like kubernetes which means it runs a reconciliation loop for the Objects it care for
and tries to match the exixting state of the Object with the desired state. It uses a KVStore similar to K8S to store Objects and information
about workflows and pipelines.

Basic building blocks of xene are:

* Workflow
* Pipeline
* Task
* Step
* Trigger
* Connector(To be introduced in later version)

## Workflow

A workflow is the topmost resource that xene expects and manages, this workflow defines some basic needed information for xene to
start scheduling a worker for the workflow. Xene is smart enough and can schedule different component of this workflow on different
available workers to make the best use of existing resources.

A xene workflow can have multiple pipelines associated with it, which can also define some sort of internal execution
order based on the spec. A workflow must consist of atleast one pipeline for execution.

See [WorkflowSpec](/docs/design/spec/workflow.md) to see specification.

## Pipeline

A pipeline is a logical building block of xene. A pipeline consists of activities connected to each other in a logical
ordering(Connectors can be used between different activity/tasks to have more control over the input and output of the
pipeline). Each pipeline execution is handled by a `Trigger` associated with it.

A pipeline is scheduled on a single Agent node, which then sets up a controller for the trigger(based on the type of the trigger) and runs
the pipeline when the execution is triggered by the controller.

A Pipeline can have the following status associated with it on the APIServer side as the API server is only responsible
for the scheduling of the pipelin:

- Scheduled
- NotScheduled

The execution of the pipeline is in the hand of the Agent which can put the pipeline in any one of the
following state

- StartingUp
- Running
- Error
- Success
- ShuttingDown
- Unknown

## Task

A task is the basic execution unit managed by xene. Each task is associated with a context that helps xene to
group activities in the same environment. A single pipeline consists of the Directed Acyclic graph of these tasks
which is then resolved and executed by xene agent on the run.

A single task contains a sequential list of steps to be executed for the task. Steps in the same task shares the same
context. For example a downloaded artifact is mounted on the same container that the task is running on.

## Step

Step is the atomic unit of execution in xene. A step can be either of the following state

- Running
- Success
- Error
- NotExecuted

In a single task if any of the steps fails, it means that the task failed as a whole.

## Trigger

Each pipeline consists of a trigger which handles the job of starting the execution of pipeline. By default each pipeline can also
be triggered from the API Server using the connstructed webhook for the same.

Other customizable triggers for xene include

* Webhook Trigger
* Cron Trigger

## Connector

Connects different activites/tasks in a pipeline with each other in order to provide a logical ordering, these connectors also allow us
to build a sort of dependency graph for defining execution order.

* BasicConnector - Simple connector for pointing the execution pointer from one task to another.(Default)
* Rest other connectors are in development.

These connectors can help us couple inputs and outputs of different tasks and make more informed decisions about the
execution of the pipeline or the next step/task.
