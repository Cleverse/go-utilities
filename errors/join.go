//go:build !go1.20
// +build !go1.20

package errors

// Join returns an error that wraps the given errors with a stack trace at the point WithStack was called.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	// TODO: implement Join for backwards compatibility
	panic("not implemented")
}
