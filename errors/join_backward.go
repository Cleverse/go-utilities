//go:build !go1.20
// +build !go1.20

package errors

import errors "github.com/Cleverse/go-utilities/errors/errconstructor"

// Disclaimer: functions Join is copied from the stdlib errors
// package v1.21.0.

// Join returns an error that wraps the given errors with a stack trace at the point WithStack was called.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		// stack: errors.Callers(0),
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return errors.WithStack(e, 1)
}

type joinError struct {
	// stack *errors.Stacks
	errs []error
}

func (e *joinError) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

// TODO: implement joinError.Format method for stack traces printing
