package main

import (
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/olekukonko/tablewriter"
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
		fmt.Println("No databases found.")

		return nil
	}

	table := newTable("ID", "Name", "Type", "Image", "Status")

	for i := range dbs {
		db := &dbs[i]
		table.Rich(
			[]string{shortID(db.ID), db.Name, typeName(db.Type), db.Image, db.Status},
			[]tablewriter.Colors{
				{tablewriter.FgGreenColor, tablewriter.Bold},
				{tablewriter.FgBlueColor, tablewriter.Bold},
				{},
				{tablewriter.FgGreenColor},
				statusColor(db.State),
			},
		)
	}

	table.Render()

	return nil
}
