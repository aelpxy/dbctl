package main

import (
	"cmp"
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func (a *app) newInspectCmd() *cobra.Command {
	var output *outputFormat

	cmd := &cobra.Command{
		Use:     "inspect <container-id>",
		Short:   "Inspect a database",
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			return runInspect(ctx, c, args[0], output)
		}),
	}

	output = addOutputFlag(cmd)

	return cmd
}

func runInspect(ctx context.Context, c *docker.Client, id string, output *outputFormat) error {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return fmt.Errorf("inspect database: %w", err)
	}

	stats, err := c.Stats(ctx, db.ID)
	if err != nil {
		return fmt.Errorf("get database stats: %w", err)
	}

	if output.isJSON() {
		out := newDatabaseJSON(db)
		out.Stats = newStatsJSON(stats)

		return writeJSON(out)
	}

	table := newTable("Field", "Value")
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
	)
	table.AppendBulk([][]string{
		{"ID", shortID(db.ID)},
		{"Name", db.Name},
		{"Type", typeName(db.Type)},
		{"Image", db.Image},
		{"Status", db.State},
		{"Health", cmp.Or(db.Health, "none")},
		{"CPU Usage", fmt.Sprintf("%.2f%%", stats.CPUPercent)},
		{"Memory Usage", fmt.Sprintf("%.2f MB / %.2f MB", stats.MemoryUsageMB(), stats.MemoryLimitMB())},
		{"Ports", formatPorts(db.Ports)},
		{"Volumes", formatVolumes(db.Volumes)},
	})
	table.Render()

	return nil
}
