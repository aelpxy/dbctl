package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
)

const readyPollInterval = time.Second

// health statuses reported by the docker engine
const (
	healthy   = "healthy"
	unhealthy = "unhealthy"
)

// Start starts a stopped database.
func (c *Client) Start(ctx context.Context, id string) error {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return err
	}

	if err := c.api.ContainerStart(ctx, info.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("start container %s: %w", id, err)
	}

	return nil
}

// Stop stops a running database, keeping its container and data.
func (c *Client) Stop(ctx context.Context, id string) error {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return err
	}

	if err := c.api.ContainerStop(ctx, info.ID, container.StopOptions{}); err != nil {
		return fmt.Errorf("stop container %s: %w", id, err)
	}

	return nil
}

// Restart stops and starts a database.
func (c *Client) Restart(ctx context.Context, id string) error {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return err
	}

	if err := c.api.ContainerRestart(ctx, info.ID, container.StopOptions{}); err != nil {
		return fmt.Errorf("restart container %s: %w", id, err)
	}

	return nil
}

// WaitReady blocks until the database reports healthy, or is running when it has no healthcheck.
func (c *Client) WaitReady(ctx context.Context, id string) error {
	ticker := time.NewTicker(readyPollInterval)
	defer ticker.Stop()

	for {
		ready, err := c.ready(ctx, id)
		if err != nil {
			return err
		}

		if ready {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for database %s: %w", id, ctx.Err())
		case <-ticker.C:
		}
	}
}

func (c *Client) ready(ctx context.Context, id string) (bool, error) {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return false, err
	}

	if info.State.Dead || (!info.State.Running && !info.State.Restarting) {
		return false, fmt.Errorf("%w: %s", ErrNotRunning, id)
	}

	if info.State.Health == nil {
		return info.State.Running, nil
	}

	switch info.State.Health.Status {
	case healthy:
		return true, nil
	case unhealthy:
		return false, fmt.Errorf("%w: %s", ErrUnhealthy, id)
	default:
		return false, nil
	}
}
