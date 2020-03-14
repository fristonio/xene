# Workflow spec

What is the reasoning behind creating a workflow type above Pipeline, when there is no logical
connection between the two?

A Workflow in xene let's us keep all the related pipelines at one place, which also helps to share a lot of
functionality like sharing triggers for pipelines.

Should triggers be for a workflow or for a single pipeline?

Trigger is associated with a Workflow and not a pipeline, the reasoning is that there are going to be a very low number
of cases when we have unrelated actions to perform which are logically associated. This also supports to the point of
having a higher level object over pipeline, which will contain common properties of the pipelines that are logically
associated.

Example Workflow:

CI Job for a particular application triggered on every commit on master.

* Pipeline to run Build - Test - Push for the application.
* Pipeline to sync the git repository with multiple platforms like - Github, Gitlab etc.
``````````i
```yaml
apiVersion: v1alpha1
kind: workflow
name: ci-workflow
description: CI Workflow using xene

metadata:
  creationTimestamp:
  deletionTimestamp:
  labels:
    type: testing-ci

spec:
  triggers:
  - name: master-commit
    description: Triggered using an HTTP request to an auto configured webhook endpoint for the workflow.
    type: webhook

  pipelines:
  - name: test-pipeline
    description: Run CI pipeline for the changes landed in the master.
    tasks:
    - name: build
    - name: test
      dependencies: [build]
    - name: push
      dependencies: [test]

  - name: git-sync
    description: Sync the commit in master to all the git mirrors maintained.
    tasks:
    - name: github-sync
    - name: gitlab-sync
```

