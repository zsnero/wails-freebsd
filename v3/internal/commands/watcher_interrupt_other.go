//go:build !darwin && !linux && !freebsd && !windows

package commands

import "os"

func isInterruptProcessState(*os.ProcessState) bool {
	return false
}
