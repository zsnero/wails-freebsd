//go:build freebsd

package capabilities

// FreeBSD currently targets the GTK3 backend only, so this reports the GTK3
// / WebKit2GTK-4.1 combination unconditionally rather than probing. The Linux
// GTK3 variant declares the same values; a future FreeBSD GTK4 build will need
// a separate file selected by build tag, as Linux does with its `gtk3` tag.
func NewCapabilities() Capabilities {
	return Capabilities{
		HasNativeDrag: true,
		GTKVersion:    3,
		WebKitVersion: "4.1",
	}
}
