package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/term"
)

func ShellConnect(containerID string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dockerClient, err := DockerClient()
	if err != nil {
		return fmt.Errorf("error creating docker client: %w", err)
	}

	inspectedContainer, err := InspectContainer(containerID)
	if err != nil {
		return fmt.Errorf("error inspecting container: %w", err)
	}
	if !inspectedContainer.State.Running {
		return fmt.Errorf("container %s is not running", containerID)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("error setting terminal to raw mode: %w", err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("error getting terminal size: %w", err)
	}

	consoleSize := [2]uint{uint(height), uint(width)}

	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/sh"},
		Detach:       false,
		ConsoleSize:  &consoleSize,
	}

	exec, err := dockerClient.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return fmt.Errorf("error creating exec in container: %w", err)
	}

	hijackedResp, err := dockerClient.ContainerExecAttach(ctx, exec.ID, container.ExecStartOptions{
		Detach: false,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("error attaching to exec in container: %w", err)
	}

	defer func() {
		hijackedResp.CloseWrite()
		hijackedResp.Close()
	}()

	shutdownChan := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		select {
		case <-sigChan:
			cancel()
			close(shutdownChan)
		case <-ctx.Done():
			close(shutdownChan)
		}
	}()
	defer signal.Stop(sigChan)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer cancel()

		buf := make([]byte, 32*1024)
		for {
			nr, err := os.Stdin.Read(buf)
			if err == io.EOF || (nr == 5 && string(buf[:nr]) == "exit\n") {
				return
			}
			if err != nil {
				if err != io.ErrClosedPipe {
					fmt.Fprintf(os.Stderr, "stdin error: %v\n", err)
				}
				return
			}
			if nr > 0 {
				_, err := hijackedResp.Conn.Write(buf[:nr])
				if err != nil {
					if err != io.ErrClosedPipe {
						fmt.Fprintf(os.Stderr, "write error: %v\n", err)
					}
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer cancel()

		if execConfig.Tty {
			_, err = io.Copy(os.Stdout, hijackedResp.Reader)
		} else {
			_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, hijackedResp.Reader)
		}
		if err != nil && err != io.EOF && err != io.ErrClosedPipe {
			fmt.Fprintf(os.Stderr, "output error: %v\n", err)
		}
	}()

	<-shutdownChan

	cleanupDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(cleanupDone)
	}()

	select {
	case <-cleanupDone:
	case <-time.After(500 * time.Millisecond):
	}

	select {
	case <-ctx.Done():
		return nil
	default:
		inspectResp, err := dockerClient.ContainerExecInspect(context.Background(), exec.ID)

		if err != nil {
			return nil
		}
		if inspectResp.ExitCode != 0 && inspectResp.ExitCode != -1 {
			return fmt.Errorf("command exited with code %d", inspectResp.ExitCode)
		}
	}

	return nil
}
