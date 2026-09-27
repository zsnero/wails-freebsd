//go:build linux || freebsd

package operatingsystem

import "strings"

// parseOsRelease parses the freedesktop os-release format shared by Linux and
// FreeBSD. Fields that are absent from the file are left at the "Unknown"
// defaults rather than being treated as an error, because the spec marks every
// key except NAME optional.
func parseOsRelease(osRelease string) *OS {

	// Default value
	var result OS
	result.ID = "Unknown"
	result.Name = "Unknown"
	result.Version = "Unknown"

	// Split into lines
	lines := strings.Split(osRelease, "\n")
	// Iterate lines
	for _, line := range lines {
		// Split each line by the equals char
		splitLine := strings.SplitN(line, "=", 2)
		// Check we have
		if len(splitLine) != 2 {
			continue
		}
		switch splitLine[0] {
		case "ID":
			result.ID = strings.ToLower(strings.Trim(splitLine[1], `"`))
		case "NAME":
			result.Name = strings.Trim(splitLine[1], `"`)
		case "VERSION_ID":
			result.Version = strings.Trim(splitLine[1], `"`)
		case "VERSION":
			result.Branding = strings.Trim(splitLine[1], `"`)
		}
	}
	return &result
}
