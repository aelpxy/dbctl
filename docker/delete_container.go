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

	Tcontainer, err := dockerClient.ContainerInspect(context.Background(), containerId)
	if err != nil {
		return fmt.Errorf("error inspecting container: %w", err)
	}

	if !strings.HasPrefix(strings.TrimPrefix(Tcontainer.Name, "/"), config.DockerContainerPrefix) {
		return fmt.Errorf("this container %s is not managed by dbctl", Tcontainer.Name)
	}

	containerSpinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	containerSpinner.Suffix = fmt.Sprintf(" Stopping container %s... \n", strings.TrimPrefix(Tcontainer.Name, "/"))
	containerSpinner.Color("green")
	containerSpinner.Start()

	err = dockerClient.ContainerStop(context.Background(), Tcontainer.ID, container.StopOptions{})

	if err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}

	err = dockerClient.ContainerRemove(context.Background(), Tcontainer.ID, container.RemoveOptions{
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
			_, err := dockerClient.ContainerInspect(context.Background(), Tcontainer.ID)
			if err != nil {
				break
			}
			time.Sleep(1 * time.Second)
		}

		// faily hacky solution
		volumeName := config.DockerVolumeName + Tcontainer.Name[len(fmt.Sprintf("/%s", config.DockerContainerPrefix)):]

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
