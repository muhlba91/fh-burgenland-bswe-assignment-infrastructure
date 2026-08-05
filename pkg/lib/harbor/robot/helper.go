package robot

import (
	"os"
	"strings"
)

// harborURL returns the Harbor URL without the "https://" prefix.
func harborURL() string {
	return strings.ReplaceAll(os.Getenv("HARBOR_URL"), "https://", "")
}
