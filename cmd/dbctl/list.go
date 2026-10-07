package main

import (
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

func (a *app) newListCmd() *cobra.Command {
	var output *outputFormat

	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List all databases",
		Aliases: []string{"list"},
		Args:    cobra.NoArgs,
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, _ []string) error {
			return runList(ctx, c, output)
		}),
	}

	output = addOutputFlag(cmd)

	return cmd
}

func runList(ctx context.Context, c *docker.Client, output *outputFormat) error {
	dbs, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("list databases: %w", err)
	}

	if output.isJSON() {
		out := make([]databaseJSON, 0, len(dbs))

		for i := range dbs {
			out = append(out, newDatabaseJSON(&dbs[i]))
		}

		return writeJSON(out)
	}

	if len(dbs) == 0 {
		printOut("No databases yet.")
		printOut(hintLine("Create one with " + accent("dbctl create postgres")))

		return nil
	}

	t := newTable("Name", "Type", "Status", "Uptime", "ID")

	for i := range dbs {
		db := &dbs[i]
		t.Row(bold(db.Name), typeName(db.Type), status(db.State), muted(db.Status), muted(shortID(db.ID)))
	}

	printOut(t.Render())

	return nil
}
