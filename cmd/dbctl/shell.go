package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

func (a *app) newShellCmd() *cobra.Command {
	opts := docker.ShellOptions{Stdin: os.Stdin, Stdout: os.Stdout}

	cmd := &cobra.Command{
		Use:     "shell [name-or-id]",
		Short:   "Open the database client, or a shell with --sh",
		Example: "  dbctl shell misty-river-bold-pine\n  dbctl shell misty-river-bold-pine --sh",
		Aliases: []string{"enter", "sh", "connect"},
		Args:    cobra.MaximumNArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			id, err := databaseArg(ctx, c, args, "Open a client for which database?")
			if err != nil {
				return err
			}

			if err := c.Shell(ctx, id, opts); err != nil {
				return fmt.Errorf("open shell: %w", err)
			}

			return nil
		}),
	}

	cmd.Flags().BoolVar(&opts.Plain, "sh", false, "Open /bin/sh instead of the database client.")

	return cmd
}
