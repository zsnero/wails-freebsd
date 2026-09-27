//go:build linux && cgo && !android && !server

package application

// portalGlobalShortcutsSupported reports whether
// org.freedesktop.portal.GlobalShortcuts can be reached. On Linux it can,
// which is why a Wayland session prefers the portal over X11.
const portalGlobalShortcutsSupported = true
