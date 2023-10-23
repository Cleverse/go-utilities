//go:build go1.20
// +build go1.20

package utils

import (
	"unsafe"
)

// UnsafeString returns a string pointer without allocation
func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// UnsafeBytes returns a byte pointer without allocation
func UnsafeBytes(s string) (bs []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// CopyString copies a string to make it immutable
func CopyString(s string) string {
	return string(UnsafeBytes(s))
}
