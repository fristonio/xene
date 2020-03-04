package utils

import "os"

// FileExists checks if the file in the arguments exists or not
func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	return false
}
