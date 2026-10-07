//go:build !windows

package main

import (
	"os"

	"golang.org/x/sys/unix"
)

// shutdownSignals cancel running commands.
func shutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, unix.SIGTERM}
}
