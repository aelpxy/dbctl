//go:build windows

package main

import "os"

// shutdownSignals cancel running commands; windows only delivers interrupts.
func shutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}
