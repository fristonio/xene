# API Server

API server is the centeral component of xene, which can be though of as the master in a master-slave
architecture. There are two main functions of apiserver:

* Act as a API interface for the users/consumers to interact with.
* Scheduling jobs on agents running remotely and aggregating results.

Apart from this API server is also responsible for various small tasks which helps in assisting the
required state of the application. This involves interacting with kvstore to store/get data along with
managing health status of agents and scheduling tasks.

API Server is also responsible for issuing tokens for agents through which the nodes can join in the
worker cluster. The central apiserver also issue regular health probes for the associated agents and ensure
that tasks are scheduled only on a healthy agent. The communication between the agent and the API server is
done using GRPC protocol, this communication includes:

* Issue of a task to an agent
* Aggregating results from the tasks performed by the worker agent.

API Server also issues the required secrets to the agent based on the label that are associated with it.

