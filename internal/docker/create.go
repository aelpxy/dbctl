package docker

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"time"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
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
	port, err := network.ParsePort(strconv.Itoa(opts.Definition.Port()) + "/tcp")
	if err != nil {
		return "", fmt.Errorf("parse container port: %w", err)
	}

	cfg, err := containerConfig(opts, port)
	if err != nil {
		return "", err
	}

	volumeName := volumePrefix + opts.Name

	host, err := hostConfig(opts, port, volumeName)
	if err != nil {
		return "", err
	}

	if _, err := c.api.VolumeCreate(ctx, client.VolumeCreateOptions{Name: volumeName}); err != nil {
		return "", fmt.Errorf("create volume %s: %w", volumeName, err)
	}

	resp, err := c.api.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     cfg,
		HostConfig: host,
		Name:       ContainerPrefix + opts.Name,
	})
	if err != nil {
		// not forced so a volume still used by an existing database is never removed, a failure here is expected then
		_, _ = c.api.VolumeRemove(ctx, volumeName, client.VolumeRemoveOptions{})

		return "", fmt.Errorf("create container: %w", err)
	}

	if _, err := c.api.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("start container %s: %w", resp.ID, err)
	}

	return resp.ID, nil
}

func containerConfig(opts *CreateOptions, port network.Port) (*container.Config, error) {
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
		ExposedPorts: network.PortSet{port: struct{}{}},
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

func hostConfig(opts *CreateOptions, port network.Port, volumeName string) (*container.HostConfig, error) {
	hostIP, err := netip.ParseAddr(opts.HostIP)
	if err != nil {
		return nil, fmt.Errorf("parse host ip %q: %w", opts.HostIP, err)
	}

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
		PortBindings: network.PortMap{
			port: []network.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: strconv.Itoa(opts.HostPort),
				},
			},
		},
	}, nil
}
