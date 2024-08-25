//go:build !windows

package main

import (
	"os"
	"syscall"
)

func getShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGKILL}
}
