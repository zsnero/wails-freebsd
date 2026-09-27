//go:build linux && !android

package operatingsystem

import (
	"fmt"
	"os"
)

// platformInfo is the platform specific method to get system information
func platformInfo() (*OS, error) {
	_, err := os.Stat("/etc/os-release")
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("unable to read system information")
	}

	osRelease, _ := os.ReadFile("/etc/os-release")
	return parseOsRelease(string(osRelease)), nil
}
