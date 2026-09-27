//go:build !freebsd

package application

import "github.com/adrg/xdg"

// runtimeDir returns the per-user runtime directory.
//
// It is a call rather than a direct read of xdg.RuntimeDir because that
// library's default is a systemd layout, /run/user/<uid>, which not every
// platform Wails builds for has. Platforms where the default is wrong override
// this in their own file; FreeBSD is the one that does, see path_freebsd.go.
func runtimeDir() string {
	return xdg.RuntimeDir
}
