package utils // import github.com/mkaz/hastie/pkg/utils

import (
	"os"
	"strings"
)

// RemoveIndexHTML removes "index.html" from a URL string
func RemoveIndexHTML(str string) string {
	return strings.Replace(str, "index.html", "", 1)
}

// TrimSlash removes the trailing slash
func TrimSlash(path string) string {
	return strings.Trim(path, "/")
}

// FileExists checks if a file or directory exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
