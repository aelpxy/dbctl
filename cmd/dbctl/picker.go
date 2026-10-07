package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"charm.land/huh/v2"
	"github.com/aelpxy/dbctl/internal/docker"
	"golang.org/x/term"
)

var (
	errMissingArgument = errors.New("no database given")
	errNoDatabases     = errors.New("no databases to choose from")
	errPickCanceled    = errors.New("selection canceled")
)

// databaseArg returns the database named in args, or lets the user pick one interactively.
func databaseArg(ctx context.Context, c *docker.Client, args []string, title string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	dbs, err := c.List(ctx)
	if err != nil {
		return "", fmt.Errorf("list databases: %w", err)
	}

	if len(dbs) == 0 {
		return "", errNoDatabases
	}

	options := make([]huh.Option[string], 0, len(dbs))

	for i := range dbs {
		db := &dbs[i]
		options = append(options, huh.NewOption(db.Name+"  ("+typeName(db.Type)+", "+db.State+")", db.Name))
	}

	return pick(ctx, title, options)
}

// typeArg returns the database type in args, or lets the user pick one interactively.
func (a *app) typeArg(ctx context.Context, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	names := a.databases.Names()
	options := make([]huh.Option[string], 0, len(names))

	for _, name := range names {
		options = append(options, huh.NewOption(name, name))
	}

	return pick(ctx, "Which database do you want to create?", options)
}

// pick shows a filterable list on stderr; without a terminal it fails instead of hanging a script.
func pick(ctx context.Context, title string, options []huh.Option[string]) (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stderr.Fd())) {
		return "", errMissingArgument
	}

	var choice string

	field := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Filtering(true).
		Value(&choice)

	err := huh.NewForm(huh.NewGroup(field)).
		WithTheme(huh.ThemeFunc(huh.ThemeCharm)).
		WithOutput(os.Stderr).
		RunWithContext(ctx)
	if errors.Is(err, huh.ErrUserAborted) || ctx.Err() != nil {
		return "", errPickCanceled
	}

	if err != nil {
		return "", fmt.Errorf("pick an option: %w", err)
	}

	return choice, nil
}
