package docker

import (
	"context"
	"encoding/json"
	"io"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockerfilters "github.com/docker/docker/api/types/filters"
	dockerapi "github.com/docker/docker/client"
	dockermessage "github.com/docker/docker/pkg/jsonmessage"
	"github.com/fristonio/xene/pkg/executor/cre/runtime"
)

var (
	imagePullProgressDeadline time.Duration = time.Minute * 5
)

// ListImages lists all the images matching the criteria in ImageRequest
func (e *RuntimeExecutor) ListImages(ctx context.Context,
	r *runtime.ListImagesRequest) (*runtime.ListImagesResponse, error) {
	filter := r.Filter
	opts := dockertypes.ImageListOptions{}
	if filter != nil {
		if filter.Image.Image != "" {
			opts.Filters = dockerfilters.NewArgs()
			opts.Filters.Add("reference", filter.Image.Image)
		}
	}

	images, err := e.client.ImageList(ctx, opts)
	if ctxErr := contextError(ctx); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil {
		return nil, err
	}

	result := make([]*runtime.Image, 0, len(images))
	for _, i := range images {
		apiImage, err := imageToRuntimeAPIImage(&i)
		if err != nil {
			log.Infof("Failed to convert docker API image %+v to runtime API image: %v", i, err)
			continue
		}
		result = append(result, apiImage)
	}
	return &runtime.ListImagesResponse{Images: result}, nil
}

// ImageStatus returns the status of the image. If the image is not
// present, returns a response with ImageStatusResponse.Image set to nil.
func (e *RuntimeExecutor) ImageStatus(ctx context.Context,
	r *runtime.ImageStatusRequest) (*runtime.ImageStatusResponse, error) {
	image := r.Image

	imageInspect, _, err := e.client.ImageInspectWithRaw(ctx, image.Image)
	if err != nil {
		if isImageNotFoundError(err) {
			return &runtime.ImageStatusResponse{}, nil
		}
		return nil, err
	}

	imageStatus, err := imageInspectToRuntimeImage(&imageInspect)
	if err != nil {
		return nil, err
	}

	res := runtime.ImageStatusResponse{Image: imageStatus}
	if r.Verbose {
		res.Info = imageInspect.Config.Labels
	}
	return &res, nil
}

// PullImage pulls an image with authentication config.
func (e *RuntimeExecutor) PullImage(ctx context.Context,
	r *runtime.PullImageRequest) (*runtime.PullImageResponse, error) {
	image := r.Image
	auth := r.Auth
	authConfig := dockertypes.AuthConfig{}

	if auth != nil {
		authConfig.Username = auth.Username
		authConfig.Password = auth.Password
		authConfig.ServerAddress = auth.ServerAddress
		authConfig.IdentityToken = auth.IdentityToken
		authConfig.RegistryToken = auth.RegistryToken
	}

	// RegistryAuth is the base64 encoded credentials for the registry
	base64Auth, err := base64EncodeAuth(&authConfig)
	if err != nil {
		return nil, err
	}

	opts := dockertypes.ImagePullOptions{}
	opts.RegistryAuth = base64Auth
	resp, err := e.client.ImagePull(ctx, image.Image, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	_, cancel := context.WithCancel(ctx)
	reporter := newProgressReporter(image.Image, cancel, imagePullProgressDeadline)
	reporter.start()
	defer reporter.stop()
	decoder := json.NewDecoder(resp)
	for {
		var msg dockermessage.JSONMessage
		err := decoder.Decode(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if msg.Error != nil {
			return nil, msg.Error
		}
		reporter.set(&msg)
	}

	imageRef, err := getImageRef(e.client, image.Image)
	if err != nil {
		return nil, err
	}

	return &runtime.PullImageResponse{ImageRef: imageRef}, nil
}

// RemoveImage removes the image.
// This call is idempotent, and must not return an error if the image has
// already been removed.
func (e *RuntimeExecutor) RemoveImage(ctx context.Context, r *runtime.RemoveImageRequest) error {
	image := r.Image
	// If the image has multiple tags, we need to remove all the tags
	// TODO: We assume image.Image is image ID here, which is true in the current implementation
	// of kubelet, but we should still clarify this in CRI.
	imageInspect, _, err := e.client.ImageInspectWithRaw(ctx, image.Image)

	// dockerclient.InspectImageByID doesn't work with digest and repoTags,
	// it is safe to continue removing it since there is another check below.
	if err != nil && !isImageNotFoundError(err) {
		return err
	}

	if isImageNotFoundError(err) {
		// image is nil, assuming it doesn't exist.
		return nil
	}

	// An image can have different numbers of RepoTags and RepoDigests.
	// Iterating over both of them plus the image ID ensures the image really got removed.
	// It also prevents images from being deleted, which actually are deletable using this approach.
	var images []string
	images = append(images, imageInspect.RepoTags...)
	images = append(images, imageInspect.RepoDigests...)
	images = append(images, image.Image)

	for _, image := range images {
		_, err := e.client.ImageRemove(ctx, image, dockertypes.ImageRemoveOptions{PruneChildren: true})
		if ctxErr := contextError(ctx); ctxErr != nil {
			return ctxErr
		}
		if dockerapi.IsErrNotFound(err) {
			return imageNotFoundError{ID: image}
		}

		if err != nil && !isImageNotFoundError(err) {
			return err
		}
	}

	return nil
}
