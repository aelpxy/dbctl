package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func BackupContainer(containerID, outputPath string) (string, error) {
	docker, err := DockerClient()
	if err != nil {
		return "", fmt.Errorf("error creating docker client: %w", err)
	}

	inspectedContainer, err := InspectContainer(containerID)
	if err != nil {
		return "", fmt.Errorf("error inspecting container: %w", err)
	}

	switch inspectedContainer.Config.Image {
	case "postgres":
		var postgresUser, postgresDB string

		//  NOTE: loops through the `Config.Env` and splits them
		for _, envVar := range inspectedContainer.Config.Env {
			pair := strings.SplitN(envVar, "=", 2)
			if len(pair) == 2 {
				switch pair[0] {
				case "POSTGRES_USER":
					postgresUser = pair[1]
				case "POSTGRES_DB":
					postgresDB = pair[1]
				}
			}
		}

		exec, err := docker.ContainerExecCreate(context.Background(), containerID, types.ExecConfig{
			Cmd:          []string{"pg_dump", "-U", postgresUser, "-f", "backup.sql", postgresDB},
			AttachStdout: true,
		})

		if err != nil {
			return "", err
		}

		resp, err := docker.ContainerExecAttach(context.Background(), exec.ID, types.ExecStartCheck{
			Detach: false,
			Tty:    false,
		})

		if err != nil {
			return "", err
		}

		defer resp.Close()

		var backupBytes []byte
		backupBytes, err = io.ReadAll(resp.Reader)

		if err != nil {
			return "", err
		}

		if err := os.WriteFile(outputPath, backupBytes, 0644); err != nil {
			return "", err
		}

		return outputPath, nil
	default:
		return "", fmt.Errorf("unsupported container image")
	}
}
