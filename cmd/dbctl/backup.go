package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

const backupTimeFormat = "20060102-150405"

func (a *app) newBackupCmd() *cobra.Command {
	output := ""

	cmd := &cobra.Command{
		Use:     "backup <container-id>",
		Short:   "Backup a database",
		Example: "  dbctl backup container-id\n  dbctl backup container-id -o backup.sql",
		Aliases: []string{"cp"},
		Args:    cobra.ExactArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			path, err := a.runBackup(ctx, c, args[0], output)
			if err != nil {
				return err
			}

			fmt.Printf("Backup saved to %s\n", path)

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

	if err := c.Backup(ctx, db.ID, f); err != nil {
		return "", fmt.Errorf("back up database %s: %w", id, err)
	}

	return path, nil
}

func backupFileName(db *docker.Database, def *database.Definition) string {
	name := strings.TrimPrefix(db.Name, docker.ContainerPrefix)

	return name + "-" + time.Now().Format(backupTimeFormat) + "." + def.BackupExtension()
}
