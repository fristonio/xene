package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockerapi "github.com/docker/docker/client"
)

// ExecHandler knows how to execute a command in a running Docker container.
type ExecHandler interface {
	ExecInContainer(context.Context, *dockerapi.Client, *dockertypes.ContainerJSON,
		[]string, io.Reader, io.WriteCloser, io.WriteCloser, time.Duration, bool) (*dockertypes.ContainerExecInspect, error)
}

type dockerExitError struct {
	Inspect *dockertypes.ContainerExecInspect
}

func (d *dockerExitError) String() string {
	return d.Error()
}

func (d *dockerExitError) Error() string {
	return fmt.Sprintf("Error executing in Docker Container: %d", d.Inspect.ExitCode)
}

func (d *dockerExitError) Exited() bool {
	return !d.Inspect.Running
}

func (d *dockerExitError) ExitStatus() int {
	return d.Inspect.ExitCode
}

// ExecTimeoutError is the error corresponding to the timeout in
// exec of the step.
type ExecTimeoutError struct{}

func (d *ExecTimeoutError) Error() string {
	return fmt.Sprintf("Exec deadline excedded, timeout")
}

// NativeExecHandler executes commands in Docker containers using Docker's exec API.
type NativeExecHandler struct{}

// ExecInContainer executes the provided command in the container using docker client
func (n *NativeExecHandler) ExecInContainer(ctx context.Context, client *dockerapi.Client,
	container *dockertypes.ContainerJSON, cmd []string,
	stdin io.Reader, stdout io.WriteCloser, stderr io.WriteCloser,
	timeout time.Duration, tty bool) (*dockertypes.ContainerExecInspect, error) {
	done := make(chan struct{})
	defer close(done)

	createOpts := dockertypes.ExecConfig{
		Cmd:          cmd,
		AttachStdin:  stdin != nil,
		AttachStdout: stdout != nil,
		AttachStderr: stderr != nil,
		Tty:          tty,
	}

	execObj, err := client.ContainerExecCreate(context.TODO(), container.ID, createOpts)
	if ctxErr := contextError(ctx); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil {
		return nil, fmt.Errorf("failed to exec in container - Exec setup failed - %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.ContainerExecAttach(ctx, execObj.ID, createOpts)
	if ctxErr := contextError(ctx); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil {
		return nil, err
	}

	err = holdHijackedConnection(false, stdin, stdout, stderr, resp)
	if err != nil {
		log.Errorf("error while hijacking exec connection: %s", err)
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	var curDur time.Duration
	for {
		inspect, err2 := client.ContainerExecInspect(context.TODO(), execObj.ID)
		if err2 != nil {
			return nil, fmt.Errorf("fialed while inspecting exec container: %s", err2)
		}

		if !inspect.Running {
			if inspect.ExitCode != 0 {
				err = &dockerExitError{&inspect}
			}

			return &inspect, err
		}

		curDur = curDur + time.Second*2
		if curDur > (timeout + time.Second*2) {
			log.Errorf("Exec session %s in container %s terminated but process still running!", execObj.ID, container.ID)
			return nil, &ExecTimeoutError{}
		}

		<-ticker.C

	}
}
