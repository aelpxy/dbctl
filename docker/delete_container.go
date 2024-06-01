package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/config"
	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/container"
)

func DeleteContainer(containerId string, deleteVolume bool) error {
	dockerClient, err := DockerClient()
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	inspectedContainer, err := InspectContainer(containerId)

	if err != nil {
		return fmt.Errorf("error inspecting container: %w", err)
	}

	containerSpinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	containerSpinner.Suffix = fmt.Sprintf(" Stopping container %s... \n", strings.TrimPrefix(inspectedContainer.Name, "/"))
	containerSpinner.Color("green")
	containerSpinner.Start()

	err = dockerClient.ContainerStop(context.Background(), inspectedContainer.ID, container.StopOptions{})

	if err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}

	err = dockerClient.ContainerRemove(context.Background(), inspectedContainer.ID, container.RemoveOptions{
		RemoveVolumes: deleteVolume,
		Force:         true,
	})

	containerSpinner.Stop()

	if err != nil {
		return fmt.Errorf("error removing container: %w", err)
	}

	if deleteVolume {
		// this waits for the container to be fully removed before deleting the volume
		for {
			_, err := dockerClient.ContainerInspect(context.Background(), inspectedContainer.ID)
			if err != nil {
				break
			}
			time.Sleep(1 * time.Second)
		}

		// fairly hacky solution
		volumeName := config.DockerVolumeName + inspectedContainer.Name[len(fmt.Sprintf("/%s", config.DockerContainerPrefix)):]

		fmt.Printf("Deleted volume %s\n", volumeName)

		volumeSpinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		volumeSpinner.Suffix = fmt.Sprintf(" Deleting volume %s... \n", volumeName)
		volumeSpinner.Color("green")
		volumeSpinner.Start()

		err = dockerClient.VolumeRemove(context.Background(), volumeName, true)

		volumeSpinner.Stop()

		if err != nil {
			return fmt.Errorf("error removing volume: %w", err)
		}
	}

	return nil
}
