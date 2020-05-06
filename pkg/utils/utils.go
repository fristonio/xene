package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// RandToken generates a random toke string of the provided size
func RandToken(size uint32) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
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
