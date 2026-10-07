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
		Use:     "shell <container-id>",
		Short:   "Open the database client, or a shell with --sh",
		Example: "  dbctl shell container-id\n  dbctl shell container-id --sh",
		Aliases: []string{"enter", "sh", "connect"},
		Args:    cobra.ExactArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			if err := c.Shell(ctx, args[0], opts); err != nil {
				return fmt.Errorf("open shell: %w", err)
			}

			return nil
		}),
	}

	cmd.Flags().BoolVar(&opts.Plain, "sh", false, "Open /bin/sh instead of the database client.")

	return cmd
}
