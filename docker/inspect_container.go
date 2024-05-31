package docker

import (
	"fmt"
	"strings"

	"github.com/aelpxy/dbctl/config"
	"github.com/docker/docker/api/types"
)

func InspectContainer(containerId string) (types.ContainerJSON, error) {
	dockerClient, err := DockerClient()

	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("error creating docker client: %w", err)
	}

	containerInfo, err := dockerClient.ContainerInspect(Ctx, containerId)

	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("error inspecting container: %w", err)
	}

	if !strings.HasPrefix(strings.TrimPrefix(containerInfo.Name, "/"), config.DockerContainerPrefix) {
		return types.ContainerJSON{}, fmt.Errorf("this container %s is not managed by dbctl", containerInfo.Name)
	}

	return containerInfo, nil
}
