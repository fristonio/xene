package utils

import "os"

// FileExists checks if the file in the arguments exists or not
func FileExists(file string) bool {
	if d, err := os.Stat(file); err == nil {
		if !d.IsDir() {
			return true
		}
	}

	return false
}

// DirExists checks if the provided directory exists or not.
func DirExists(dir string) bool {
	if d, err := os.Stat(dir); err == nil {
		if d.IsDir() {
			return true
		}
	}

	return false
}
