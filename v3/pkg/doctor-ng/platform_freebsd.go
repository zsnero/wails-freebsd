//go:build freebsd

package doctorng

import (
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/doctor-ng/packagemanager"
)

func collectPlatformExtras() map[string]string {
	extras := make(map[string]string)

	// FreeBSD desktops are X11-only in practice: there is no Wayland session in
	// the base system, and GDK_BACKEND is pinned to x11 at startup. Report the
	// display that is actually in use rather than the XDG variable, which is
	// typically unset because FreeBSD sessions are started via start(8) rather
	// than a logind-style session manager.
	extras["DISPLAY"] = getEnvOrDefault("DISPLAY", "unset")
	extras["Desktop Environment"] = getEnvOrDefault("XDG_CURRENT_DESKTOP", "unset")
	extras["Compiler"] = getFreeBSDCompilerInfo()

	return extras
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getFreeBSDCompilerInfo reports the C compiler cgo will use. FreeBSD's base
// system provides clang as `cc`, and the gcc port (when installed) is not the
// default, so report whichever is actually first on PATH.
func getFreeBSDCompilerInfo() string {
	if cc := os.Getenv("CC"); cc != "" {
		return cc
	}
	for _, candidate := range []string{"cc", "clang", "gcc"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
	}
	return "not found"
}

func (d *Doctor) collectDependencies() error {
	info, _ := operatingsystem.Info()
	pm := packagemanager.Find(info.ID)
	if pm == nil {
		return nil
	}

	deps, err := packagemanager.Dependencies(pm)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		status := StatusMissing
		if dep.Installed {
			status = StatusOK
		}

		d.report.Dependencies = append(d.report.Dependencies, &Dependency{
			Name:           dep.Name,
			PackageName:    dep.PackageName,
			Version:        dep.Version,
			Status:         status,
			Required:       !dep.Optional,
			InstallCommand: dep.InstallCommand,
			Category:       categorizeFreeBSDDep(dep.Name),
		})
	}

	return nil
}

func categorizeFreeBSDDep(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "gtk"):
		return "gtk"
	case strings.Contains(lower, "webkit"):
		return "webkit"
	case name == "gcc" || name == "pkg-config":
		return "build-tools"
	case name == "npm":
		return "frontend"
	case name == "docker":
		return "optional"
	default:
		return "other"
	}
}

func (d *Doctor) runDiagnostics() {
	d.checkGoInstallation()
	d.checkFreeBSDSpecific()
}

func (d *Doctor) checkGoInstallation() {
	if d.report.Build.GoVersion == "" {
		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Go Installation",
			Message:  "Go installation not found",
			Severity: SeverityError,
			HelpURL:  "/getting-started/installation/",
			Fix: &Fix{
				Description: "Install Go from https://go.dev/dl/",
				ManualSteps: []string{
					"Download Go from https://go.dev/dl/",
					"Extract and add to PATH",
				},
			},
		})
	}
}

func (d *Doctor) checkFreeBSDSpecific() {
	missingRequired := d.report.Dependencies.RequiredMissing()
	if len(missingRequired) > 0 {
		var commands []string
		for _, dep := range missingRequired {
			if dep.InstallCommand != "" {
				commands = append(commands, dep.InstallCommand)
			}
		}

		d.report.Diagnostics = append(d.report.Diagnostics, DiagnosticResult{
			Name:     "Missing Dependencies",
			Message:  "Required system packages are not installed",
			Severity: SeverityError,
			HelpURL:  "/getting-started/installation/#freebsd",
			Fix: &Fix{
				Description:  "Install missing packages",
				Command:      strings.Join(commands, " && "),
				RequiresSudo: true,
			},
		})
	}
}
