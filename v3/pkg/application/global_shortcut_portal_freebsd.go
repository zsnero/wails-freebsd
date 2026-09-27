//go:build freebsd && cgo && !android && !server

package application

// portalGlobalShortcutsSupported reports whether
// org.freedesktop.portal.GlobalShortcuts can be reached. The interface is
// implemented by xdg-desktop-portal, which is a freedesktop component and is
// not part of a FreeBSD desktop, so a Wayland session there must fall back to
// X11 instead of selecting a backend that will never answer.
const portalGlobalShortcutsSupported = false

// newGlobalShortcutsPortalImpl is the per-OS shim that keeps
// newPortalGlobalShortcuts - which is defined in a Linux-only file - from
// being referenced by the shared dispatcher.
//
// The dispatcher only calls this when portalGlobalShortcutsSupported is true,
// so this is unreachable on FreeBSD. It returns the X11 backend rather than
// panicking so that a future change to the guard degrades to a working
// implementation instead of crashing an app at shortcut-registration time.
func newGlobalShortcutsPortalImpl(manager *GlobalShortcutManager) globalShortcutImpl {
	return newX11GlobalShortcuts(manager)
}
