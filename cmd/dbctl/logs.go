package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

func (a *app) newLogsCmd() *cobra.Command {
	opts := docker.LogsOptions{Stdout: os.Stdout, Stderr: os.Stderr}

	cmd := &cobra.Command{
		Use:     "logs [name-or-id]",
		Short:   "Stream live logs of a database",
		Example: "  dbctl logs misty-river-bold-pine\n  dbctl logs misty-river-bold-pine --tail 100 --follow=false",
		Aliases: []string{"tail"},
		Args:    cobra.MaximumNArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			id, err := databaseArg(ctx, c, args, "Show logs of which database?")
			if err != nil {
				return err
			}

			if err := c.Logs(ctx, id, opts); err != nil {
				return fmt.Errorf("stream database logs: %w", err)
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&opts.Follow, "follow", "f", true, "Follow log output.")
	cmd.Flags().StringVarP(&opts.Tail, "tail", "t", "all", "Number of lines to show from the end of the logs.")

	return cmd
}
