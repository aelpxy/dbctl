package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func StreamLogs(containerId string, follow bool, tail string) error {
	dockerClient, err := DockerClient()
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	containerInfo, err := InspectContainer(containerId)
	if err != nil {
		return err
	}

	logs, err := dockerClient.ContainerLogs(Ctx, containerInfo.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Tail:       tail,
	})
	if err != nil {
		return fmt.Errorf("error getting container logs: %w", err)
	}
	defer logs.Close()

	// non-tty containers multiplex stdout and stderr into a single stream
	if containerInfo.Config.Tty {
		_, err = io.Copy(os.Stdout, logs)
	} else {
		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, logs)
	}
	if err != nil {
		return fmt.Errorf("error streaming logs: %w", err)
	}

	return nil
}
