//go:build freebsd && !server

package application

import "errors"

// errSystemTrayUnsupported is returned by every operation that needs a real
// tray icon. The Linux implementation is a StatusNotifierItem, and
// StatusNotifier is specified in terms of the freedesktop D-Bus interfaces; it
// is not a freedesktop standard, so nothing on FreeBSD implements it just
// because a D-Bus daemon is present.
//
// The stub keeps the API shape so an app that guards its tray code with
// runtime.GOOS still compiles and runs, and so bounds/positionWindow report
// failure instead of quietly returning a zero rect that callers would use to
// place a window at 0,0.
var errSystemTrayUnsupported = errors.New("system tray is not supported on FreeBSD: StatusNotifierItem has no implementation here")

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &freebsdSystemTray{parent: s}
}

// freebsdSystemTray is a no-op system tray. See errSystemTrayUnsupported.
type freebsdSystemTray struct {
	parent *SystemTray
}

func (t *freebsdSystemTray) setLabel(label string)            {}
func (t *freebsdSystemTray) setTooltip(tooltip string)        {}
func (t *freebsdSystemTray) run()                             {}
func (t *freebsdSystemTray) setIcon(icon []byte)              {}
func (t *freebsdSystemTray) setMenu(menu *Menu)               {}
func (t *freebsdSystemTray) setIconPosition(pos IconPosition) {}
func (t *freebsdSystemTray) setTemplateIcon(icon []byte)      {}
func (t *freebsdSystemTray) destroy()                         {}
func (t *freebsdSystemTray) setDarkModeIcon(icon []byte)      {}
func (t *freebsdSystemTray) openMenu()                        {}
func (t *freebsdSystemTray) Show()                            {}
func (t *freebsdSystemTray) Hide()                            {}

func (t *freebsdSystemTray) bounds() (*Rect, error) {
	return nil, errSystemTrayUnsupported
}

func (t *freebsdSystemTray) getScreen() (*Screen, error) {
	return nil, errSystemTrayUnsupported
}

func (t *freebsdSystemTray) positionWindow(w Window, o int) error {
	return errSystemTrayUnsupported
}
