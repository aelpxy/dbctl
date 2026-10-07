package docker

import (
	"fmt"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/config"
	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func DeleteContainer(containerId string, deleteVolume bool) error {
	dockerClient, err := DockerClient()
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	inspectedContainer, err := InspectContainer(containerId)
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Deleting container %s...", strings.TrimPrefix(inspectedContainer.Name, "/"))
	s.Color("green")
	s.Start()

	err = dockerClient.ContainerRemove(Ctx, inspectedContainer.ID, container.RemoveOptions{Force: true})

	s.Stop()

	if err != nil {
		return fmt.Errorf("error removing container: %w", err)
	}

	if !deleteVolume {
		return nil
	}

	for _, m := range inspectedContainer.Mounts {
		if m.Type != mount.TypeVolume || !strings.HasPrefix(m.Name, config.DockerVolumeName) {
			continue
		}

		if err := dockerClient.VolumeRemove(Ctx, m.Name, true); err != nil {
			return fmt.Errorf("error removing volume %s: %w", m.Name, err)
		}

		fmt.Printf("Deleted volume %s\n", m.Name)
	}

	return nil
}
