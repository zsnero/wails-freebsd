//go:build freebsd && cgo && !android && !server

package application

// portalGlobalShortcutsSupported reports whether
// org.freedesktop.portal.GlobalShortcuts can be reached. The interface is
// implemented by xdg-desktop-portal, which is a freedesktop component and is
// not part of a FreeBSD desktop, so a Wayland session there must fall back to
// X11 instead of selecting a backend that will never answer.
const portalGlobalShortcutsSupported = false
