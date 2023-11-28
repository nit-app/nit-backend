//go:build !windows
// +build !windows

package main

import "os"

func getShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGKILL}
}
