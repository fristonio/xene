package docker

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerfilters "github.com/docker/docker/api/types/filters"
	dockerstrslice "github.com/docker/docker/api/types/strslice"
	dockerapi "github.com/docker/docker/client"
	"github.com/fristonio/xene/pkg/defaults"
	"github.com/fristonio/xene/pkg/executor/cre/runtime"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
)

var (
	dockerRuntimeName string = "docker"

	containerTypeLabelKey       = "io.xene.docker.type"
	containerTypeLabelContainer = "container"
	containerLogPathLabelKey    = "io.xene.container.logpath"

	conflictRE = regexp.MustCompile(`Conflict. (?:.)+ is already in use by container \"?([0-9a-z]+)\"?`)

	defaultDockerTimeout = time.Second * 10
)

// RuntimeExecutor contains the runtime executor for docker provider.
type RuntimeExecutor struct {
	client *dockerapi.Client

	execHandler ExecHandler

	timeout time.Duration
}

// NewCRE returns an instance of new docker container runtime executor.
func NewCRE() (*RuntimeExecutor, error) {
	// getDockerClient from the environment
	cli, err := getDockerClient()
	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %s", err)
	}
	return &RuntimeExecutor{
		client:      cli,
		timeout:     defaultDockerTimeout,
		execHandler: &NativeExecHandler{},
	}, nil
}

// Type returns the type of the container runtime executor, for docker
// this is docker.
func (e *RuntimeExecutor) Type() string {
	return string(v1alpha1.DockerExecutor)
}

// Version returns the version information about the cotainer runtime.
func (e *RuntimeExecutor) Version(ctx context.Context, vr *runtime.VersionRequest) (*runtime.VersionResponse, error) {
	v, err := e.client.ServerVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker version: %v", err)
	}
	// Docker API version (e.g., 1.23) is not semver compatible. Add a ".0"
	// suffix to remedy this.
	v.APIVersion = fmt.Sprintf("%s.0", v.APIVersion)

	return &runtime.VersionResponse{
		Version:           defaults.AgentVersion,
		RuntimeName:       dockerRuntimeName,
		RuntimeVersion:    v.Version,
		RuntimeAPIVersion: v.APIVersion,
	}, nil
}

// CreateContainer creates a new container
// Docker cannot store the log to an arbitrary location (yet), so we create an
// symlink at LogPath, linking to the actual path of the log.
func (e *RuntimeExecutor) CreateContainer(ctx context.Context,
	r *runtime.CreateContainerRequest) (*runtime.CreateContainerResponse, error) {
	config := r.Config

	if config == nil {
		return nil, fmt.Errorf("container config is nil")
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}
	// Apply a the container type label.
	config.Labels[containerTypeLabelKey] = containerTypeLabelContainer
	// Write the container log path in the labels.
	config.Labels[containerLogPathLabelKey] = filepath.Join(config.LogDirectory, config.LogPath)

	image := ""
	if config.Image != nil {
		image = config.Image.Image
	}
	containerName := config.Metadata.Name
	createConfig := dockertypes.ContainerCreateConfig{
		Name: containerName,
		Config: &dockercontainer.Config{
			// TODO: set User.
			Entrypoint: dockerstrslice.StrSlice(config.Command),
			Cmd:        dockerstrslice.StrSlice(config.Args),
			Env:        generateEnvList(config.Envs),
			Image:      image,
			WorkingDir: config.WorkingDir,
			Labels:     config.Labels,
			// Interactive containers:
			OpenStdin: config.Stdin,
			StdinOnce: config.StdinOnce,
			Tty:       config.Tty,
			// Disable Docker's health check until we officially support it
			// (https://github.com/kubernetes/kubernetes/issues/25829).
			// TODO: find a fix  for xene
			Healthcheck: &dockercontainer.HealthConfig{
				Test: []string{"NONE"},
			},
		},
		HostConfig: &dockercontainer.HostConfig{
			Binds: generateMountBindings(config.Mounts),
			RestartPolicy: dockercontainer.RestartPolicy{
				Name: "no",
			},
		},
	}

	// Set devices for container.
	devices := make([]dockercontainer.DeviceMapping, len(config.Devices))
	for i, device := range config.Devices {
		devices[i] = dockercontainer.DeviceMapping{
			PathOnHost:        device.HostPath,
			PathInContainer:   device.ContainerPath,
			CgroupPermissions: device.Permissions,
		}
	}
	createConfig.HostConfig.Resources.Devices = devices

	createResp, createErr := e.client.ContainerCreate(ctx,
		createConfig.Config, createConfig.HostConfig, createConfig.NetworkingConfig, createConfig.Name)
	if createErr != nil {
		createResp, createErr = recoverFromCreationConflictIfNeeded(e.client, createConfig, createErr)
	}

	if createResp.ID != "" {
		containerID := createResp.ID

		return &runtime.CreateContainerResponse{ContainerID: containerID}, nil
	}

	return nil, createErr
}

// StartContainer starts running the container in context.
func (e *RuntimeExecutor) StartContainer(ctx context.Context, r *runtime.StartContainerRequest) error {
	err := e.client.ContainerStart(ctx, r.ContainerID, dockertypes.ContainerStartOptions{})
	if ctxErr := contextError(ctx); ctxErr != nil {
		return ctxErr
	}

	return err
}

// StopContainer stops a running container.
func (e *RuntimeExecutor) StopContainer(ctx context.Context, r *runtime.StopContainerRequest) error {
	err := e.client.ContainerStop(ctx, r.ContainerID, &e.timeout)
	if ctxErr := contextError(ctx); ctxErr != nil {
		return ctxErr
	}

	return err
}

// RemoveContainer removes a container using the docker container
// runtime.
func (e *RuntimeExecutor) RemoveContainer(ctx context.Context, r *runtime.RemoveContainerRequest) error {
	err := e.client.ContainerRemove(ctx, r.ContainerID, dockertypes.ContainerRemoveOptions{
		Force: true,
	})
	if ctxErr := contextError(ctx); ctxErr != nil {
		return ctxErr
	}

	return err
}

// ListContainers list all the containers that match the required criteria
func (e *RuntimeExecutor) ListContainers(ctx context.Context,
	r *runtime.ListContainersRequest) (*runtime.ListContainersResponse, error) {
	filter := r.Filter
	opts := dockertypes.ContainerListOptions{All: true}

	opts.Filters = dockerfilters.NewArgs()
	f := newDockerFilter(&opts.Filters)
	// Add filter to get only xene containers
	f.AddLabel(containerTypeLabelKey, containerTypeLabelContainer)

	if filter != nil {
		if filter.ID != "" {
			f.Add("id", filter.ID)
		}
		if filter.State != nil {
			f.Add("status", toDockerContainerStatus(filter.State.State))
		}

		if filter.LabelSelector != nil {
			for k, v := range filter.LabelSelector {
				f.AddLabel(k, v)
			}
		}
	}
	containers, err := e.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Convert docker to runtime api containers.
	result := []*runtime.Container{}
	for i := range containers {
		c := containers[i]

		converted, err := toRuntimeContainer(&c)
		if err != nil {
			log.Infof("Unable to convert docker to runtime API container: %v", err)
			continue
		}

		result = append(result, converted)
	}

	return &runtime.ListContainersResponse{Containers: result}, nil
}

// ContainerStatus returns status of the container. If the container is not
// present, returns an error.
func (e *RuntimeExecutor) ContainerStatus(ctx context.Context,
	req *runtime.ContainerStatusRequest) (*runtime.ContainerStatusResponse, error) {
	containerID := req.ContainerID
	r, err := e.client.ContainerInspect(ctx, containerID)
	if ctxErr := contextError(ctx); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil {
		return nil, err
	}

	// Parse the timestamps.
	createdAt, startedAt, finishedAt, err := getContainerTimestamps(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp for container %q: %v", containerID, err)
	}

	// Convert the image id to a pullable id.
	resp, _, err := e.client.ImageInspectWithRaw(ctx, r.Image)
	if ctxErr := contextError(ctx); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil {
		if dockerapi.IsErrNotFound(err) {
			err = imageNotFoundError{ID: r.Image}
		}
		return nil, err
	}

	if !matchImageIDOnly(resp, r.Image) {
		return nil, imageNotFoundError{ID: r.Image}
	}

	if err != nil {
		if !isImageNotFoundError(err) {
			return nil, fmt.Errorf("unable to inspect docker image %q while inspecting docker container %q: %v",
				r.Image, containerID, err)
		}
		log.Warningf("ignore error image %q not found while inspecting docker container %q: %v",
			r.Image, containerID, err)
	}
	imageID := toPullableImageID(r.Image, &resp)

	// Convert the mounts.
	mounts := make([]*runtime.Mount, 0, len(r.Mounts))
	for i := range r.Mounts {
		m := r.Mounts[i]
		readonly := !m.RW
		mounts = append(mounts, &runtime.Mount{
			HostPath:      m.Source,
			ContainerPath: m.Destination,
			Readonly:      readonly,
			// Note: Can't set SeLinuxRelabel
		})
	}
	// Interpret container states.
	var state runtime.ContainerState
	var reason, message string
	if r.State.Running {
		// Container is running.
		state = runtime.ContainerStateRunning
	} else {
		// Container is *not* running. We need to get more details.
		//    * Case 1: container has run and exited with non-zero finishedAt
		//              time.
		//    * Case 2: container has failed to start; it has a zero finishedAt
		//              time, but a non-zero exit code.
		//    * Case 3: container has been created, but not started (yet).
		if !finishedAt.IsZero() { // Case 1
			state = runtime.ContainerStateExited
			switch {
			case r.State.OOMKilled:
				// TODO: consider exposing OOMKilled via the runtime.
				// Note: if an application handles OOMKilled gracefully, the
				// exit code could be zero.
				reason = "OOMKilled"
			case r.State.ExitCode == 0:
				reason = "Completed"
			default:
				reason = "Error"
			}
		} else if r.State.ExitCode != 0 { // Case 2
			state = runtime.ContainerStateExited
			// Adjust finshedAt and startedAt time to createdAt time to avoid
			// the confusion.
			finishedAt, startedAt = createdAt, createdAt
			reason = "ContainerCannotRun"
		} else { // Case 3
			state = runtime.ContainerStateCreated
		}
		message = r.State.Error
	}

	// Convert to unix timestamps.
	ct, st, ft := createdAt.UnixNano(), startedAt.UnixNano(), finishedAt.UnixNano()
	exitCode := int32(r.State.ExitCode)

	imageName := r.Config.Image
	if resp.ID == "" && len(resp.RepoTags) > 0 {
		imageName = resp.RepoTags[0]
	}
	status := &runtime.ContainerStatus{
		ID: r.ID,
		Metadata: &runtime.ContainerMetadata{
			Name: r.Name,
		},
		Image:      &runtime.ImageSpec{Image: imageName},
		ImageRef:   imageID,
		Mounts:     mounts,
		ExitCode:   exitCode,
		State:      state,
		CreatedAt:  ct,
		StartedAt:  st,
		FinishedAt: ft,
		Reason:     reason,
		Message:    message,
		Labels:     r.Config.Labels,
		LogPath:    r.Config.Labels[containerLogPathLabelKey],
	}
	return &runtime.ContainerStatusResponse{Status: status}, nil
}

// ExecSync executes the command in the docker container
func (e *RuntimeExecutor) ExecSync(ctx context.Context, r *runtime.ExecRequest) (*runtime.ExecResponse, error) {
	cj, err := e.checkContainerStatus(r.ContainerID)
	if err != nil {
		return nil, err
	}

	res, err := e.execHandler.ExecInContainer(ctx, e.client, cj, r.Cmd, r.Stdin, r.Stdout, r.Stdout, r.Timeout, r.Tty)
	if err != nil {
		return nil, err
	}

	return &runtime.ExecResponse{
		ExitCode:    res.ExitCode,
		ContainerID: res.ContainerID,
	}, nil
}

func (e *RuntimeExecutor) checkContainerStatus(containerID string) (*dockertypes.ContainerJSON, error) {
	container, err := e.client.ContainerInspect(context.TODO(), containerID)
	if err != nil {
		return nil, err
	}
	if !container.State.Running {
		return nil, fmt.Errorf("container not running (%s)", container.ID)
	}
	return &container, nil
}

// Status returns the status of the underlying runtime.
func (e *RuntimeExecutor) Status(ctx context.Context, r *runtime.StatusRequest) (*runtime.StatusResponse, error) {
	runtimeReady := &runtime.RuntimeCondition{
		Type:   runtime.RuntimeReady,
		Status: true,
	}
	networkReady := &runtime.RuntimeCondition{
		Type:   runtime.NetworkReady,
		Status: true,
	}
	conditions := []*runtime.RuntimeCondition{runtimeReady, networkReady}
	if _, err := e.client.ServerVersion(ctx); err != nil {
		runtimeReady.Status = false
		runtimeReady.Reason = "DockerDaemonNotReady"
		runtimeReady.Message = fmt.Sprintf("docker: failed to get docker version: %v", err)
	}

	status := &runtime.StatusResponse{Conditions: conditions}
	return status, nil
}
