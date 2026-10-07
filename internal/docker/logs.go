package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

// LogsOptions configures where and how database logs are streamed.
type LogsOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	Tail   string
	Follow bool
}

// Logs streams the database logs until they end or ctx is canceled.
func (c *Client) Logs(ctx context.Context, id string, opts LogsOptions) error {
	info, err := c.inspect(ctx, id)
	if err != nil {
		return err
	}

	logs, err := c.api.ContainerLogs(ctx, info.ID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     opts.Follow,
		Tail:       opts.Tail,
	})
	if err != nil {
		return fmt.Errorf("get container logs %s: %w", id, err)
	}

	defer func() { _ = logs.Close() }()

	// dbctl containers run without a tty so stdout and stderr are multiplexed into one stream
	_, err = stdcopy.StdCopy(opts.Stdout, opts.Stderr, logs)
	if err != nil && ctx.Err() == nil {
		return fmt.Errorf("stream logs: %w", err)
	}

	return nil
}
