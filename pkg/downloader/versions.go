package downloader

import (
	"strings"
)

// GetBaseVersion removes any version prefixes (ex. "v") from version numbers
func GetBaseVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version[1:]
	}

	return version
}
