//go:build darwin || linux || freebsd

package commands

import (
	"os"
	"syscall"
)

func isInterruptProcessState(state *os.ProcessState) bool {
	status, ok := state.Sys().(syscall.WaitStatus)
	return ok && status.Signal() == syscall.SIGINT
}
