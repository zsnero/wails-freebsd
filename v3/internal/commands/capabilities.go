package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/pterm/pterm"
)

// GTKCapabilities reports which GTK and WebKitGTK toolkits are installed,
// detected via pkg-config. It is shared by every host that uses the GTK
// toolkits (Linux and FreeBSD); only the JSON key it is reported under differs.
type GTKCapabilities struct {
	GTK4Available        bool   `json:"gtk4_available"`
	GTK3Available        bool   `json:"gtk3_available"`
	WebKitGTK6Available  bool   `json:"webkitgtk_6_available"`
	WebKit2GTK4Available bool   `json:"webkit2gtk_4_1_available"`
	Recommended          string `json:"recommended"`
}

type Capabilities struct {
	Platform string           `json:"platform"`
	Arch     string           `json:"arch"`
	Linux    *GTKCapabilities `json:"linux,omitempty"`
	FreeBSD  *GTKCapabilities `json:"freebsd,omitempty"`
}

type ToolCapabilitiesOptions struct{}

func ToolCapabilities(_ *ToolCapabilitiesOptions) error {
	caps := Capabilities{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	switch runtime.GOOS {
	case "linux":
		caps.Linux = detectGTKCapabilities()
	case "freebsd":
		caps.FreeBSD = detectGTKCapabilities()
	}

	pterm.Println(capsToJSON(caps))
	return nil
}

func detectGTKCapabilities() *GTKCapabilities {
	caps := &GTKCapabilities{}

	caps.GTK4Available = pkgConfigExists("gtk4")
	caps.WebKitGTK6Available = pkgConfigExists("webkitgtk-6.0")
	caps.GTK3Available = pkgConfigExists("gtk+-3.0")
	caps.WebKit2GTK4Available = pkgConfigExists("webkit2gtk-4.1")

	if caps.GTK4Available && caps.WebKitGTK6Available {
		caps.Recommended = "gtk4"
	} else if caps.GTK3Available && caps.WebKit2GTK4Available {
		caps.Recommended = "gtk3"
	} else {
		caps.Recommended = "none"
	}

	return caps
}

func pkgConfigExists(pkg string) bool {
	cmd := exec.Command("pkg-config", "--exists", pkg)
	return cmd.Run() == nil
}

func capsToJSON(caps Capabilities) string {
	result := fmt.Sprintf(`{"platform":"%s","arch":"%s"`, caps.Platform, caps.Arch)
	// Ordered, not a map range: only one entry is ever populated, but map
	// iteration order would make the output nondeterministic if that ever
	// changed, and this output is diffed in tests and by users.
	for _, entry := range []struct {
		key string
		gtk *GTKCapabilities
	}{
		{"linux", caps.Linux},
		{"freebsd", caps.FreeBSD},
	} {
		if entry.gtk == nil {
			continue
		}
		gtk := entry.gtk
		result += fmt.Sprintf(`,"%s":{"gtk4_available":%t,"gtk3_available":%t,"webkitgtk_6_available":%t,"webkit2gtk_4_1_available":%t,"recommended":"%s"}`,
			entry.key, gtk.GTK4Available, gtk.GTK3Available, gtk.WebKitGTK6Available, gtk.WebKit2GTK4Available, gtk.Recommended)
	}
	result += "}"
	return result
}
