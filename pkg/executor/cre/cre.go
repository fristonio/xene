package cre

import (
	"context"
	"fmt"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/executor/cre/docker"
	"github.com/fristonio/xene/pkg/executor/cre/runtime"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/sirupsen/logrus"
)

// CRE is the interface corresponding to the container runtime executor.
type CRE interface {
	// Type returns the container runtime executor type, by default docker
	// is used.
	Type() string

	runtime.RuntimeService

	runtime.ImageService
}

// CRExecutor is the type corresponding to the container runtime executor.
type CRExecutor struct {
	// runtime is the underlying container runtime to use for the
	// ContainerRuntimeExecutor
	runtime string

	// name for the container runtime executor
	// all the resources created by the executor will be related to
	// this name.
	name string

	// cre is the container runtime executor for the Runtime
	// executor
	cre CRE

	// id is the ID corresponding to the current run, it should be unique
	// for each of the CRExecutor instance.
	id string

	spec *v1alpha1.PipelineSpecWithName

	imageRef string

	// log contains the logger for the pipeline executor
	log *logrus.Entry
}

// NewCRExecutor returns a new instance of the container runtime
// executor.
func NewCRExecutor(runtime, id, name string, spec *v1alpha1.PipelineSpecWithName) *CRExecutor {
	return &CRExecutor{
		runtime: runtime,
		name:    name,
		id:      id,
		log: logrus.WithFields(logrus.Fields{
			"pipeline": name,
			"id":       id,
			"runtime":  runtime,
		}),
		spec: spec,
	}
}

func (c *CRExecutor) getResName() string {
	return fmt.Sprintf("%s-%s", c.name, c.id)
}

// Configure configures the container runtime executor
func (c *CRExecutor) Configure() error {
	c.log.Infof("Configuring the container runtime executor")
	var (
		cre CRE
		err error
	)
	switch v1alpha1.Executor(c.runtime) {
	case v1alpha1.DockerExecutor:
		cre, err = docker.NewCRE()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("not a valid Executor: %s", c.runtime)
	}

	c.cre = cre

	img := parseImageCanonicalURL(c.spec.Executor.ContainerConfig.Image)

	// First set up the image for the container to run.
	// Here we assume that the name is unique for each pipeline run
	ctx, cancel := context.WithTimeout(context.Background(), defaults.ImagePullDeadline)
	defer cancel()
	res, err := cre.PullImage(ctx, &runtime.PullImageRequest{
		Image: &runtime.ImageSpec{
			Image: img,
		},
	})
	if err != nil {
		return fmt.Errorf("error while pulling pipeline image: %s", err)
	}

	c.imageRef = res.ImageRef

	// After pulling in the image for the container
	// create the container with the configuration required.
	ctx, cancel = context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
	defer cancel()
	resp, err := cre.CreateContainer(ctx, &runtime.CreateContainerRequest{
		Config: &runtime.ContainerConfig{
			Metadata: &runtime.ContainerMetadata{
				Name: c.getResName(),
			},
			Image: &runtime.ImageSpec{
				Image: res.ImageRef,
			},
			Command:    []string{"sleep", "100000"},
			WorkingDir: "/",
		},
	})
	if err != nil {
		return fmt.Errorf("error while creating container: %s", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
	defer cancel()
	err = cre.StartContainer(ctx, &runtime.StartContainerRequest{
		ContainerID: resp.ContainerID,
	})
	if err != nil {
		return fmt.Errorf("error while starting container: %s", err)
	}
	return nil
}

// StepExecError is the error in the step of the task, it is returned
// when the step failed exitting in a code other than 0
type StepExecError struct {
	name string
}

func (s *StepExecError) Error() string {
	return fmt.Sprintf("Failed step: %s", s.name)
}

// RunTask runs the task defined by the spec in the container.
func (c *CRExecutor) RunTask(name string, task *v1alpha1.TaskSpec) error {
	c.log.Infof("Running task: %s", name)

	for _, step := range task.Steps {
		c.log.Infof("Running step: %s", step.Name)
		l := newLogger(c.name, c.id, name, step.Name)
		w := l.getLogWriter()

		ctx, cancel := context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
		defer cancel()
		res, err := c.cre.ExecSync(ctx, &runtime.ExecRequest{
			ContainerID: c.getResName(),
			Cmd:         step.Cmd,
			Tty:         false,
			Stdin:       nil,
			Stdout:      w,
			Stderr:      w,
		})
		if err != nil {
			if w != nil {
				w.Close()
			}
			return err
		}

		if res.ExitCode != 0 {
			return &StepExecError{
				name: step.Name,
			}
		}

		if w != nil {
			_ = w.Close()
		}
		c.log.Infof("Step completed: %s", step.Name)
	}

	return nil
}

// Shutdown shuts down the container runtime executor.
func (c *CRExecutor) Shutdown() error {
	c.log.Infof("shutting down the container runtime executor")
	ctx, cancel := context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
	defer cancel()
	err := c.cre.StopContainer(ctx, &runtime.StopContainerRequest{
		ContainerID: c.getResName(),
	})
	if err != nil {
		return fmt.Errorf("error while stopping container: %s", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
	defer cancel()
	err = c.cre.RemoveContainer(ctx, &runtime.RemoveContainerRequest{
		ContainerID: c.getResName(),
	})
	if err != nil {
		return fmt.Errorf("error while removing container: %s", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
	defer cancel()
	err = c.cre.RemoveImage(ctx, &runtime.RemoveImageRequest{
		Image: &runtime.ImageSpec{
			Image: c.imageRef,
		},
	})
	if err != nil {
		return fmt.Errorf("error while removing image: %s", err)
	}

	return nil
}
