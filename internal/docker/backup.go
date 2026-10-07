package docker

import (
	"context"
	"fmt"
	"io"
)

// Backup dumps the database to w.
func (c *Client) Backup(ctx context.Context, id string, w io.Writer) error {
	db, err := c.runningDatabase(ctx, id)
	if err != nil {
		return err
	}

	def, err := c.definition(db.Labels, db.Image)
	if err != nil {
		return err
	}

	backup, err := def.Backup(db.Env)
	if err != nil {
		return fmt.Errorf("render backup command: %w", err)
	}

	return c.exec(ctx, db.ID, &execOptions{Stdout: w, Exec: backup})
}

// Restore loads a dump read from r into the database.
func (c *Client) Restore(ctx context.Context, id string, r io.Reader) error {
	db, err := c.runningDatabase(ctx, id)
	if err != nil {
		return err
	}

	def, err := c.definition(db.Labels, db.Image)
	if err != nil {
		return err
	}

	restore, err := def.Restore(db.Env)
	if err != nil {
		return fmt.Errorf("render restore command: %w", err)
	}

	return c.exec(ctx, db.ID, &execOptions{Stdin: r, Stdout: io.Discard, Exec: restore})
}

func (c *Client) runningDatabase(ctx context.Context, id string) (*Database, error) {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return nil, err
	}

	if !db.Running {
		return nil, fmt.Errorf("%w: %s", ErrNotRunning, id)
	}

	return db, nil
}
