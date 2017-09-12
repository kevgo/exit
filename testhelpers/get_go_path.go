package testhelpers

import (
	"os"
	"strings"
)

// GetGoPath returns the GOPATH environment variable content to use
func GetGoPath() string {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	return gopaths[len(gopaths)-1]
}
