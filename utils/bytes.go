package utils

import (
	"crypto/rand"
)

// RandomBytes returns a random byte slice with the given length with crypto/rand.
func RandomBytes(length int) []byte {
	b := make([]byte, length)
	Must(rand.Read(b))
	return b
}
