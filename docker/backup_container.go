package docker

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func BackupContainer(containerID, outputPath string) (string, error) {
	dockerClient, err := DockerClient()
	if err != nil {
		return "", fmt.Errorf("error creating docker client: %w", err)
	}

	inspectedContainer, err := InspectContainer(containerID)
	if err != nil {
		return "", err
	}

	if !inspectedContainer.State.Running {
		return "", fmt.Errorf("container %s is not running", containerID)
	}

	env := parseEnv(inspectedContainer.Config.Env)

	var cmd, execEnv []string
	extension := "sql"

	switch databaseType(inspectedContainer.Config) {
	case "postgres":
		user := valueOr(env["POSTGRES_USER"], "postgres")
		cmd = []string{"pg_dump", "-U", user, valueOr(env["POSTGRES_DB"], user)}
	case "mysql":
		cmd = []string{"mysqldump", "-uroot", "--all-databases"}
		execEnv = []string{"MYSQL_PWD=" + env["MYSQL_ROOT_PASSWORD"]}
	case "mariadb":
		cmd = []string{"mariadb-dump", "-uroot", "--all-databases"}
		execEnv = []string{"MYSQL_PWD=" + env["MARIADB_ROOT_PASSWORD"]}
	case "mongo":
		cmd = []string{"mongodump", "--archive", "--username", env["MONGO_INITDB_ROOT_USERNAME"], "--password", env["MONGO_INITDB_ROOT_PASSWORD"], "--authenticationDatabase", "admin"}
		extension = "archive"
	default:
		return "", fmt.Errorf("backups are not supported for image %s", inspectedContainer.Config.Image)
	}

	if outputPath == "" {
		name := strings.TrimPrefix(strings.TrimPrefix(inspectedContainer.Name, "/"), config.DockerContainerPrefix)
		outputPath = fmt.Sprintf("%s-%s.%s", name, time.Now().Format("20060102-150405"), extension)
	}

	exec, err := dockerClient.ContainerExecCreate(Ctx, inspectedContainer.ID, container.ExecOptions{
		Cmd:          cmd,
		Env:          execEnv,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("error creating backup exec: %w", err)
	}

	resp, err := dockerClient.ContainerExecAttach(Ctx, exec.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("error attaching to backup exec: %w", err)
	}
	defer resp.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("error creating backup file: %w", err)
	}
	defer file.Close()

	var stderr bytes.Buffer

	if _, err := stdcopy.StdCopy(file, &stderr, resp.Reader); err != nil {
		os.Remove(outputPath)
		return "", fmt.Errorf("error writing backup: %w", err)
	}

	inspectResp, err := dockerClient.ContainerExecInspect(Ctx, exec.ID)
	if err != nil {
		os.Remove(outputPath)
		return "", fmt.Errorf("error inspecting backup exec: %w", err)
	}

	if inspectResp.ExitCode != 0 {
		os.Remove(outputPath)
		return "", fmt.Errorf("backup command exited with code %d: %s", inspectResp.ExitCode, strings.TrimSpace(stderr.String()))
	}

	return outputPath, nil
}

// falls back to the image name for containers created before the type label existed
func databaseType(cfg *container.Config) string {
	if dbType, ok := cfg.Labels[config.DockerTypeLabel]; ok {
		return dbType
	}

	repository, _, _ := strings.Cut(cfg.Image, ":")

	return repository[strings.LastIndex(repository, "/")+1:]
}

func parseEnv(vars []string) map[string]string {
	env := make(map[string]string, len(vars))

	for _, v := range vars {
		if key, value, ok := strings.Cut(v, "="); ok {
			env[key] = value
		}
	}

	return env
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
