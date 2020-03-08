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
