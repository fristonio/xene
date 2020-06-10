package cre

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/executor/cre/docker"
	"github.com/fristonio/xene/pkg/executor/cre/runtime"
	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/templates"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"github.com/fristonio/xene/pkg/utils"
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

	status *v1alpha1.PipelineRunStatus

	imageRef string

	// log contains the logger for the pipeline executor
	log *logrus.Entry

	// useStore specifies wheather to use kvstore interaction during the execution
	// of the pipeline.
	useStore bool

	mux *sync.Mutex
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
		spec:     spec,
		useStore: true,
		mux:      &sync.Mutex{},
	}
}

func (c *CRExecutor) getResName() string {
	return fmt.Sprintf("%s-%s", c.name, c.id)
}

// WithoutStore returns the CRExecutor disabling useStore option
func (c *CRExecutor) WithoutStore() {
	c.useStore = false
}

// WithStatus sets the status in the executor
func (c *CRExecutor) WithStatus(status *v1alpha1.PipelineRunStatus) {
	c.status = status
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

	if !utils.FileExists(defaults.AgentMountScript) {
		if !utils.DirExists(defaults.AgentAssetsDir) {
			if err := os.MkdirAll(defaults.AgentAssetsDir, os.ModePerm); err != nil {
				return fmt.Errorf("error while creating mount script: %s", err)
			}
		}

		data := templates.GetAgentMountScript()
		err := ioutil.WriteFile(defaults.AgentMountScript, []byte(data), 0777)
		if err != nil {
			return fmt.Errorf("error while writing agent script: %s", err)
		}
	}

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

	envs := make([]*runtime.KeyValue, 0)
	for key, val := range c.spec.Envs {
		envs = append(envs, &runtime.KeyValue{
			Key:   key,
			Value: val,
		})
	}
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
			Command:    []string{"sleep", fmt.Sprintf("%d", int64(defaults.GlobalPipelineTimeout/time.Second))},
			WorkingDir: "/",
			Mounts: []*runtime.Mount{
				{
					ContainerPath: defaults.AgentMountContainerScript,
					HostPath:      defaults.AgentMountScript,
				},
			},
			Envs: envs,
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

	c.mux.Lock()
	c.status.Tasks[name].Status = v1alpha1.StatusRunning
	err := c.SaveStatusToStore()
	if err != nil {
		c.log.Errorf("error while saving status to the store: %s", err)
	}
	c.mux.Unlock()

	for _, step := range task.Steps {
		c.log.Infof("Running step: %s", step.Name)

		l := newLogger(c.name, c.id, name, step.Name)
		w := l.getLogWriter()

		c.mux.Lock()
		c.status.Tasks[name].Steps[step.Name].Status = v1alpha1.StatusRunning
		c.status.Tasks[name].Steps[step.Name].LogFile = l.getLogFileName()
		err := c.SaveStatusToStore()
		if err != nil {
			c.log.Errorf("error while saving status to the store: %s", err)
		}
		c.mux.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
		defer cancel()
		start := time.Now()
		res, err := c.cre.ExecSync(ctx, &runtime.ExecRequest{
			ContainerID: c.getResName(),
			Cmd:         []string{"xene-cmd-run.sh", task.WorkingDirectory, step.Cmd},
			Tty:         false,
			Stdin:       nil,
			Stdout:      w,
			Stderr:      w,
			Timeout:     time.Minute * 30,
		})

		if res != nil {
			c.log.Debugf("exec info: %s: %d", res.ContainerID, res.ExitCode)
		}

		if err != nil || res == nil {
			c.mux.Lock()
			c.status.Tasks[name].Steps[step.Name] = &v1alpha1.StepRunStatus{
				Status:  v1alpha1.StatusError,
				LogFile: l.getLogFileName(),
				Time:    time.Since(start),
			}
			c.status.Tasks[name].Status = v1alpha1.StatusError
			err := c.SaveStatusToStore()
			if err != nil {
				c.log.Errorf("error while saving status to the store: %s", err)
			}
			c.mux.Unlock()
			if w != nil {
				w.Close()
			}

			return err
		}

		if res.ExitCode != 0 {
			c.mux.Lock()
			c.status.Tasks[name].Steps[step.Name] = &v1alpha1.StepRunStatus{
				Status:  v1alpha1.StatusError,
				LogFile: l.getLogFileName(),
				Time:    time.Since(start),
			}
			c.status.Tasks[name].Status = v1alpha1.StatusError
			err := c.SaveStatusToStore()
			if err != nil {
				c.log.Errorf("error while saving status to the store: %s", err)
			}
			c.mux.Unlock()

			return &StepExecError{
				name: step.Name,
			}
		}

		c.mux.Lock()
		c.status.Tasks[name].Steps[step.Name] = &v1alpha1.StepRunStatus{
			Status:  v1alpha1.StatusSuccess,
			LogFile: l.getLogFileName(),
			Time:    time.Since(start),
		}
		err = c.SaveStatusToStore()
		if err != nil {
			c.log.Errorf("error while saving status to the store: %s", err)
		}
		c.mux.Unlock()

		if w != nil {
			_ = w.Close()
		}
		c.log.Infof("Step completed: %s", step.Name)
	}

	c.mux.Lock()
	c.status.Tasks[name].Status = v1alpha1.StatusSuccess
	err = c.SaveStatusToStore()
	if err != nil {
		c.log.Errorf("error while saving status to the store: %s", err)
	}
	c.mux.Unlock()

	return nil
}

// Shutdown shuts down the container runtime executor.
func (c *CRExecutor) Shutdown() error {
	c.log.Infof("shutting down the container runtime executor")

	ctx, cancel := context.WithTimeout(context.Background(), defaults.CreateContainerTimeout)
	defer cancel()
	err := c.cre.StopContainer(ctx, &runtime.StopContainerRequest{
		ContainerID: c.getResName(),
		Timeout:     int64(defaults.CreateContainerTimeout / time.Second),
	})
	if err != nil {
		c.log.Warnf("error while stopping container: %s", err)
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

// SaveStatusToStore saves the status in the store
func (c *CRExecutor) SaveStatusToStore() error {
	if !c.useStore {
		return nil
	}

	statusKey := fmt.Sprintf("%s/%s/%s", v1alpha1.PipelineStatusKeyPrefix, c.name, c.id)

	val, err := json.Marshal(c.status)
	if err != nil {
		return fmt.Errorf("error while marshalling status: %s", err)
	}

	err = store.KVStore.Set(context.TODO(), statusKey, val)
	if err != nil {
		return fmt.Errorf("error while setting status key in kvstore: %s", err)
	}

	return nil
}
