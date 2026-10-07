package docker

import (
	"fmt"
	"strconv"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

func CreateContainer(imageName, dbType, containerName string, externalPort int, password string, envVars ...string) (string, error) {
	dockerClient, err := DockerClient()
	if err != nil {
		return "", fmt.Errorf("error creating docker client: %w", err)
	}

	var internalPort int
	var mountTarget string
	var cmd []string

	switch dbType {
	case "postgres":
		internalPort = 5432
		// postgres 18+ images store data under a versioned subdirectory of /var/lib/postgresql
		mountTarget = "/var/lib/postgresql"
	case "redis":
		internalPort = 6379
		mountTarget = "/data"
		cmd = []string{"redis-server", "--appendonly", "yes", "--requirepass", password}
	case "mysql", "mariadb":
		internalPort = 3306
		mountTarget = "/var/lib/mysql"
	case "mongo":
		internalPort = 27017
		mountTarget = "/data/db"
	case "meilisearch":
		internalPort = 7700
		mountTarget = "/meili_data"
		cmd = []string{"meilisearch", "--master-key", password}
	case "keydb":
		internalPort = 6379
		mountTarget = "/data"
		cmd = []string{"keydb-server", "/etc/keydb/keydb.conf", "--appendonly", "yes", "--requirepass", password}
	case "couchdb":
		internalPort = 5984
		mountTarget = "/opt/couchdb/data"
	case "clickhouse":
		internalPort = 9000
		mountTarget = "/var/lib/clickhouse"
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}

	volumeName := config.DockerVolumeName + containerName

	_, err = dockerClient.VolumeCreate(Ctx, volume.CreateOptions{Name: volumeName})
	if err != nil {
		return "", fmt.Errorf("error creating volume: %w", err)
	}

	port := nat.Port(strconv.Itoa(internalPort) + "/tcp")

	containerConfig := &container.Config{
		Image:        imageName,
		Env:          envVars,
		Cmd:          cmd,
		ExposedPorts: nat.PortSet{port: struct{}{}},
		Labels:       map[string]string{config.DockerTypeLabel: dbType},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: mountTarget,
			},
		},
		NetworkMode:   container.NetworkMode(config.DockerNetworkName),
		RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyAlways},
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{
					HostIP:   utils.GetIP().String(),
					HostPort: strconv.Itoa(externalPort),
				},
			},
		},
	}

	resp, err := dockerClient.ContainerCreate(Ctx, containerConfig, hostConfig, nil, nil, config.DockerContainerPrefix+containerName)
	if err != nil {
		// not forced so a volume still used by an existing database is never removed
		_ = dockerClient.VolumeRemove(Ctx, volumeName, false)
		return "", fmt.Errorf("error creating container: %w", err)
	}

	err = dockerClient.ContainerStart(Ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("error starting container: %w", err)
	}

	return resp.ID, nil
}
