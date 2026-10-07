// Package docker manages dbctl databases through the Docker Engine API.
package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// ContainerPrefix prefixes the name of every container managed by dbctl.
const ContainerPrefix = "dbctl."

const (
	networkName   = "dbctl.network"
	volumePrefix  = "dbctl.volume."
	typeLabel     = "dbctl.type"
	passwordLabel = "dbctl.password"
)

var (
	ErrNotFound   = errors.New("database not found")
	ErrNotManaged = errors.New("container is not managed by dbctl")
	ErrNotRunning = errors.New("database is not running")
	ErrUnhealthy  = errors.New("database is unhealthy")
)

// Client talks to the Docker daemon on behalf of dbctl.
type Client struct {
	api       *client.Client
	databases *database.Registry
}

// New connects to the Docker daemon from the environment and ensures the dbctl network exists.
func New(ctx context.Context, databases *database.Registry) (*Client, error) {
	api, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}

	c := &Client{api: api, databases: databases}

	if err := c.setup(ctx); err != nil {
		return nil, errors.Join(err, c.Close())
	}

	return c, nil
}

// Close releases the connection to the Docker daemon.
func (c *Client) Close() error {
	if err := c.api.Close(); err != nil {
		return fmt.Errorf("close docker client: %w", err)
	}

	return nil
}

func (c *Client) setup(ctx context.Context) error {
	if _, err := c.api.Ping(ctx); err != nil {
		return fmt.Errorf("reach docker daemon, make sure docker is installed and running: %w", err)
	}

	if _, err := c.api.NetworkInspect(ctx, networkName, network.InspectOptions{}); err == nil {
		return nil
	}

	if _, err := c.api.NetworkCreate(ctx, networkName, network.CreateOptions{}); err != nil {
		return fmt.Errorf("create network %s: %w", networkName, err)
	}

	return nil
}

// docker reports container names with a leading slash
func displayName(name string) string {
	return strings.TrimPrefix(name, "/")
}

func isManaged(name string) bool {
	return strings.HasPrefix(displayName(name), ContainerPrefix)
}

// containers created before the type label existed fall back to their image
func (c *Client) definition(labels map[string]string, image string) (*database.Definition, error) {
	def, err := c.databases.Resolve(labels[typeLabel], image)
	if err != nil {
		return nil, fmt.Errorf("resolve database type: %w", err)
	}

	return def, nil
}

func (c *Client) typeName(labels map[string]string, image string) string {
	def, err := c.definition(labels, image)
	if err != nil {
		return ""
	}

	return def.Name()
}
