//go:build (linux && gtk3 && !android && !server) || (freebsd && !android && !server)

package application

func gtkDispatch(fn func()) {
	InvokeAsync(fn)
}
