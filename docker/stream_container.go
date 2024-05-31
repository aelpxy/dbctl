package docker

import (
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
)

func StreamLogs(containerId string) {
	dockerClient, err := DockerClient()
	if err != nil {
		log.Fatalf("error creating docker client: %v", err)
	}

	containerInfo, err := InspectContainer(containerId)

	if err != nil {
		log.Fatalf("error inspecting container: %v", err)
	}

	logOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	}

	logs, err := dockerClient.ContainerLogs(Ctx, containerInfo.ID, logOptions)

	if err != nil {
		log.Fatalf("error getting container logs: %v", err)
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)

	if err != nil {
		log.Fatalf("error streaming logs: %v", err)
	}
}
