package main

import (
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

func (a *app) newDeleteCmd() *cobra.Command {
	removeVolumes := true

	cmd := &cobra.Command{
		Use:     "delete <container-id>...",
		Short:   "Stop and delete one or more databases",
		Long:    "Stop and remove database containers along with their volumes (use --force=false to keep the volumes).",
		Aliases: []string{"rm"},
		Args:    cobra.MinimumNArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			opts := docker.DeleteOptions{KeepVolumes: !removeVolumes}

			for _, id := range args {
				if err := deleteDatabase(ctx, c, id, opts); err != nil {
					return err
				}
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&removeVolumes, "force", "f", true, "Delete the associated volume.")

	return cmd
}

func deleteDatabase(ctx context.Context, c *docker.Client, id string, opts docker.DeleteOptions) error {
	s := startSpinner("Deleting database " + id + "...")
	volumes, err := c.Delete(ctx, id, opts)

	s.Stop()

	for _, v := range volumes {
		fmt.Printf("Deleted volume %s\n", v)
	}

	if err != nil {
		return fmt.Errorf("delete database %s: %w", id, err)
	}

	fmt.Printf("Database %s has been deleted.\n", id)

	return nil
}
