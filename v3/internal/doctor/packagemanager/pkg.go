//go:build linux || freebsd

package packagemanager

import (
	"os/exec"
	"strings"
)

// Pkg represents the FreeBSD pkg package manager
type Pkg struct {
	name string
	osid string
}

// NewPkg creates a new Pkg instance
func NewPkg(osid string) *Pkg {
	return &Pkg{
		name: "pkg",
		osid: osid,
	}
}

// Packages returns the packages needed to build Wails on FreeBSD.
//
// The names are the FreeBSD *package* names, which frequently differ from the
// pkg-config module names. Notably www/webkit2-gtk is a single port with three
// flavours and the flavour suffix is part of the package name, so the GTK3 +
// libsoup3 build we require is `webkit2-gtk_41` — that is the flavour that
// installs webkit2gtk-4.1.pc, which is what the cgo directives link against.
//
// GTK3 is listed as a hard requirement rather than optional because the FreeBSD
// port currently targets the GTK3 backend only. The GTK4 pair is reported as
// optional so the eventual GTK4 work has somewhere to hang its dependency
// reporting.
func (p *Pkg) Packages() Packagemap {
	return Packagemap{
		"gtk3": []*Package{
			{Name: "gtk3", SystemPackage: true, Library: true},
		},
		"webkit2gtk-4.1": []*Package{
			{Name: "webkit2-gtk_41", SystemPackage: true, Library: true},
		},
		"gtk4": []*Package{
			{Name: "gtk4", SystemPackage: true, Library: true, Optional: true},
		},
		"webkitgtk-6.0": []*Package{
			{Name: "webkit2-gtk_60", SystemPackage: true, Library: true, Optional: true},
		},
		"pkg-config": []*Package{
			{Name: "pkgconf", SystemPackage: true},
		},
		// cgo needs a working C compiler. FreeBSD's base system ships clang as
		// `cc`, so this is deliberately not a SystemPackage check: requiring the
		// gcc port here would report a missing dependency on a machine that can
		// already build. InstallCheck probes for a compiler on PATH instead, and
		// the install command is only surfaced if none is found.
		"gcc": []*Package{
			{
				Name:           "gcc",
				InstallCommand: "sudo pkg install -y gcc",
				InstallCheck: func() bool {
					return commandExists("cc") || commandExists("clang") || commandExists("gcc")
				},
			},
		},
		"npm": []*Package{
			{Name: "www/npm", SystemPackage: true},
		},
	}
}

// Name returns the name of the package manager
func (p *Pkg) Name() string {
	return p.name
}

// queryRepo looks a package up in the configured repositories rather than in
// the installed set, so "available" and "installed" stay distinct questions.
//
// `pkg rquery` exits 0 whether or not anything matches, so availability is
// determined from the output being non-empty. The format string yields
// "<name>-<version>", matching the shape the other backends parse.
func (p *Pkg) queryRepo(name string) (string, error) {
	return execCmd("pkg", "rquery", "-e", `%n == "`+name+`"`, "%n-%v")
}

// PackageAvailable tests if the given package exists in the repositories.
// It records the repository's version on pkg as a side effect, matching the
// contract the other backends rely on.
func (p *Pkg) PackageAvailable(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		return true, nil
	}
	output, err := p.queryRepo(pkg.Name)
	available := strings.TrimSpace(output) != ""
	if available {
		p.setPackageVersion(pkg, output)
	}
	return available, err
}

// PackageInstalled tests if the given package is installed, and records the
// installed version on pkg.
func (p *Pkg) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		if pkg.InstallCheck != nil {
			return pkg.InstallCheck(), nil
		}
		return false, nil
	}
	// `pkg info -e` exits 0 if the package is installed and 1 otherwise. It
	// cannot distinguish "not installed" from "no such package", but
	// PackageAvailable has already established that the package exists in the
	// repositories by the time this runs.
	if err := exec.Command("pkg", "info", "-e", pkg.Name).Run(); err != nil {
		return false, nil
	}
	// Prefer the installed version over the repository's; a machine can lag the
	// repo by many revisions, and reporting the repo version would misreport
	// what the build will actually link against.
	if output, err := execCmd("pkg", "query", "%n-%v", pkg.Name); err == nil {
		p.setPackageVersion(pkg, output)
	}
	return true, nil
}

// InstallCommand returns the pkg command to install a package
func (p *Pkg) InstallCommand(pkg *Package) string {
	if !pkg.SystemPackage {
		return pkg.InstallCommand
	}
	return "sudo pkg install -y " + pkg.Name
}

func (p *Pkg) setPackageVersion(pkg *Package, output string) {
	// Output is "<name>-<version>"; strip the name prefix to get the version.
	name := pkg.Name + "-"
	if trimmed := strings.TrimSpace(output); strings.HasPrefix(trimmed, name) {
		pkg.Version = strings.TrimPrefix(trimmed, name)
	}
}
