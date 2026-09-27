//go:build freebsd && cgo && !server

package application

/*
#cgo pkg-config: glib-2.0 gtk+-3.0
#include <gtk/gtk.h>
#include <stdlib.h>

// wails_theme_settings reports how the desktop expresses its light/dark
// preference, and returns 0 when the answer is not available yet.
//
// GtkSettings does not exist until GTK has been initialised, and GTK is not
// initialised until g_application_run: gtk_application_new() only constructs
// the object. So gtk_settings_get_default() returns NULL for the whole window
// between process start and the GTK main loop coming up, and a caller that
// asks too early silently gets "not dark". Callers must ask after startup.
static int wails_theme_settings(int *prefer_dark, char **theme_name) {
	GtkSettings *settings = gtk_settings_get_default();
	if (settings == NULL) {
		return 0;
	}
	gboolean prefer = FALSE;
	g_object_get(settings, "gtk-application-prefer-dark-theme", &prefer, "gtk-theme-name", theme_name, NULL);
	*prefer_dark = prefer ? 1 : 0;
	return 1;
}
*/
import "C"

import (
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/events"
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
	seed := func() {
		if wailsThemeFromSettings() {
			a.theme = "dark"
			return
		}
		a.theme = "light"
	}

	// Best effort for the case where something has already forced GTK to
	// initialise, for example a window built before Run was called.
	seed()

	// The authoritative read has to wait for ApplicationStartup, which is
	// emitted once the GTK main loop is running. Reading earlier is the trap
	// this indirection exists to avoid: it used to happen inline here, and
	// gtk_settings_get_default() was still NULL, so a dark-themed desktop was
	// reported as light for the whole session.
	a.parent.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(_ *ApplicationEvent) {
		seed()
	})
}

// wailsThemeFromSettings asks GTK how the desktop is themed, and falls back to
// GTK_THEME for an app started with an explicit theme on the command line or
// from a session that sets the variable without touching the settings file.
//
// Both signals GTK knows about count. The prefer-dark boolean is what a
// desktop's "dark mode" switch sets, but a theme name is just as common on
// the BSDs - a user who picks Gruvbox-Dark or Adwaita-dark in a settings file
// often leaves the boolean alone, and reading only the boolean would call that
// a light desktop. Substring-matching the name for "dark" is the same rule
// isDarkMode() applies to a.theme.
func wailsThemeFromSettings() bool {
	var preferDark C.int
	var themeName *C.char
	if C.wails_theme_settings(&preferDark, &themeName) == 0 {
		// GTK is not up yet, so the answer is unknown rather than light.
		return false
	}
	if preferDark == 1 {
		return true
	}
	if themeName != nil && strings.Contains(strings.ToLower(C.GoString(themeName)), "dark") {
		return true
	}
	// GTK_THEME is a colon-separated preference list, e.g. "Adwaita:dark".
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
