//go:build freebsd

package operatingsystem

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// platformInfo reports the host FreeBSD release.
//
// FreeBSD 11 and later ship a freedesktop os-release file, which is the cheap
// path: it is already parsed for us. Trimmed images and jail base systems can be
// built without one, so fall back to uname(2) rather than failing — the caller
// treats an error here as "no system information", which is a worse outcome
// than a slightly coarser version string.
func platformInfo() (*OS, error) {
	if osRelease, err := os.ReadFile("/etc/os-release"); err == nil {
		return parseOsRelease(string(osRelease)), nil
	}

	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return nil, fmt.Errorf("unable to read system information: %w", err)
	}

	// Release looks like "15.1-RELEASE-p3". Split the base version from the
	// branch/patch suffix so Version stays comparable to what the os-release
	// path reports on other systems.
	release := cstrToString(uts.Release[:])
	version := release
	if idx := indexByte(release, '-'); idx >= 0 {
		version = release[:idx]
	}

	return &OS{
		ID:       "freebsd",
		Name:     cstrToString(uts.Sysname[:]),
		Version:  version,
		Branding: release,
	}, nil
}

func cstrToString(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
