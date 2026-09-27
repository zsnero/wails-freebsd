//go:build linux || freebsd

package dock

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type x11Dock struct{}

// New creates a new Dock Service.
// On the X11 desktops (Linux, FreeBSD) this returns a stub implementation since dock icon visibility
// and badge functionality are not standardized across desktop environments.
func New() *DockService {
	return &DockService{
		impl: &x11Dock{},
	}
}

// NewWithOptions creates a new dock service with badge options.
// On the X11 desktops (Linux, FreeBSD) this returns a stub implementation since badge functionality
// is not standardized across desktop environments. Options are ignored.
func NewWithOptions(options BadgeOptions) *DockService {
	return New()
}

func (l *x11Dock) Startup(ctx context.Context, options application.ServiceOptions) error {
	// No-op: the X11 desktops have no standardized dock/badge support
	return nil
}

func (l *x11Dock) Shutdown() error {
	// No-op: the X11 desktops have no standardized dock/badge support
	return nil
}

// HideAppIcon is a stub on the X11 desktops since dock icon visibility is not
// standardized across desktop environments.
func (l *x11Dock) HideAppIcon() {
	// No-op: the X11 desktops have no standardized dock icon visibility support
}

// ShowAppIcon is a stub on the X11 desktops since dock icon visibility is not
// standardized across desktop environments.
func (l *x11Dock) ShowAppIcon() {
	// No-op: the X11 desktops have no standardized dock icon visibility support
}

// SetBadge is a stub on the X11 desktops since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *x11Dock) SetBadge(label string) error {
	// No-op: the X11 desktops have no standardized badge support
	return nil
}

// SetCustomBadge is a stub on the X11 desktops since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *x11Dock) SetCustomBadge(label string, options BadgeOptions) error {
	// No-op: the X11 desktops have no standardized badge support
	return nil
}

// RemoveBadge is a stub on the X11 desktops since most desktop environments don't support
// application dock badges. This method exists for cross-platform compatibility.
func (l *x11Dock) RemoveBadge() error {
	// No-op: the X11 desktops have no standardized badge support
	return nil
}

func (l *x11Dock) GetBadge() *string {
	// No-op: the X11 desktops have no standardized badge support
	return nil
}
