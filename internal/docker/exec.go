package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

// execOptions wires a non-interactive command to the caller's streams.
type execOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
	Exec   database.Exec
}

// exec runs a command in the container without a tty, streaming stdin in and stdout out.
func (c *Client) exec(ctx context.Context, containerID string, opts *execOptions) error {
	created, err := c.api.ExecCreate(ctx, containerID, client.ExecCreateOptions{
		Cmd:          opts.Exec.Command,
		Env:          opts.Exec.Env,
		AttachStdin:  opts.Stdin != nil,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return fmt.Errorf("create exec in container %s: %w", containerID, err)
	}

	resp, err := c.api.ExecAttach(ctx, created.ID, client.ExecAttachOptions{})
	if err != nil {
		return fmt.Errorf("attach to exec %s: %w", created.ID, err)
	}

	defer resp.Close()

	var (
		wg      sync.WaitGroup
		copyErr error
		stderr  bytes.Buffer
	)

	if opts.Stdin != nil {
		wg.Go(func() {
			_, copyErr = io.Copy(resp.Conn, opts.Stdin)
			copyErr = errors.Join(copyErr, resp.CloseWrite())
		})
	}

	_, err = stdcopy.StdCopy(opts.Stdout, &stderr, resp.Reader)

	wg.Wait()

	if err := errors.Join(err, copyErr); err != nil {
		return fmt.Errorf("stream exec: %w", err)
	}

	if err := c.execResult(ctx, created.ID); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return nil
}

// ExitError reports a non-zero exit code from a command run inside a database container.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("command exited with code %d", e.Code)
}

// the exit code is still read after ctx is canceled, e.g. when a shell is interrupted
func (c *Client) execResult(ctx context.Context, execID string) error {
	result, err := c.api.ExecInspect(context.WithoutCancel(ctx), execID, client.ExecInspectOptions{})
	if err != nil {
		return fmt.Errorf("inspect exec %s: %w", execID, err)
	}

	if result.ExitCode != 0 {
		return &ExitError{Code: result.ExitCode}
	}

	return nil
}

func envMap(vars []string) map[string]string {
	env := make(map[string]string, len(vars))

	for _, v := range vars {
		if key, value, ok := strings.Cut(v, "="); ok {
			env[key] = value
		}
	}

	return env
}
