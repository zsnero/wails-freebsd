//go:build freebsd && cgo && !server

package application

/*
#cgo pkg-config: glib-2.0 gtk+-3.0
#include <gtk/gtk.h>
#include <stdlib.h>

// wails_prefers_dark reports the value of the gtk-application-prefer-dark-theme
// GtkSetting. That property is what a desktop's "prefer dark" switch ends up
// setting, on FreeBSD as much as on Linux, and unlike the freedesktop
// color-scheme portal interface it needs no D-Bus session to read.
static int wails_prefers_dark(void) {
	GtkSettings *settings = gtk_settings_get_default();
	if (settings == NULL) {
		return 0;
	}
	gboolean prefer = FALSE;
	g_object_get(settings, "gtk-application-prefer-dark-theme", &prefer, NULL);
	return prefer ? 1 : 0;
}
*/
import "C"

import (
	"os"
	"strings"
)

// monitorThemeChanges seeds the theme once and then stops.
//
// On Linux this subscribes to the freedesktop Settings portal and keeps
// a.theme current for the life of the process. There is no equivalent
// subscription to make here: the portal is a freedesktop component and
// xdg-desktop-portal does not run on FreeBSD, so a live feed is not available
// without a second mechanism such as watching the GTK settings file.
//
// Reading the setting once is still much better than leaving a.theme empty,
// because isDarkMode() tests it for the substring "dark" and would otherwise
// report light mode for the entire session on a dark-themed desktop - the
// frontend would render light while the desktop is dark. A theme change made
// while the app is running needs a restart to be picked up.
func (a *linuxApp) monitorThemeChanges() {
	if wailsThemeFromSettings() {
		a.theme = "dark"
		return
	}
	a.theme = "light"
}

// wailsThemeFromSettings asks GTK first and falls back to GTK_THEME, which
// covers the case of an app started with an explicit theme on the command line
// or from a session that sets the variable without touching the GSettings
// database.
func wailsThemeFromSettings() bool {
	if C.wails_prefers_dark() == 1 {
		return true
	}
	// GTK_THEME is a colon-separated preference list, e.g. "Adwaita:dark".
	// Anything carrying "dark" is treated as a dark preference, matching how
	// isDarkMode classifies the value the portal reports.
	for _, candidate := range strings.Split(os.Getenv("GTK_THEME"), ":") {
		if strings.Contains(strings.ToLower(candidate), "dark") {
			return true
		}
	}
	return false
}

// monitorPowerEvents is a no-op.
//
// On Linux this follows org.freedesktop.login1's PrepareForSleep signal on the
// system bus. FreeBSD exposes no comparable notification to a desktop
// application, so SystemWillSleep and SystemDidWake are never emitted here.
// Apps that suspend on those events need to poll or be told by the user.
func (a *linuxApp) monitorPowerEvents() {}
