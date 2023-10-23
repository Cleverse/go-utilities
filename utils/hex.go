package utils

import (
	"encoding/hex"
	"fmt"
)

// RandomHex returns a random hex string with the given length.
func RandomHex(length int) string {
	// TODO: reduce memory allocation by using a same buffer for random and hex encoding.
	return hex.EncodeToString(RandomBytes(length))
}

// Has0xPrefix checks if the input string has 0x prefix or not.
//
// Returns `true“ if the input string has 0x prefix, otherwise `false`.
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// Trim0xPrefix returns the input string without 0x prefix.
func Trim0xPrefix(input string) string {
	if Has0xPrefix(input) {
		return input[2:]
	}
	return input
}

// Add0xPrefix returns the input string with 0x prefix.
func Add0xPrefix(input string) string {
	if !Has0xPrefix(input) {
		return "0x" + input
	}
	return input
}

// Flip0xPrefix returns the input string with 0x prefix if it doesn't have 0x prefix, otherwise returns the input string without 0x prefix.
func Flip0xPrefix(input string) string {
	if Has0xPrefix(input) {
		return input[2:]
	}
	return "0x" + input
}

// IsHex verifies whether a string can represent a valid hex-encoded or not.
func IsHex(str string) bool {
	str = Trim0xPrefix(str)
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// DecodeHex decodes a hex string into a byte slice. str can be prefixed with 0x.
func DecodeHex(str string) ([]byte, error) {
	b, err := hex.DecodeString(Trim0xPrefix(str))
	return b, fmt.Errorf("decode hex: %w", err)
}
