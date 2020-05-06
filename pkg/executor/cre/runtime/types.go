// Package runtime is a derived form of the container runtime api provided
// by kubernetes removing all the parts related to the Pod sandboxes.
package runtime

import (
	"io"
	"time"
)

// VersionRequest contains the request of the version of the container runtime
// executor.
type VersionRequest struct {
	// Version of the agent runtime API.
	Version string `json:"version,omitempty"`
}

// VersionResponse contains the response to the version request from the container runtime
// API.
type VersionResponse struct {
	// Version of the agent runtime API.
	Version string `json:"version,omitempty"`

	// RuntimeName contains the name of the container runtime.
	RuntimeName string `json:"runtimeName,omitempty"`

	// Version of the container runtime. The string must be
	// semver-compatible.
	RuntimeVersion string `json:"runtimeVersion,omitempty"`

	// RuntimeAPIVersion is the API version of the container runtime. The string must be
	// semver-compatible.
	RuntimeAPIVersion string `json:"runtimeApiVersion,omitempty"`
}

// CreateContainerRequest contains the reqeust to create a container from the config provided.
type CreateContainerRequest struct {
	// Config of the container.
	Config *ContainerConfig `json:"config,omitempty"`
}

// ContainerConfig holds all the required and optional fields for creating a
// container.
type ContainerConfig struct {
	// Metadata of the container. This information will uniquely identify the
	// container, and the runtime should leverage this to ensure correct
	// operation. The runtime may also use this information to improve UX, such
	// as by constructing a readable name.
	Metadata *ContainerMetadata `json:"metadata,omitempty"`

	// Image to use for the container.
	Image *ImageSpec `json:"image,omitempty"`

	// Command to execute (i.e., entrypoint for docker)
	Command []string `json:"command,omitempty"`

	// Args for the Command (i.e., command for docker)
	Args []string `json:"args,omitempty"`

	// Current working directory of the command.
	WorkingDir string `json:"workingDir,omitempty"`

	// List of environment variable to set in the container.
	Envs []*KeyValue `json:"envs,omitempty"`

	// Mounts for the container.
	Mounts []*Mount `json:"mounts,omitempty"`

	// Devices for the container.
	Devices []*Device `json:"devices,omitempty"`

	// Key-value pairs that may be used to scope and select individual resources.
	Labels map[string]string `json:"labels,omitempty"`

	Stdin     bool `json:"stdin,omitempty"`
	StdinOnce bool `json:"stdinOnce,omitempty"`
	Tty       bool `json:"tty,omitempty"`

	// Configuration specific to Linux containers.
	Linux *LinuxContainerConfig `json:"linux,omitempty"`

	// LogDirectory contains the directory to store the logs into
	LogDirectory string `json:"logDirectory"`

	LogPath string `json:"logPath"`
}

// ContainerMetadata holds all necessary information for building the container
// name. The container runtime is encouraged to expose the metadata in its user
// interface for better user experience. E.g., runtime can construct a unique
// container name based on the metadata. Note that (name, attempt) is unique
// within a sandbox for the entire lifetime of the sandbox.
type ContainerMetadata struct {
	// Name of the container. Same as the container name.
	Name string `json:"name,omitempty"`
}

// ImageSpec is an internal representation of an image.  Currently, it wraps the
// value of a Container's Image field (e.g. imageID or imageDigest), but in the
// future it will include more detailed information about the different image types.
type ImageSpec struct {
	Image string `json:"image,omitempty"`
}

// KeyValue contains a key value pair of strings.
type KeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Mount specifies a host volume to mount into a container.
type Mount struct {
	// Path of the mount within the container.
	ContainerPath string `json:"containerPath,omitempty"`

	// Path of the mount on the host. If the hostPath doesn't exist, then runtimes
	// should report error. If the hostpath is a symbolic link, runtimes should
	// follow the symlink and mount the real destination to container.
	HostPath string `json:"hostPath,omitempty"`

	// If set, the mount is read-only.
	Readonly bool `json:"readonly,omitempty"`

	// Requested propagation mode.
	Propagation int32 `json:"propagation,omitempty"`
}

// MountPropagation defines the kind of mount propogation to do for the container.
type MountPropagation int32

const (
	// MountPropagationPrivate means no mount propagation ("private" in Linux terminology).
	MountPropagationPrivate MountPropagation = 0

	// MountPropagationHostToContainer means mounts get propagated
	// from the host to the container ("rslave" in Linux).
	MountPropagationHostToContainer MountPropagation = 1

	// MountPropagationBiDirectional means mounts get propagated from the host to the container and from the
	// container to the host ("rshared" in Linux).
	MountPropagationBiDirectional MountPropagation = 2
)

// Device specifies a host device to mount into a container.
type Device struct {
	// Path of the device within the container.
	ContainerPath string `json:"containerPath,omitempty"`

	// Path of the device on the host.
	HostPath string `json:"hostPath,omitempty"`

	// Cgroups permissions of the device, candidates are one or more of
	// * r - allows container to read from the specified device.
	// * w - allows container to write to the specified device.
	// * m - allows container to create device files that do not yet exist.
	Permissions string `json:"permissions,omitempty"`
}

// LinuxContainerConfig contains platform-specific configuration for
// Linux-based containers.
type LinuxContainerConfig struct {
	// Resources specification for the container.
	Resources *LinuxContainerResources `json:"resources,omitempty"`
}

// LinuxContainerResources specifies Linux specific configuration for
// resources.
type LinuxContainerResources struct {
	// CPU CFS (Completely Fair Scheduler) period. Default: 0 (not specified).
	CPUPeriod int64 `json:"cpuPeriod,omitempty"`
	// CPU CFS (Completely Fair Scheduler) quota. Default: 0 (not specified).
	CPUQuota int64 `json:"cpuQuota,omitempty"`
	// CPU shares (relative weight vs. other containers). Default: 0 (not specified).
	CPUShares int64 `json:"cpuShares,omitempty"`
	// Memory limit in bytes. Default: 0 (not specified).
	MemoryLimitInBytes int64 `json:"memoryLimitInBytes,omitempty"`
	// OOMScoreAdj adjusts the oom-killer score. Default: 0 (not specified).
	OomScoreAdj int64 `json:"oomScoreAdj,omitempty"`
	// CpusetCpus constrains the allowed set of logical CPUs. Default: "" (not specified).
	CpusetCpus string `json:"cpusetCpus,omitempty"`
	// CpusetMems constrains the allowed set of memory nodes. Default: "" (not specified).
	CpusetMems string `json:"cpusetMems,omitempty"`
	// List of HugepageLimits to limit the HugeTLB usage of container per page size. Default: nil (not specified).
	HugepageLimits []*HugepageLimit `json:"hugepageLimits,omitempty"`
}

// HugepageLimit corresponds to the file`hugetlb.<hugepagesize>.limit_in_byte` in container level cgroup.
// For example, `PageSize=1GB`, `Limit=1073741824` means setting `1073741824` bytes to hugetlb.1GB.limit_in_bytes.
type HugepageLimit struct {
	// The value of PageSize has the format <size><unit-prefix>B (2MB, 1GB),
	// and must match the <hugepagesize> of the corresponding control
	// file found in `hugetlb.<hugepagesize>.limit_in_bytes`.
	// The values of <unit-prefix> are intended to be parsed
	// using base 1024("1KB" = 1024, "1MB" = 1048576, etc).
	PageSize string `json:"pageSize,omitempty"`
	// limit in bytes of hugepagesize HugeTLB usage.
	Limit uint64 `json:"limit,omitempty"`
}

// CreateContainerResponse contains the response corresponding to create container
// request.
type CreateContainerResponse struct {
	// ID of the created container.
	ContainerID string `json:"containerID,omitempty"`
}

// StartContainerRequest contains the information about the containerk which needs to be restarted.
type StartContainerRequest struct {
	// ID of the container to start.
	ContainerID string `json:"containerID,omitempty"`
}

// StopContainerRequest contains the information about the container which
// needs to be stopped.
type StopContainerRequest struct {
	// ID of the container to stop.
	ContainerID string `json:"containerID,omitempty"`
	// Timeout in seconds to wait for the container to stop before forcibly
	// terminating it. Default: 0 (forcibly terminate the container immediately)
	Timeout int64 `json:"timeout,omitempty"`
}

// RemoveContainerRequest contains the information about the container which
// needs to be removed.
type RemoveContainerRequest struct {
	// ID of the container to remove.
	ContainerID string `json:"containerID,omitempty"`
}

// ListContainersRequest contains the information about the container which
// needs to be listed.
type ListContainersRequest struct {
	Filter *ContainerFilter `json:"filter,omitempty"`
}

// ContainerFilter is used to filter containers.
// All those fields are combined with 'AND'
type ContainerFilter struct {
	// ID of the container.
	ID string `json:"id,omitempty"`

	// State of the container.
	State *ContainerStateValue `json:"state,omitempty"`

	// LabelSelector to select matches.
	// Only api.MatchLabels is supported for now and the requirements
	// are ANDed. MatchExpressions is not supported yet.
	LabelSelector map[string]string `json:"labelSelector,omitempty"`
}

// ContainerState contains the state of the container.
type ContainerState int32

// ContainerStateValue is the wrapper of ContainerState.
type ContainerStateValue struct {
	// State of the container.
	State ContainerState `json:"state,omitempty"`
}

const (
	// ContainerStateCreated is the state of the created container
	ContainerStateCreated ContainerState = 0
	// ContainerStateRunning is the state of the running container
	ContainerStateRunning ContainerState = 1
	// ContainerStateExited is the state of the container which exited.
	ContainerStateExited ContainerState = 2
	// ContainerStateUnknown is the state of the container which is not known.
	ContainerStateUnknown ContainerState = 3
)

// ContainerStateName stores the mapping from value to string representation of
// state of the container.
var ContainerStateName = map[int32]string{
	0: "CONTAINER_CREATED",
	1: "CONTAINER_RUNNING",
	2: "CONTAINER_EXITED",
	3: "CONTAINER_UNKNOWN",
}

// ContainerStateValues contains the values corresponding to string representation of
// the container state
var ContainerStateValues = map[string]int32{
	"CONTAINER_CREATED": 0,
	"CONTAINER_RUNNING": 1,
	"CONTAINER_EXITED":  2,
	"CONTAINER_UNKNOWN": 3,
}

// ListContainersResponse contains the response of the List containers request.
type ListContainersResponse struct {
	// List of containers.
	Containers []*Container `json:"containers,omitempty"`
}

// Container provides the runtime information for a container, such as ID, hash,
// state of the container.
type Container struct {
	// ID of the container, used by the container runtime to identify
	// a container.
	ID string `json:"id,omitempty"`

	// Metadata of the container.
	Metadata *ContainerMetadata `json:"metadata,omitempty"`

	// Spec of the image.
	Image *ImageSpec `json:"image,omitempty"`

	// Reference to the image in use. For most runtimes, this should be an
	// image ID.
	ImageRef string `json:"imageRef,omitempty"`

	// State of the container.
	State ContainerState `json:"state,omitempty"`

	// Creation time of the container in nanoseconds.
	CreatedAt int64 `json:"createdAt,omitempty"`

	// Key-value pairs that may be used to scope and select individual resources.
	Labels map[string]string `json:"labels,omitempty"`
}

// ContainerStatusRequest contains the information about the container whose status
// we are requesting.
type ContainerStatusRequest struct {
	// ID of the container for which to retrieve status.
	ContainerID string `json:"containerID,omitempty"`
	// Verbose indicates whether to return extra information about the container.
	Verbose bool `json:"verbose,omitempty"`
}

// ContainerStatusResponse contains the response to container status request
type ContainerStatusResponse struct {
	// Status of the container.
	Status *ContainerStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`

	// Info is extra information of the Container. The key could be arbitrary string, and
	// value should be in json format. The information could include anything useful for
	// debug, e.g. pid for linux container based container runtime.
	// It should only be returned non-empty when Verbose is true.
	Info map[string]string `json:"info,omitempty"`
}

// ContainerStatus represents the status of a container.
type ContainerStatus struct {
	// ID of the container.
	ID string `json:"id,omitempty"`
	// Metadata of the container.
	Metadata *ContainerMetadata `json:"metadata,omitempty"`
	// Status of the container.
	State ContainerState `json:"state,omitempty"`
	// Creation time of the container in nanoseconds.
	CreatedAt int64 `json:"createdAt,omitempty"`
	// Start time of the container in nanoseconds. Default: 0 (not specified).
	StartedAt int64 `json:"startedAt,omitempty"`
	// Finish time of the container in nanoseconds. Default: 0 (not specified).
	FinishedAt int64 `json:"finishedAt,omitempty"`
	// Exit code of the container. Only required when finished_at != 0. Default: 0.
	ExitCode int32 `json:"exitCode,omitempty"`
	// Spec of the image.
	Image *ImageSpec `json:"image,omitempty"`
	// Reference to the image in use. For most runtimes, this should be an
	// image ID
	ImageRef string `json:"image_Ref,omitempty"`
	// Brief CamelCase string explaining why container is in its current state.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about why container is in its
	// current state.
	Message string `json:"message,omitempty"`
	// Key-value pairs that may be used to scope and select individual resources.
	Labels map[string]string `json:"labels,omitempty"`
	// Mounts for the container.
	Mounts []*Mount `json:"mounts,omitempty"`
	// Log path of container.
	LogPath string `json:"logPath,omitempty"`
}

// UpdateContainerResourcesRequest is the request structure for update container resource
// request.
type UpdateContainerResourcesRequest struct {
	// ID of the container to update.
	ContainerID string `json:"containerID,omitempty"`
	// Resource configuration specific to Linux containers.
	Linux *LinuxContainerResources `json:"linux,omitempty"`
}

// ReopenContainerLogRequest is the request structure for reopen container log
// request.
type ReopenContainerLogRequest struct {
	// ID of the container for which to reopen the log.
	ContainerID string `json:"containerID,omitempty"`
}

// ExecRequest contains the request for executing the command in the container
// using the container runtime.
type ExecRequest struct {
	// ID of the container in which to execute the command.
	ContainerID string `json:"containerID,omitempty"`
	// Command to execute.
	Cmd []string `json:"cmd,omitempty"`
	// Whether to exec the command in a TTY.
	Tty bool `json:"tty,omitempty"`
	// Whether to stream stdin.
	// One of `stdin`, `stdout`, and `stderr` MUST be true.
	Stdin io.Reader `json:"stdin,omitempty"`
	// Whether to stream stdout.
	// One of `stdin`, `stdout`, and `stderr` MUST be true.
	Stdout io.WriteCloser `json:"stdout,omitempty"`
	// Whether to stream stderr.
	// One of `stdin`, `stdout`, and `stderr` MUST be true.
	// If `tty` is true, `stderr` MUST be false. Multiplexing is not supported
	// in this case. The output of stdout and stderr will be combined to a
	// single stream.
	Stderr io.WriteCloser `json:"stderr,omitempty"`

	Timeout time.Duration `json:"timeout,omitemtpy"`
}

// ExecResponse contains the reponse of the exec command in the container
type ExecResponse struct {
	ContainerID string
	ExitCode    int
}

// StatusRequest represents status request of the underlying container runtime.
type StatusRequest struct {
	// Verbose indicates whether to return extra information about the runtime.
	Verbose bool `json:"verbose,omitempty"`
}

// StatusResponse contains the response for the container runtime status
// reqeuest
type StatusResponse struct {
	// Conditions of the Runtime.
	Conditions []*RuntimeCondition `json:"status,omitempty"`
}

// RuntimeCondition contains condition information for the runtime.
// There are 2 kinds of runtime conditions:
// 1. Required conditions: Conditions are required for kubelet to work
// properly. If any required condition is unmet, the node will be not ready.
// The required conditions include:
//   * RuntimeReady: RuntimeReady means the runtime is up and ready to accept
//   basic containers e.g. container only needs host network.
//   * NetworkReady: NetworkReady means the runtime network is up and ready to
//   accept containers which require container network.
// 2. Optional conditions: Conditions are informative to the user, but kubelet
// will not rely on. Since condition type is an arbitrary string, all conditions
// not required are optional. These conditions will be exposed to users to help
// them understand the status of the system.
//nolint
type RuntimeCondition struct {
	// Type of runtime condition.
	Type string `json:"type,omitempty"`
	// Status of the condition, one of true/false. Default: false.
	Status bool `json:"status,omitempty"`
	// Brief CamelCase string containing reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

// ListImagesRequest is the type corresponding request for image listing using the
// container runtime.
type ListImagesRequest struct {
	// Filter to list images.
	Filter *ImageFilter `json:"filter,omitempty"`
}

// ImageFilter is a wrapper around ImageSpec which can be used to filter images
// based on image spec.
type ImageFilter struct {
	// Spec of the image.
	Image *ImageSpec `json:"image,omitempty"`
}

// ListImagesResponse contains the response for listing images from the container
// runtime.
type ListImagesResponse struct {
	// List of images.
	Images []*Image `json:"images,omitempty"`
}

// Image contains basic information about a container image.
type Image struct {
	// ID of the image.
	ID string `json:"id,omitempty"`
	// Other names by which this image is known.
	RepoTags []string `json:"repoTags,omitempty"`
	// Digests by which this image is known.
	RepoDigests []string `json:"repoDigests,omitempty"`
	// Size of the image in bytes. Must be > 0.
	Size uint64 `json:"size,omitempty"`
	// UID that will run the command(s). This is used as a default if no user is
	// specified when creating the container. UID and the following user name
	// are mutually exclusive.
	UID int64 `json:"uid,omitempty"`
	// User name that will run the command(s). This is used if UID is not set
	// and no user is specified when creating container.
	Username string `json:"username,omitempty"`
}

// ImageStatusRequest is the request for Image status.
type ImageStatusRequest struct {
	// Spec of the image.
	Image *ImageSpec `json:"image,omitempty"`
	// Verbose indicates whether to return extra information about the image.
	Verbose bool `json:"verbose,omitempty"`
}

// ImageStatusResponse is the response for status request for a image from container
// runtime.
type ImageStatusResponse struct {
	// Status of the image.
	Image *Image `json:"image,omitempty"`
	// Info is extra information of the Image. The key could be arbitrary string, and
	// value should be in json format. The information could include anything useful
	// for debug, e.g. image config for oci image based container runtime.
	// It should only be returned non-empty when Verbose is true.
	Info map[string]string `json:"info,omitempty"`
}

// PullImageRequest is the request to pull the image from the container
// registry
type PullImageRequest struct {
	// Spec of the image.
	Image *ImageSpec `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	// Authentication configuration for pulling the image.
	Auth *AuthConfig `protobuf:"bytes,2,opt,name=auth,proto3" json:"auth,omitempty"`
}

// AuthConfig contains authorization information for connecting to a registry.
type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth,omitempty"`
	ServerAddress string `json:"server_address,omitempty"`
	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `json:"identity_token,omitempty"`
	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `json:"registry_token,omitempty"`
}

// PullImageResponse contains the response to pull image request on container runtime.
type PullImageResponse struct {
	// Reference to the image in use. For most runtimes, this should be an
	// image ID or digest.
	ImageRef string `json:"imageRef,omitempty"`
}

// RemoveImageRequest is the request to remove the image from the container
// registry
type RemoveImageRequest struct {
	// Spec of the image to remove.
	Image *ImageSpec `json:"image,omitempty"`
}

var (
	// RuntimeReady means the runtime is up and ready to accept basic containers.
	RuntimeReady = "RuntimeReady"
	// NetworkReady means the runtime network is up and ready to accept containers which require network.
	NetworkReady = "NetworkReady"
)
