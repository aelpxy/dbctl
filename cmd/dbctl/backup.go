package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

const backupTimeFormat = "20060102-150405"

func (a *app) newBackupCmd() *cobra.Command {
	output := ""

	cmd := &cobra.Command{
		Use:               "backup [name-or-id]",
		Short:             "Backup a database",
		Example:           "  dbctl backup misty-river-bold-pine\n  dbctl backup misty-river-bold-pine -o backup.sql",
		Aliases:           []string{"cp"},
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: a.completeDatabases(false),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			id, err := databaseArg(ctx, c, args, "Back up which database?")
			if err != nil {
				return err
			}

			path, err := a.runBackup(ctx, c, id, output)
			if err != nil {
				return err
			}

			printOut(hintLine("Restore it with " + accent("dbctl restore "+id+" "+path)))

			return nil
		}),
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Specify an output path for the database backup (default: <name>-<timestamp>.<ext>).")

	return cmd
}

func (a *app) runBackup(ctx context.Context, c *docker.Client, id, output string) (_ string, err error) {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return "", fmt.Errorf("inspect database: %w", err)
	}

	def, err := a.databases.Get(db.Type)
	if err != nil {
		return "", fmt.Errorf("resolve database type: %w", err)
	}

	if !def.CanBackup() {
		return "", fmt.Errorf("back up database %s: %w", id, database.ErrBackupUnsupported)
	}

	path := cmp.Or(output, backupFileName(db, def))

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create backup file: %w", err)
	}

	defer func() {
		err = errors.Join(err, f.Close())
		if err != nil {
			// a partial dump is worse than none, the original error is already being returned
			_ = os.Remove(path)
		}
	}()

	err = step("Backed up "+bold(db.Name)+" to "+bold(path), func() error { return c.Backup(ctx, db.ID, f) })
	if err != nil {
		return "", fmt.Errorf("back up database %s: %w", id, err)
	}

	return path, nil
}

func backupFileName(db *docker.Database, def *database.Definition) string {
	return db.Name + "-" + time.Now().Format(backupTimeFormat) + "." + def.BackupExtension()
}
