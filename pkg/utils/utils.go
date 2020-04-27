package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandToken generates a random toke string of the provided size
func RandToken(size uint32) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// CheckStringSliceEqual checks if the two slices provided
// are equal or not.
func CheckStringSliceEqual(a, b []string) bool {
	m := make(map[string]struct{})
	for _, x := range a {
		m[x] = struct{}{}
	}

	for _, x := range b {
		if _, ok := m[x]; !ok {
			return false
		}
	}

	return true
}
