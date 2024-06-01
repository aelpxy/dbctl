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

	_, err = dockerClient.VolumeCreate(Ctx, volume.CreateOptions{
		Name: config.DockerVolumeName + containerName,
	})

	if err != nil {
		return "", fmt.Errorf("error creating volume: %w", err)
	}

	var internalPort int
	var mountSource string
	var mountTarget string
	var cmd []string

	hostIP := utils.GetIP().String()

	switch dbType {
	case "postgres":
		internalPort = 5432
	case "redis":
		internalPort = 6379
	case "mysql":
		internalPort = 3306
	case "mariadb":
		internalPort = 3306
	case "mongo":
		internalPort = 27017
	case "meilisearch":
		internalPort = 7700
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}

	containerConfig := &container.Config{
		Image: imageName,
		Env:   envVars,
		ExposedPorts: nat.PortSet{
			nat.Port(strconv.Itoa(internalPort) + "/tcp"): struct{}{},
		},
	}

	switch dbType {
	case "postgres":
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/var/lib/postgresql/data"

		containerConfig.Env = append(containerConfig.Env,
			"POSTGRES_DB=postgres",
			"POSTGRES_USER=postgres",
		)
	case "redis":
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/data"

		cmd = []string{"redis-server", "--requirepass", password}
	case "mysql":
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/var/lib/mysql"

		containerConfig.Env = append(containerConfig.Env, "MYSQL_DATABASE=db")
	case "mariadb":
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/var/lib/mysql"

		containerConfig.Env = append(containerConfig.Env, "MARIADB_DATABASE=db")
	case "mongo":
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/data/db"

		containerConfig.Env = append(containerConfig.Env,
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_DATABASE=db",
		)
	case "meilisearch":
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/meili_data"

		cmd = []string{"meilisearch", "--master-key", password}
	default:
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/data"
	}

	if len(cmd) > 0 {
		containerConfig.Cmd = cmd
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: mountSource,
				Target: mountTarget,
			},
		},
		NetworkMode:   container.NetworkMode(config.DockerNetworkName),
		RestartPolicy: container.RestartPolicy{Name: "always"},
		PortBindings: nat.PortMap{
			nat.Port(strconv.Itoa(internalPort) + "/tcp"): []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: strconv.Itoa(externalPort),
				},
			},
		},
	}

	containerName = fmt.Sprintf("%s%s", config.DockerContainerPrefix, containerName)

	resp, err := dockerClient.ContainerCreate(Ctx, containerConfig, hostConfig, nil, nil, containerName)

	if err != nil {
		return "", fmt.Errorf("error creating container: %w", err)
	}

	err = dockerClient.ContainerStart(Ctx, resp.ID, container.StartOptions{})

	if err != nil {
		return "", fmt.Errorf("error starting container: %w", err)
	}

	return resp.ID, nil
}
