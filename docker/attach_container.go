package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func ShellConnect(containerID string) error {
	dockerClient, err := DockerClient()

	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	inspectedContainer, err := InspectContainer(containerID)

	if err != nil {
		return fmt.Errorf("error inspecting container: %w", err)
	}

	out, err := dockerClient.ContainerAttach(context.Background(), containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})

	if err != nil {
		return fmt.Errorf("error attaching to container: %w", err)
	}

	defer out.Close()

	exec, err := dockerClient.ContainerExecCreate(context.Background(), inspectedContainer.ID, types.ExecConfig{
		Cmd:          []string{"sh", "-c", "TERM=xterm-256color; export TERM; exec /bin/sh -i"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})

	if err != nil {
		return fmt.Errorf("error creating exec in container: %w", err)
	}

	execObj, err := dockerClient.ContainerExecAttach(context.Background(), exec.ID, types.ExecStartCheck{})

	if err != nil {
		return fmt.Errorf("error attaching to exec in container: %w", err)
	}

	defer execObj.Close()

	// this basically makes a connection between terminal and the shell session
	go io.Copy(execObj.Conn, os.Stdin)
	_, err = io.Copy(os.Stdout, execObj.Reader)

	if err != nil {
		log.Printf("error copying output: %v", err)
	}

	return nil
}
