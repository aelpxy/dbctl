package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
)

// DeleteOptions configures how a database is deleted.
type DeleteOptions struct {
	KeepVolumes bool
}

// Delete stops and removes a database container and, unless KeepVolumes is set, its data volumes.
// It returns the names of the removed volumes.
func (c *Client) Delete(ctx context.Context, id string, opts DeleteOptions) ([]string, error) {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := c.api.ContainerRemove(ctx, info.ID, container.RemoveOptions{Force: true}); err != nil {
		return nil, fmt.Errorf("remove container %s: %w", id, err)
	}

	if opts.KeepVolumes {
		return nil, nil
	}

	vols := volumes(info.Mounts)
	removed := make([]string, 0, len(vols))

	for i := range vols {
		if err := c.api.VolumeRemove(ctx, vols[i].Name, true); err != nil {
			return removed, fmt.Errorf("remove volume %s: %w", vols[i].Name, err)
		}

		removed = append(removed, vols[i].Name)
	}

	return removed, nil
}
