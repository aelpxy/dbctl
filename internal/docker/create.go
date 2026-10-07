package docker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

const (
	healthInterval = 2 * time.Second
	healthTimeout  = 5 * time.Second
	healthRetries  = 30
)

// CreateOptions configures a new database container.
type CreateOptions struct {
	Definition *database.Definition
	Image      string
	Name       string
	Password   string
	HostIP     string
	HostPort   int
}

// Create creates and starts a database container with its data volume and returns the container ID.
func (c *Client) Create(ctx context.Context, opts *CreateOptions) (string, error) {
	cfg, err := containerConfig(opts)
	if err != nil {
		return "", err
	}

	volumeName := volumePrefix + opts.Name

	if _, err := c.api.VolumeCreate(ctx, volume.CreateOptions{Name: volumeName}); err != nil {
		return "", fmt.Errorf("create volume %s: %w", volumeName, err)
	}

	resp, err := c.api.ContainerCreate(ctx, cfg, hostConfig(opts, volumeName), nil, nil, ContainerPrefix+opts.Name)
	if err != nil {
		// not forced so a volume still used by an existing database is never removed, a failure here is expected then
		_ = c.api.VolumeRemove(ctx, volumeName, false)

		return "", fmt.Errorf("create container: %w", err)
	}

	if err := c.api.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("start container %s: %w", resp.ID, err)
	}

	return resp.ID, nil
}

func containerConfig(opts *CreateOptions) (*container.Config, error) {
	def := opts.Definition

	env, err := def.Env(opts.Password)
	if err != nil {
		return nil, fmt.Errorf("render env: %w", err)
	}

	cmd, err := def.Command(opts.Password)
	if err != nil {
		return nil, fmt.Errorf("render command: %w", err)
	}

	healthcheck, err := def.Healthcheck(opts.Password)
	if err != nil {
		return nil, fmt.Errorf("render healthcheck: %w", err)
	}

	return &container.Config{
		Image:        opts.Image,
		Env:          env,
		Cmd:          cmd,
		Healthcheck:  healthConfig(healthcheck),
		ExposedPorts: nat.PortSet{containerPort(def): struct{}{}},
		Labels: map[string]string{
			typeLabel:     def.Name(),
			passwordLabel: opts.Password,
		},
	}, nil
}

// a nil config keeps any healthcheck defined by the image
func healthConfig(test []string) *container.HealthConfig {
	if len(test) == 0 {
		return nil
	}

	return &container.HealthConfig{
		Test:     append([]string{"CMD"}, test...),
		Interval: healthInterval,
		Timeout:  healthTimeout,
		Retries:  healthRetries,
	}
}

func hostConfig(opts *CreateOptions, volumeName string) *container.HostConfig {
	return &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: opts.Definition.DataDir(),
			},
		},
		NetworkMode:   container.NetworkMode(networkName),
		RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyUnlessStopped},
		PortBindings: nat.PortMap{
			containerPort(opts.Definition): []nat.PortBinding{
				{
					HostIP:   opts.HostIP,
					HostPort: strconv.Itoa(opts.HostPort),
				},
			},
		},
	}
}

func containerPort(def *database.Definition) nat.Port {
	return nat.Port(strconv.Itoa(def.Port()) + "/tcp")
}
