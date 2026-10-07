package docker

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Database is a container managed by dbctl.
type Database struct {
	Created time.Time
	Labels  map[string]string
	Env     map[string]string
	ID      string
	Name    string
	Image   string
	State   string
	Status  string
	Health  string
	Type    string
	// Password is empty for databases created before dbctl stored it.
	Password string
	Ports    []Port
	Volumes  []Volume
	Running  bool
}

// Port is a container port published on the host.
type Port struct {
	HostIP        string
	HostPort      string
	ContainerPort string
}

// Volume is a volume mounted into a database container.
type Volume struct {
	Name   string
	Source string
	Target string
}

// Inspect returns the database with the given container ID or name.
func (c *Client) Inspect(ctx context.Context, id string) (*Database, error) {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return nil, err
	}

	return c.fromInspect(&info), nil
}

// List returns every database managed by dbctl, including stopped ones.
func (c *Client) List(ctx context.Context) ([]Database, error) {
	containers, err := c.api.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", ContainerPrefix)),
	})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	dbs := make([]Database, 0, len(containers))

	for i := range containers {
		ctr := &containers[i]
		if len(ctr.Names) == 0 || !isManaged(ctr.Names[0]) {
			continue
		}

		dbs = append(dbs, Database{
			Created: time.Unix(ctr.Created, 0),
			Labels:  ctr.Labels,
			ID:      ctr.ID,
			Name:    displayName(ctr.Names[0]),
			Image:   ctr.Image,
			State:   ctr.State,
			Status:  ctr.Status,
			Type:    c.typeName(ctr.Labels, ctr.Image),
			Running: ctr.State == "running",
		})
	}

	return dbs, nil
}

func (c *Client) inspect(ctx context.Context, id string) (types.ContainerJSON, error) {
	info, err := c.api.ContainerInspect(ctx, id)
	if client.IsErrNotFound(err) {
		return types.ContainerJSON{}, fmt.Errorf("%w: %s", ErrNotFound, id)
	}

	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("inspect container %s: %w", id, err)
	}

	if !isManaged(info.Name) {
		return types.ContainerJSON{}, fmt.Errorf("%w: %s", ErrNotManaged, displayName(info.Name))
	}

	return info, nil
}

func (c *Client) fromInspect(info *types.ContainerJSON) *Database {
	db := &Database{
		Labels:   info.Config.Labels,
		Env:      envMap(info.Config.Env),
		ID:       info.ID,
		Name:     displayName(info.Name),
		Image:    info.Config.Image,
		State:    info.State.Status,
		Status:   info.State.Status,
		Type:     c.typeName(info.Config.Labels, info.Config.Image),
		Password: info.Config.Labels[passwordLabel],
		Volumes:  volumes(info.Mounts),
		Running:  info.State.Running,
	}

	if created, err := time.Parse(time.RFC3339Nano, info.Created); err == nil {
		db.Created = created
	}

	if info.NetworkSettings != nil {
		db.Ports = ports(info.NetworkSettings.Ports)
	}

	if info.State.Health != nil {
		db.Health = info.State.Health.Status
	}

	return db
}

func ports(pm nat.PortMap) []Port {
	out := make([]Port, 0, len(pm))

	for _, port := range slices.Sorted(maps.Keys(pm)) {
		for _, b := range pm[port] {
			out = append(out, Port{HostIP: b.HostIP, HostPort: b.HostPort, ContainerPort: port.Port()})
		}
	}

	return out
}

func volumes(mounts []types.MountPoint) []Volume {
	out := make([]Volume, 0, len(mounts))

	for i := range mounts {
		m := &mounts[i]
		if !strings.HasPrefix(m.Name, volumePrefix) {
			continue
		}

		out = append(out, Volume{Name: m.Name, Source: m.Source, Target: m.Destination})
	}

	return out
}
