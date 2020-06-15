package runtime

import (
	"context"
)

// RuntimeService defines the public APIs for remote container runtimes
type RuntimeService interface { //nolint
	// Version returns the runtime name, runtime version, and runtime API version.
	Version(context.Context, *VersionRequest) (*VersionResponse, error)

	// CreateContainer creates a new container in specified PodSandbox
	CreateContainer(context.Context, *CreateContainerRequest) (*CreateContainerResponse, error)

	// StartContainer starts the container.
	StartContainer(context.Context, *StartContainerRequest) error

	// StopContainer stops a running container with a grace period (i.e., timeout).
	// This call is idempotent, and must not return an error if the container has
	// already been stopped.
	StopContainer(context.Context, *StopContainerRequest) error

	// RemoveContainer removes the container. If the container is running, the
	// container must be forcibly removed.
	// This call is idempotent, and must not return an error if the container has
	// already been removed.
	RemoveContainer(context.Context, *RemoveContainerRequest) error

	// ListContainers lists all containers by filters.
	ListContainers(context.Context, *ListContainersRequest) (*ListContainersResponse, error)

	// ContainerStatus returns status of the container. If the container is not
	// present, returns an error.
	ContainerStatus(context.Context, *ContainerStatusRequest) (*ContainerStatusResponse, error)

	// TODO: Complete this
	// UpdateContainerResources updates ContainerConfig of the container.
	// UpdateContainerResources(context.Context, *UpdateContainerResourcesRequest) error

	// TODO: Complete this
	// ReopenContainerLog asks runtime to reopen the stdout/stderr log file
	// for the container. This is often called after the log file has been
	// rotated. If the container is not running, container runtime can choose
	// to either create a new log file and return nil, or return an error.
	// Once it returns error, new container log file MUST NOT be created.
	// ReopenContainerLog(context.Context, *ReopenContainerLogRequest) error

	// ExecSync runs a command in a container synchronously.
	ExecSync(context.Context, *ExecRequest) (*ExecResponse, error)

	// TODO: Implement ExecSync
	// Exec prepares a streaming endpoint to execute a command in the container.
	// Exec(context.Context, *ExecRequest) error

	// Status returns the status of the runtime.
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
}

// ImageService defines the public APIs for managing images.
type ImageService interface {
	// ListImages lists existing images.
	ListImages(context.Context, *ListImagesRequest) (*ListImagesResponse, error)

	// ImageStatus returns the status of the image. If the image is not
	// present, returns a response with ImageStatusResponse.Image set to
	// nil.
	ImageStatus(context.Context, *ImageStatusRequest) (*ImageStatusResponse, error)

	// PullImage pulls an image with authentication config.
	PullImage(context.Context, *PullImageRequest, string) (*PullImageResponse, error)

	// RemoveImage removes the image.
	// This call is idempotent, and must not return an error if the image has
	// already been removed.
	RemoveImage(context.Context, *RemoveImageRequest) error
}
