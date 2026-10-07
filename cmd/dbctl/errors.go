package main

import (
	"errors"

	"github.com/aelpxy/dbctl/internal/database"
	"github.com/aelpxy/dbctl/internal/docker"
)

var errNotConfirmed = errors.New("deletion not confirmed")

// printError shows err with a suggested next step when the failure is a known one.
func printError(err error) {
	printErr(errorLine(err.Error()))

	if hint := hintFor(err); hint != "" {
		printErr(hintLine(hint))
	}
}

func hintFor(err error) string {
	switch {
	case errors.Is(err, docker.ErrDaemonUnreachable):
		return "Is Docker running? Check with " + accent("docker info") + "."
	case errors.Is(err, docker.ErrNotFound), errors.Is(err, docker.ErrNotManaged):
		return "Run " + accent("dbctl ls") + " to see your databases."
	case errors.Is(err, docker.ErrNotRunning):
		return "Start it with " + accent("dbctl start <name>") + "."
	case errors.Is(err, docker.ErrUnhealthy):
		return "Check what went wrong with " + accent("dbctl logs <name>") + "."
	case errors.Is(err, database.ErrUnsupported):
		return "Run " + accent("dbctl create --help") + " for the supported databases."
	case errors.Is(err, errNoPassword):
		return "Recreate the database to store its password."
	case errors.Is(err, errMissingArgument):
		return "Pass a database name, or run in a terminal to pick one."
	case errors.Is(err, errNoDatabases):
		return "Create one with " + accent("dbctl create postgres") + "."
	case errors.Is(err, errNotConfirmed):
		return "Pass " + accent("--yes") + " to skip the confirmation."
	default:
		return ""
	}
}
