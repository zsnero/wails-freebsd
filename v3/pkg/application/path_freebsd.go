//go:build freebsd

package application

import (
	"os"
	"path/filepath"
	"strconv"
)

// runtimeDir returns $XDG_RUNTIME_DIR when the session set one, and otherwise a
// per-user directory under the system temporary directory.
//
// The fallback exists because adrg/xdg defaults to /run/user/<uid>, which is a
// systemd layout: FreeBSD has no /run/user and no component that creates one, so
// on a FreeBSD box that is not started from a desktop session launcher - ssh, a
// service, a cron job - xdg.RuntimeDir names a directory that does not exist.
// Anything using Path(PathRuntimeDir) for a socket or a lock would fail there
// with ENOENT rather than at anything to do with its own logic.
//
// A desktop session on FreeBSD does export XDG_RUNTIME_DIR, so the first branch
// is the normal case and this only covers the sessions that do not.
//
// The directory is not created here. Path is a lookup and callers may use it
// for display; a caller that needs it to exist should MkdirAll it, and must
// create it with mode 0700 because the system temp directory is world-readable
// and the runtime directory is expected to be private to the user.
func runtimeDir() string {
	if dir := os.Getenv("XDG_RUNTIME_DIR"); dir != "" {
		return dir
	}
	return filepath.Join(os.TempDir(), "wails-"+strconv.Itoa(os.Getuid()))
}
