# Xene Design

Basic building blocks of xene workflow are:

* Workflow
* Pipeline
* Task
* Connectors
* Triggers

## Workflow

A workflow is the topmost resource that xene expects, this workflow defines some basic needed information for xene to
start schedule a worker for the workflow. Xene is smart enough and can schedule different component on different
available workers to make the best use of existing resources.

A xene workflow can have multiple pipelines associated with it, which can also define some sort of internal execution
order based on the spec. A workflow must consist of atleast one pipeline for execution.

See [WorkflowSpec](/docs/spec/workflow.md) to see specification.

## Pipeline

A pipeline is a logical building block of xene. A pipeline consists of activities connected to each other in a logical
ordering defined the connecters between activities. Each pipeline execution is handled by a `Trigger` associated with
it.

## Task

A task is the basic execution unit managed by xene. Each task is associated with a context which helps xene to
group activities in the same environment. If no context is specified then xene assigns it a randomly generated isolated
context. Activites under same context share the same execution environment and thus can affect the execution of each
other.

## Connectors

Connects different activites in a pipeline with each other provide a logical ordering, these connectors also allow us
to build a sort of dependency graph for defining execution order.

* BasicConnector - Simple connector for pointing the execution pointer from one task to another.

## Triggers

Each pipeline consists of a trigger which handles the job of starting the execution of pipeline.

* Webhook Trigger
* PipelineTrigger

