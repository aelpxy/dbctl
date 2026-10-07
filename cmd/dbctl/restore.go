package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

// container id and backup file
const restoreArgCount = 2

func (a *app) newRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "restore <name-or-id> <backup-file>",
		Short:   "Restore a database from a backup file",
		Example: "  dbctl restore misty-river-bold-pine misty-river-bold-pine-20261007-120000.sql",
		Args:    cobra.ExactArgs(restoreArgCount),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			return a.runRestore(ctx, c, args[0], args[1])
		}),
	}
}

func (a *app) runRestore(ctx context.Context, c *docker.Client, id, path string) error {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return fmt.Errorf("inspect database: %w", err)
	}

	def, err := a.databases.Get(db.Type)
	if err != nil {
		return fmt.Errorf("resolve database type: %w", err)
	}

	if !def.CanRestore() {
		return fmt.Errorf("restore database %s: %w", id, database.ErrRestoreUnsupported)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open backup file: %w", err)
	}

	defer func() { _ = f.Close() }()

	err = step("Restored "+bold(path)+" into "+bold(db.Name), func() error { return c.Restore(ctx, db.ID, f) })
	if err != nil {
		return fmt.Errorf("restore database %s: %w", id, err)
	}

	return nil
}
