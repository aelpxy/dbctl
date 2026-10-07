package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

// completion runs on every tab press, so a slow or missing docker must not hang the shell
const completionTimeout = 2 * time.Second

var errNoContext = errors.New("command has no context")

type completionFunc = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// completeDatabases suggests database names; multiple allows more than one, skipping names already given.
func (a *app) completeDatabases(multiple bool) completionFunc {
	return func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if !multiple && len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		names, err := a.databaseCompletions(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		names = slices.DeleteFunc(names, func(entry string) bool {
			name, _, _ := strings.Cut(entry, "\t")

			return slices.Contains(args, name)
		})

		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

// completeRestore suggests a database first, then falls back to file completion for the backup.
func (a *app) completeRestore(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return a.completeDatabases(false)(cmd, args, toComplete)
	}

	return nil, cobra.ShellCompDirectiveDefault
}

// databaseCompletions returns "name\tdescription" entries, which shells show as name plus a hint.
func (a *app) databaseCompletions(ctx context.Context) ([]string, error) {
	if ctx == nil {
		return nil, errNoContext
	}

	ctx, cancel := context.WithTimeout(ctx, completionTimeout)
	defer cancel()

	c, err := docker.New(ctx, a.databases)
	if err != nil {
		return nil, fmt.Errorf("connect to docker: %w", err)
	}

	defer func() { _ = c.Close() }()

	dbs, err := c.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list databases: %w", err)
	}

	names := make([]string, 0, len(dbs))

	for i := range dbs {
		db := &dbs[i]
		names = append(names, db.Name+"\t"+typeName(db.Type)+" · "+db.State)
	}

	return names, nil
}
