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

func CreateContainer(imageName, dbType, containerName string, port int, password string, envVars []string) (string, error) {
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

	containerConfig := &container.Config{
		Image: imageName,
		Env:   envVars,
		ExposedPorts: nat.PortSet{
			nat.Port(strconv.Itoa(port) + "/tcp"): struct{}{},
		},
	}

	hostIP := utils.GetIP().String()

	var mountSource string
	var mountTarget string

	if dbType == "postgres" {
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/var/lib/postgresql/data"
		containerConfig.Env = append(containerConfig.Env,
			"POSTGRES_DB=postgres",
			"POSTGRES_USER=postgres",
		)
	} else {
		mountSource = config.DockerVolumeName + containerName
		mountTarget = "/data"
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
			nat.Port(strconv.Itoa(port) + "/tcp"): []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: strconv.Itoa(port),
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
