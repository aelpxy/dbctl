package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moby/moby/client"
	"golang.org/x/term"
)

// ShellOptions configures an interactive session inside a database container.
type ShellOptions struct {
	Stdin  *os.File
	Stdout io.Writer
	// Plain opens /bin/sh instead of the database client.
	Plain bool
}

// Shell opens the database client, or /bin/sh when the template has none or Plain is set.
func (c *Client) Shell(ctx context.Context, id string, opts ShellOptions) error {
	db, err := c.runningDatabase(ctx, id)
	if err != nil {
		return err
	}

	cmd, err := c.shellCommand(db, opts.Plain)
	if err != nil {
		return err
	}

	width, height, err := term.GetSize(int(opts.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("get terminal size: %w", err)
	}

	exec, err := c.api.ExecCreate(ctx, db.ID, client.ExecCreateOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          true,
		Cmd:          cmd,
		ConsoleSize:  client.ConsoleSize{Height: uint(height), Width: uint(width)},
	})
	if err != nil {
		return fmt.Errorf("create exec in container %s: %w", id, err)
	}

	if err := c.attach(ctx, exec.ID, opts.Stdin, opts.Stdout); err != nil {
		return err
	}

	return c.execResult(ctx, exec.ID)
}

func (c *Client) shellCommand(db *Database, plain bool) ([]string, error) {
	sh := []string{"/bin/sh"}

	if plain {
		return sh, nil
	}

	def, err := c.definition(db.Labels, db.Image)
	if err != nil {
		return nil, err
	}

	client, err := def.Client(db.Password, db.Env)
	if err != nil {
		return nil, fmt.Errorf("render client command: %w", err)
	}

	if len(client) == 0 {
		return sh, nil
	}

	return client, nil
}

func (c *Client) attach(ctx context.Context, execID string, stdin *os.File, stdout io.Writer) error {
	resp, err := c.api.ExecAttach(ctx, execID, client.ExecAttachOptions{TTY: true})
	if err != nil {
		return fmt.Errorf("attach to exec %s: %w", execID, err)
	}

	defer resp.Close()

	stop := context.AfterFunc(ctx, resp.Close)
	defer stop()

	fd := int(stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("set terminal to raw mode: %w", err)
	}

	defer func() { _ = term.Restore(fd, state) }()

	// reading stdin cannot be interrupted, so this goroutine exits on the next keypress or when stdin closes
	go func() {
		_, _ = io.Copy(resp.Conn, stdin)
		_ = resp.CloseWrite()
	}()

	if _, err := io.Copy(stdout, resp.Reader); err != nil && ctx.Err() == nil {
		return fmt.Errorf("read shell output: %w", err)
	}

	return nil
}
