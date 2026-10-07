// Command dbctl manages containerized databases.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/aelpxy/dbctl/internal/database"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), shutdownSignals()...)
	defer stop()

	databases, err := database.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)

		return 1
	}

	a := &app{databases: databases}

	if err := a.newRootCmd().ExecuteContext(ctx); err != nil {
		return 1
	}

	return 0
}
