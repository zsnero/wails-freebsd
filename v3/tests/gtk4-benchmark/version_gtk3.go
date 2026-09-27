//go:build (linux && gtk3) || freebsd

package main

func getGTKVersionString() string {
	return "GTK3 (WebKit2GTK 4.1)"
}
