package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ImageExists reports whether the image is available locally.
func (c *Client) ImageExists(ctx context.Context, ref string) (bool, error) {
	_, _, err := c.api.ImageInspectWithRaw(ctx, ref)
	if client.IsErrNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("inspect image %s: %w", ref, err)
	}

	return true, nil
}

// Pull downloads the image from its registry.
func (c *Client) Pull(ctx context.Context, ref string) error {
	out, err := c.api.ImagePull(ctx, ref, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("pull image %s: %w", ref, err)
	}

	defer func() { _ = out.Close() }()

	// the pull only completes once the progress stream has been fully read
	if _, err := io.Copy(io.Discard, out); err != nil {
		return fmt.Errorf("read pull progress for %s: %w", ref, err)
	}

	return nil
}
