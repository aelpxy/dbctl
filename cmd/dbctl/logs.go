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
		Use:     "logs <container-id>",
		Short:   "Stream live logs of a database",
		Example: "  dbctl logs container-id\n  dbctl logs container-id --tail 100 --follow=false",
		Aliases: []string{"tail"},
		Args:    cobra.ExactArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			if err := c.Logs(ctx, args[0], opts); err != nil {
				return fmt.Errorf("stream database logs: %w", err)
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&opts.Follow, "follow", "f", true, "Follow log output.")
	cmd.Flags().StringVarP(&opts.Tail, "tail", "t", "all", "Number of lines to show from the end of the logs.")

	return cmd
}
