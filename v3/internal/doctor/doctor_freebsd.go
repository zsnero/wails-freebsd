//go:build freebsd

package doctor

import (
	"os"
	"os/exec"

	"github.com/wailsapp/wails/v3/internal/doctor/packagemanager"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)

	// FreeBSD has no XDG session type in practice (sessions come from start(8),
	// not a logind-style session manager) and the backend is always X11, so
	// report the display that is actually in use.
	result["DISPLAY"] = getEnvOrDefault("DISPLAY", "unset")

	// Check desktop environment
	result["Desktop Environment"] = getEnvOrDefault("XDG_CURRENT_DESKTOP", "unset")

	// The C compiler cgo will use. Base FreeBSD provides clang as `cc`.
	result["C Compiler"] = getFreeBSDCompilerInfo()

	return result, true
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

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

func checkPlatformDependencies(result map[string]string, ok *bool) {
	info, _ := operatingsystem.Info()

	pm := packagemanager.Find(info.ID)
	deps, _ := packagemanager.Dependencies(pm)
	for _, dep := range deps {
		var status string

		switch true {
		case !dep.Installed:
			if dep.Optional {
				status = "[Optional] "
			} else {
				*ok = false
			}
			status += "not installed."
			if dep.InstallCommand != "" {
				status += " Install with: " + dep.InstallCommand
			}
		case dep.Version != "":
			status = dep.Version
		}

		result[dep.Name] = status
	}

	checkCommonDependencies(result, ok)
}
