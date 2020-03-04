package oauth

import (
	"crypto/rand"
	"encoding/base64"
)

// randToken generates a random toke string of the provided size
func randToken(size uint32) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
