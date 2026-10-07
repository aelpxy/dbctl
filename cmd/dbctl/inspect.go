package main

import (
	"context"
	"fmt"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

func (a *app) newInspectCmd() *cobra.Command {
	var output *outputFormat

	cmd := &cobra.Command{
		Use:     "inspect [name-or-id]",
		Short:   "Inspect a database",
		Aliases: []string{"show"},
		Args:    cobra.MaximumNArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			id, err := databaseArg(ctx, c, args, "Inspect which database?")
			if err != nil {
				return err
			}

			return runInspect(ctx, c, id, output)
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

	printOut(accent(db.Name) + "  " + muted(typeName(db.Type)+" · "+shortID(db.ID)))
	printOut("")
	printOut(keyValues([][2]string{
		{"Status", status(db.State)},
		{"Health", health(db.Health)},
		{"Image", db.Image},
		{"CPU", fmt.Sprintf("%.2f%%", stats.CPUPercent)},
		{"Memory", fmt.Sprintf("%.1f MB %s %.1f MB", stats.MemoryUsageMB(), muted("of"), stats.MemoryLimitMB())},
		{"Ports", formatPorts(db.Ports)},
		{"Volumes", formatVolumes(db.Volumes)},
	}))

	return nil
}
