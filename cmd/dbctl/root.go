package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

const version = "2.0.0"

type app struct {
	databases *database.Registry
}

func (a *app) newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "dbctl",
		Short:         "A CLI tool for managing containerized databases",
		Long:          "A command-line tool designed to simplify the management of databases, including creating, deleting, and other operations.",
		SuggestFor:    []string{"db"},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		a.newBackupCmd(),
		a.newCreateCmd(),
		a.newDeleteCmd(),
		a.newHTTPCmd(),
		a.newInspectCmd(),
		a.newListCmd(),
		a.newLogsCmd(),
		a.newRestartCmd(),
		a.newRestoreCmd(),
		a.newShellCmd(),
		a.newStartCmd(),
		a.newStopCmd(),
		a.newURLCmd(),
		newVersionCmd(),
	)

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the current dbctl version",
		Args:  cobra.NoArgs,
		Run: func(*cobra.Command, []string) {
			fmt.Println(version)
		},
	}
}

// withDocker connects to the docker daemon for the duration of a command.
func (a *app) withDocker(fn func(ctx context.Context, c *docker.Client, args []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()

		c, err := docker.New(ctx, a.databases)
		if err != nil {
			return fmt.Errorf("connect to docker: %w", err)
		}

		defer func() { err = errors.Join(err, c.Close()) }()

		return fn(ctx, c, args)
	}
}
