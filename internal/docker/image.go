package docker

import (
	"context"
	"fmt"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/client"
)

// ImageExists reports whether the image is available locally.
func (c *Client) ImageExists(ctx context.Context, ref string) (bool, error) {
	_, err := c.api.ImageInspect(ctx, ref)
	if cerrdefs.IsNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("inspect image %s: %w", ref, err)
	}

	return true, nil
}

// Pull downloads the image from its registry.
func (c *Client) Pull(ctx context.Context, ref string) error {
	resp, err := c.api.ImagePull(ctx, ref, client.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pull image %s: %w", ref, err)
	}

	// Wait drains the progress stream and reports errors the registry sent inside it
	if err := resp.Wait(ctx); err != nil {
		return fmt.Errorf("pull image %s: %w", ref, err)
	}

	return nil
}
