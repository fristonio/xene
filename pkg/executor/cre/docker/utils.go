package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	dockerref "github.com/docker/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerfilters "github.com/docker/docker/api/types/filters"
	dockerapi "github.com/docker/docker/client"
	dockermessage "github.com/docker/docker/pkg/jsonmessage"
	dockerstdcopy "github.com/docker/docker/pkg/stdcopy"
	"github.com/fristonio/xene/pkg/executor/cre/runtime"
	godigest "github.com/opencontainers/go-digest"
)

var (
	// Status of a container returned by ListContainers.
	statusRunningPrefix = "Up"
	statusCreatedPrefix = "Created"
	statusExitedPrefix  = "Exited"

	containerNotFoundErrorRegx = regexp.MustCompile(`No such container: [0-9a-z]+`)

	// Delimiter used to construct docker container names.
	nameDelimiter = "_"

	// DockerImageIDPrefix is the prefix of image id in container status.
	DockerImageIDPrefix = "docker://"
	// DockerPullableImageIDPrefix is the prefix of pullable image id in container status.
	DockerPullableImageIDPrefix = "docker-pullable://"

	// defaultImagePullingProgressReportInterval is the default interval of image pulling progress reporting.
	defaultImagePullingProgressReportInterval = 10 * time.Second
)

func generateEnvList(envs []*runtime.KeyValue) (result []string) {
	for _, env := range envs {
		result = append(result, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	return
}

// generateMountBindings converts the mount list to a list of strings that
// can be understood by docker.
// '<HostPath>:<ContainerPath>[:options]', where 'options'
// is a comma-separated list of the following strings:
// 'ro', if the path is read only
// 'Z', if the volume requires SELinux relabeling
// propagation mode such as 'rslave'
func generateMountBindings(mounts []*runtime.Mount) []string {
	result := make([]string, 0, len(mounts))
	for _, m := range mounts {
		bind := fmt.Sprintf("%s:%s", m.HostPath, m.ContainerPath)
		var attrs []string
		if m.Readonly {
			attrs = append(attrs, "ro")
		}

		switch runtime.MountPropagation(m.Propagation) {
		case runtime.MountPropagationPrivate:
			// noop, private is default
		case runtime.MountPropagationBiDirectional:
			attrs = append(attrs, "rshared")
		case runtime.MountPropagationHostToContainer:
			attrs = append(attrs, "rslave")
		default:
			log.Warnf("unknown propagation mode for hostPath %q", m.HostPath)
		}

		if len(attrs) > 0 {
			bind = fmt.Sprintf("%s:%s", bind, strings.Join(attrs, ","))
		}
		result = append(result, bind)
	}

	return result
}

// IsContainerNotFoundError checks whether the error is container not found error.
func IsContainerNotFoundError(err error) bool {
	return containerNotFoundErrorRegx.MatchString(err.Error())
}

func recoverFromCreationConflictIfNeeded(client *dockerapi.Client,
	createConfig dockertypes.ContainerCreateConfig,
	err error) (dockercontainer.ContainerCreateCreatedBody, error) {
	resp := dockercontainer.ContainerCreateCreatedBody{}

	matches := conflictRE.FindStringSubmatch(err.Error())
	if len(matches) != 2 {
		return resp, err
	}

	id := matches[1]
	log.Warningf("Unable to create container due to conflict. Attempting to remove container %q", id)
	if rmErr := client.ContainerRemove(
		context.TODO(), id, dockertypes.ContainerRemoveOptions{RemoveVolumes: true}); rmErr == nil {
		log.Infof("Successfully removed conflicting container %q", id)
		return resp, err
	} else { //nolint
		log.Errorf("Failed to remove the conflicting container %q: %v", id, rmErr)
		// Return if the error is not container not found error.
		if !IsContainerNotFoundError(rmErr) {
			return resp, err
		}
	}

	// randomize the name to avoid conflict.
	createConfig.Name = randomizeName(createConfig.Name)
	log.Infof("create the container with randomized name %s", createConfig.Name)
	return client.ContainerCreate(context.TODO(),
		createConfig.Config, createConfig.HostConfig, createConfig.NetworkingConfig, createConfig.Name)
}

func randomizeName(name string) string {
	return strings.Join([]string{
		name,
		fmt.Sprintf("%08x", rand.Uint32()),
	}, nameDelimiter)
}

// operationTimeout is the error returned when the docker operations are timeout.
type operationTimeout struct {
	err error
}

func (e operationTimeout) Error() string {
	return fmt.Sprintf("operation timeout: %v", e.err)
}

func contextError(ctx context.Context) error {
	if ctx.Err() == context.DeadlineExceeded {
		return operationTimeout{err: ctx.Err()}
	}
	return ctx.Err()
}

// dockerFilter wraps around dockerfilters.Args and provides methods to modify
// the filter easily.
type dockerFilter struct {
	args *dockerfilters.Args
}

func newDockerFilter(args *dockerfilters.Args) *dockerFilter {
	return &dockerFilter{args: args}
}

func (f *dockerFilter) Add(key, value string) {
	f.args.Add(key, value)
}

func (f *dockerFilter) AddLabel(key, value string) {
	f.Add("label", fmt.Sprintf("%s=%s", key, value))
}

func toDockerContainerStatus(state runtime.ContainerState) string {
	switch state {
	case runtime.ContainerStateCreated:
		return "created"
	case runtime.ContainerStateRunning:
		return "running"
	case runtime.ContainerStateExited:
		return "exited"
	case runtime.ContainerStateUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

func toRuntimeContainer(c *dockertypes.Container) (*runtime.Container, error) {
	state := toRuntimeContainerState(c.Status)
	if len(c.Names) == 0 {
		return nil, fmt.Errorf("unexpected empty container name: %+v", c)
	}
	// Docker adds a "/" prefix to names. so trim it.
	name := strings.TrimPrefix(c.Names[0], "/")
	metadata := &runtime.ContainerMetadata{
		Name: name,
	}

	// The timestamp in dockertypes.Container is in seconds.
	createdAt := c.Created * int64(time.Second)
	return &runtime.Container{
		ID:        c.ID,
		Metadata:  metadata,
		Image:     &runtime.ImageSpec{Image: c.Image},
		ImageRef:  c.ImageID,
		State:     state,
		CreatedAt: createdAt,
		Labels:    c.Labels,
	}, nil
}

func toRuntimeContainerState(state string) runtime.ContainerState {
	// Parse the state string in dockertypes.Container. This could break when
	// we upgrade docker.
	switch {
	case strings.HasPrefix(state, statusRunningPrefix):
		return runtime.ContainerStateRunning
	case strings.HasPrefix(state, statusExitedPrefix):
		return runtime.ContainerStateExited
	case strings.HasPrefix(state, statusCreatedPrefix):
		return runtime.ContainerStateCreated
	default:
		return runtime.ContainerStateUnknown
	}
}

// parseDockerTimestamp parses the timestamp returned by Interface from string to time.Time
func parseDockerTimestamp(s string) (time.Time, error) {
	// Timestamp returned by Docker is in time.RFC3339Nano format.
	return time.Parse(time.RFC3339Nano, s)
}

func getContainerTimestamps(r *dockertypes.ContainerJSON) (time.Time, time.Time, time.Time, error) {
	var createdAt, startedAt, finishedAt time.Time
	var err error

	createdAt, err = parseDockerTimestamp(r.Created)
	if err != nil {
		return createdAt, startedAt, finishedAt, err
	}
	startedAt, err = parseDockerTimestamp(r.State.StartedAt)
	if err != nil {
		return createdAt, startedAt, finishedAt, err
	}
	finishedAt, err = parseDockerTimestamp(r.State.FinishedAt)
	if err != nil {
		return createdAt, startedAt, finishedAt, err
	}
	return createdAt, startedAt, finishedAt, nil
}

// matchImageIDOnly checks that the given image specifier is a digest-only
// reference, and that it matches the given image.
func matchImageIDOnly(inspected dockertypes.ImageInspect, image string) bool {
	// If the image ref is literally equal to the inspected image's ID,
	// just return true here (this might be the case for Docker 1.9,
	// where we won't have a digest for the ID)
	if inspected.ID == image {
		return true
	}

	// Otherwise, we should try actual parsing to be more correct
	ref, err := dockerref.Parse(image)
	if err != nil {
		log.Infof("couldn't parse image reference %q: %v", image, err)
		return false
	}

	digest, isDigested := ref.(dockerref.Digested)
	if !isDigested {
		log.Infof("the image reference %q was not a digest reference", image)
		return false
	}

	id, err := godigest.Parse(inspected.ID)
	if err != nil {
		log.Infof("couldn't parse image ID reference %q: %v", id, err)
		return false
	}

	if digest.Digest().Algorithm().String() == id.Algorithm().String() && digest.Digest().Hex() == id.Hex() {
		return true
	}

	log.Infof("The reference %s does not directly refer to the given image's ID (%q)", image, inspected.ID)
	return false
}

func toPullableImageID(id string, image *dockertypes.ImageInspect) string {
	// Default to the image ID, but if RepoDigests is not empty, use
	// the first digest instead.
	imageID := DockerImageIDPrefix + id
	if image != nil && len(image.RepoDigests) > 0 {
		imageID = DockerPullableImageIDPrefix + image.RepoDigests[0]
	}
	return imageID
}

func imageToRuntimeAPIImage(image *dockertypes.ImageSummary) (*runtime.Image, error) {
	if image == nil {
		return nil, fmt.Errorf("unable to convert a nil pointer to a runtime API image")
	}

	size := uint64(image.VirtualSize)
	return &runtime.Image{
		ID:          image.ID,
		RepoTags:    image.RepoTags,
		RepoDigests: image.RepoDigests,
		Size:        size,
	}, nil
}

func imageInspectToRuntimeImage(image *dockertypes.ImageInspect) (*runtime.Image, error) {
	if image == nil || image.Config == nil {
		return nil, fmt.Errorf("unable to convert a nil pointer to a runtime API image")
	}

	size := uint64(image.VirtualSize)
	runtimeImage := &runtime.Image{
		ID:          image.ID,
		RepoTags:    image.RepoTags,
		RepoDigests: image.RepoDigests,
		Size:        size,
	}

	uid, username := getUserFromImageUser(image.Config.User)
	if uid != nil {
		runtimeImage.UID = *uid
	}
	runtimeImage.Username = username
	return runtimeImage, nil
}

// parseUserFromImageUser splits the user out of an user:group string.
func parseUserFromImageUser(id string) string {
	if id == "" {
		return id
	}
	// split instances where the id may contain user:group
	if strings.Contains(id, ":") {
		return strings.Split(id, ":")[0]
	}
	// no group, just return the id
	return id
}

// getUserFromImageUser gets uid or user name of the image user.
// If user is numeric, it will be treated as uid; or else, it is treated as user name.
func getUserFromImageUser(imageUser string) (*int64, string) {
	user := parseUserFromImageUser(imageUser)
	// return both nil if user is not specified in the image.
	if user == "" {
		return nil, ""
	}
	// user could be either uid or user name. Try to interpret as numeric uid.
	uid, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		// If user is non numeric, assume it's user name.
		return nil, user
	}
	// If user is a numeric uid.
	return &uid, ""
}

func base64EncodeAuth(auth *dockertypes.AuthConfig) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(auth); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

// progress is a wrapper of dockermessage.JSONMessage with a lock protecting it.
type progress struct {
	sync.RWMutex
	// message stores the latest docker json message.
	message *dockermessage.JSONMessage
	// timestamp of the latest update.
	timestamp time.Time
}

func newProgress() *progress {
	return &progress{timestamp: time.Now()}
}

func (p *progress) set(msg *dockermessage.JSONMessage) {
	p.Lock()
	defer p.Unlock()
	p.message = msg
	p.timestamp = time.Now()
}

func (p *progress) get() (string, time.Time) {
	p.RLock()
	defer p.RUnlock()
	if p.message == nil {
		return "No progress", p.timestamp
	}
	// The following code is based on JSONMessage.Display
	var prefix string
	if p.message.ID != "" {
		prefix = fmt.Sprintf("%s: ", p.message.ID)
	}
	if p.message.Progress == nil {
		return fmt.Sprintf("%s%s", prefix, p.message.Status), p.timestamp
	}
	return fmt.Sprintf("%s%s %s", prefix, p.message.Status, p.message.Progress.String()), p.timestamp
}

// progressReporter keeps the newest image pulling progress and periodically report the newest progress.
type progressReporter struct {
	*progress
	image                     string
	cancel                    context.CancelFunc
	stopCh                    chan struct{}
	imagePullProgressDeadline time.Duration
}

// newProgressReporter creates a new progressReporter for specific image with specified reporting interval
func newProgressReporter(image string, cancel context.CancelFunc,
	imagePullProgressDeadline time.Duration) *progressReporter {
	return &progressReporter{
		progress:                  newProgress(),
		image:                     image,
		cancel:                    cancel,
		stopCh:                    make(chan struct{}),
		imagePullProgressDeadline: imagePullProgressDeadline,
	}
}

// start starts the progressReporter
func (p *progressReporter) start() {
	go func() {
		ticker := time.NewTicker(defaultImagePullingProgressReportInterval)
		defer ticker.Stop()
		for {
			// TODO(random-liu): Report as events.
			select {
			case <-ticker.C:
				progress, timestamp := p.progress.get()
				// If there is no progress for p.imagePullProgressDeadline, cancel the operation.
				if time.Since(timestamp) > p.imagePullProgressDeadline {
					log.Errorf("Cancel pulling image %q because of no progress for %v, latest progress: %q",
						p.image, p.imagePullProgressDeadline, progress)
					p.cancel()
					return
				}
				log.Infof("Pulling image %q: %q", p.image, progress)
			case <-p.stopCh:
				progress, _ := p.progress.get()
				log.Infof("Stop pulling image %q: %q", p.image, progress)
				return
			}
		}
	}()
}

// stop stops the progressReporter
func (p *progressReporter) stop() {
	close(p.stopCh)
}

// getImageRef returns the image digest if exists, or else returns the image ID.
func getImageRef(client *dockerapi.Client, image string) (string, error) {
	img, _, err := client.ImageInspectWithRaw(context.TODO(), image)
	if err != nil {
		return "", err
	}
	if img.ID == "" {
		return "", fmt.Errorf("unable to inspect image %s", image)
	}

	// Returns the digest if it exist.
	if len(img.RepoDigests) > 0 {
		return img.RepoDigests[0], nil
	}

	return img.ID, nil
}

// holdHijackedConnection hold the HijackedResponse, redirect the inputStream to the connection,
// and redirect the response
// stream to stdout and stderr. NOTE: If needed, we could also add context in this function.
func holdHijackedConnection(tty bool, inputStream io.Reader, outputStream, errorStream io.Writer,
	resp dockertypes.HijackedResponse) error {
	receiveStdout := make(chan error)
	if outputStream != nil || errorStream != nil {
		go func() {
			receiveStdout <- redirectResponseToOutputStream(tty, outputStream, errorStream, resp.Reader)
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			_, _ = io.Copy(resp.Conn, inputStream)

		}
		_ = resp.CloseWrite()
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		return err
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			return <-receiveStdout
		}
	}
	return nil
}

// redirectResponseToOutputStream redirect the response stream to stdout and stderr. When tty is true, all stream will
// only be redirected to stdout.
func redirectResponseToOutputStream(tty bool, outputStream, errorStream io.Writer, resp io.Reader) error {
	if outputStream == nil {
		outputStream = ioutil.Discard
	}
	if errorStream == nil {
		errorStream = ioutil.Discard
	}
	var err error
	if tty {
		_, err = io.Copy(outputStream, resp)
	} else {
		_, err = dockerstdcopy.StdCopy(outputStream, errorStream, resp)
	}
	return err
}
