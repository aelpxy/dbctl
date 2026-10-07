package main

import (
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	removeVolumes bool
	yes           bool
}

func (a *app) newDeleteCmd() *cobra.Command {
	opts := deleteOptions{}

	cmd := &cobra.Command{
		Use:               "delete [name-or-id...]",
		Short:             "Stop and delete one or more databases",
		Long:              "Stop and remove databases along with their data volumes (use --force=false to keep the volumes).",
		Example:           "  dbctl delete misty-river-bold-pine\n  dbctl delete misty-river-bold-pine --yes",
		Aliases:           []string{"rm"},
		ValidArgsFunction: a.completeDatabases(true),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			if len(args) == 0 {
				id, err := databaseArg(ctx, c, args, "Delete which database?")
				if err != nil {
					return err
				}

				args = []string{id}
			}

			for _, id := range args {
				if err := deleteDatabase(ctx, c, id, &opts); err != nil {
					return err
				}
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&opts.removeVolumes, "force", "f", true, "Delete the associated data volume.")
	cmd.Flags().BoolVarP(&opts.yes, "yes", "y", false, "Skip the confirmation prompt.")

	return cmd
}

func deleteDatabase(ctx context.Context, c *docker.Client, id string, opts *deleteOptions) error {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return fmt.Errorf("find database: %w", err)
	}

	if !opts.yes {
		if err := confirmDelete(ctx, db.Name, opts.removeVolumes); err != nil {
			return err
		}
	}

	var volumes []string

	err = step("Deleted "+bold(db.Name), func() error {
		var err error

		volumes, err = c.Delete(ctx, db.ID, docker.DeleteOptions{KeepVolumes: !opts.removeVolumes})
		if err != nil {
			return fmt.Errorf("delete database %s: %w", db.Name, err)
		}

		return nil
	})

	for _, v := range volumes {
		printErr(successLine("Removed volume " + muted(v)))
	}

	return err
}
