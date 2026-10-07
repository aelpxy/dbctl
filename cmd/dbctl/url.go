package main

import (
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

func (a *app) newURLCmd() *cobra.Command {
	var output *outputFormat

	cmd := &cobra.Command{
		Use:     "url [name-or-id]",
		Short:   "Print the connection string and credentials of a database",
		Example: "  dbctl url misty-river-bold-pine\n  dbctl url misty-river-bold-pine -o json",
		Aliases: []string{"credentials", "creds"},
		Args:    cobra.MaximumNArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			id, err := databaseArg(ctx, c, args, "Show the connection string of which database?")
			if err != nil {
				return err
			}

			return a.runURL(ctx, c, id, output)
		}),
	}

	output = addOutputFlag(cmd)

	return cmd
}

func (a *app) runURL(ctx context.Context, c *docker.Client, id string, output *outputFormat) error {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return fmt.Errorf("inspect database: %w", err)
	}

	def, err := a.databases.Get(db.Type)
	if err != nil {
		return fmt.Errorf("resolve database type: %w", err)
	}

	conn, err := connectionFor(def, db)
	if err != nil {
		return fmt.Errorf("build connection details for %s: %w", id, err)
	}

	conn.ID = db.ID
	conn.Image = db.Image

	if output.isJSON() {
		return writeJSON(conn)
	}

	fmt.Println(conn.URL)

	return nil
}
