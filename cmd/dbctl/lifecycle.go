package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

type lifecycleAction struct {
	run   func(ctx context.Context, c *docker.Client, id string) error
	use   string
	short string
	done  string
}

func (a *app) newStartCmd() *cobra.Command {
	return a.newLifecycleCmd(&lifecycleAction{
		run:   func(ctx context.Context, c *docker.Client, id string) error { return c.Start(ctx, id) },
		use:   "start",
		short: "Start one or more stopped databases",
		done:  "Started",
	})
}

func (a *app) newStopCmd() *cobra.Command {
	return a.newLifecycleCmd(&lifecycleAction{
		run:   func(ctx context.Context, c *docker.Client, id string) error { return c.Stop(ctx, id) },
		use:   "stop",
		short: "Stop one or more databases without deleting them",
		done:  "Stopped",
	})
}

func (a *app) newRestartCmd() *cobra.Command {
	return a.newLifecycleCmd(&lifecycleAction{
		run:   func(ctx context.Context, c *docker.Client, id string) error { return c.Restart(ctx, id) },
		use:   "restart",
		short: "Restart one or more databases",
		done:  "Restarted",
	})
}

func (a *app) newLifecycleCmd(action *lifecycleAction) *cobra.Command {
	return &cobra.Command{
		Use:               action.use + " [name-or-id...]",
		Short:             action.short,
		ValidArgsFunction: a.completeDatabases(true),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			if len(args) == 0 {
				id, err := databaseArg(ctx, c, args, strings.ToUpper(action.use[:1])+action.use[1:]+" which database?")
				if err != nil {
					return err
				}

				args = []string{id}
			}

			for _, id := range args {
				if err := runLifecycle(ctx, c, id, action); err != nil {
					return err
				}
			}

			return nil
		}),
	}
}

func runLifecycle(ctx context.Context, c *docker.Client, id string, action *lifecycleAction) error {
	err := step(action.done+" "+bold(id), func() error { return action.run(ctx, c, id) })
	if err != nil {
		return fmt.Errorf("%s database %s: %w", action.use, id, err)
	}

	return nil
}
