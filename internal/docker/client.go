// Package docker manages dbctl databases through the Docker Engine API.
package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/moby/moby/client"
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

	ErrDaemonUnreachable = errors.New("cannot reach the docker daemon")
)

// Client talks to the Docker daemon on behalf of dbctl.
type Client struct {
	api       *client.Client
	databases *database.Registry
}

// New connects to the Docker daemon from the environment and ensures the dbctl network exists.
func New(ctx context.Context, databases *database.Registry) (*Client, error) {
	api, err := client.New(client.FromEnv)
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
	if _, err := c.api.Ping(ctx, client.PingOptions{}); err != nil {
		// docker's own message repeats the socket path and advice the caller already gives
		return fmt.Errorf("%w at %s", ErrDaemonUnreachable, c.api.DaemonHost())
	}

	if _, err := c.api.NetworkInspect(ctx, networkName, client.NetworkInspectOptions{}); err == nil {
		return nil
	}

	if _, err := c.api.NetworkCreate(ctx, networkName, client.NetworkCreateOptions{}); err != nil {
		return fmt.Errorf("create network %s: %w", networkName, err)
	}

	return nil
}

// docker reports container names with a leading slash
func containerName(name string) string {
	return strings.TrimPrefix(name, "/")
}

// displayName is the short name users type, without the dbctl prefix.
func displayName(name string) string {
	return strings.TrimPrefix(containerName(name), ContainerPrefix)
}

func isManaged(name string) bool {
	return strings.HasPrefix(containerName(name), ContainerPrefix)
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
