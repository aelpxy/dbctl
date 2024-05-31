package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

func DeleteContainer(containerID string, deleteVolume bool) error {
	dockerClient, err := DockerClient()
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	Tcontainer, err := dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return fmt.Errorf("error inspecting container: %w", err)
	}

	if !strings.HasPrefix(strings.TrimPrefix(Tcontainer.Name, "/"), config.DockerContainerPrefix) {
		return fmt.Errorf("this container %s is not managed by dbctl", Tcontainer.Name)
	}

	err = dockerClient.ContainerStop(Ctx, containerID, container.StopOptions{})

	if err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}

	err = dockerClient.ContainerRemove(Ctx, containerID, container.RemoveOptions{
		RemoveVolumes: deleteVolume,
		Force:         true,
	})
	if err != nil {
		return fmt.Errorf("error removing container: %w", err)
	}

	if deleteVolume {
		// this waits for the container to be fully removed before deleting the volume
		for {
			_, err := dockerClient.ContainerInspect(Ctx, containerID)
			if err != nil {
				break
			}
			time.Sleep(1 * time.Second)
		}

		filter := filters.NewArgs()
		filter.Add("name", config.DockerVolumeName+containerID)
		volumes, err := dockerClient.VolumeList(Ctx, volume.ListOptions{Filters: filter})

		if err != nil {
			return fmt.Errorf("error listing volumes: %w", err)
		}

		for _, volume := range volumes.Volumes {
			err = dockerClient.VolumeRemove(Ctx, volume.Name, true)
			if err != nil {
				return fmt.Errorf("error removing volume: %w", err)
			}
		}
	}

	return nil
}
