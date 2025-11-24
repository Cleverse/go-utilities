// Package errors provides Pure Golang errors library with stacktrace support (for wrapping, formatting and joining an errors).
//
// Add errors context to your report just by wrapping it with stacktrace.
//
//	b, err := json.Unmarshal(data, &Data)
//	if err != nil {
//	        return errors.WithStack(err)
//	}
//
// Describe more error details by wrapping it with message.
//
//	data := "foo"
//	if err := Do("foo"); err != nil {
//	        return errors.Wrapf(err, "failed to do %q", data)
//	}
//
// Debugging will be easier when you can see the stacktrace.
//
//	if err := HugeStackCall(ctx); err != nil {
//	        fmt.Printf("[ERROR]: %+v", err)
//	}
//
// Deprecated: Use the github.com/cockroachdb/errors instead.
package errors

import (
	"fmt"
	"io"

	errors "github.com/Cleverse/go-utilities/errors/errconstructor"
)

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
//
// this function is supports to skip stack frames. default is 0 (begin at your current function call)
// for more information see WithStack function in github.com/Cleverse/go-utilities/errors/errconstructor
func New(text string, skip ...int) error {
	skipNum := 0
	if len(skip) > 0 {
		skipNum = skip[0]
	}
	return &stdError{
		s:     text,
		stack: errors.Callers(skipNum),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &stdError{
		s:     fmt.Sprintf(format, args...),
		stack: errors.Callers(0),
	}
}

// stdError is an standard error that contains a message and a stack trace.
type stdError struct {
	stack *errors.Stacks
	s     string
}

func (f *stdError) Error() string {
	return f.s
}

func (f *stdError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, f.s)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, f.s)
	case 'q':
		fmt.Fprintf(s, "%q", f.s)
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	return errors.WithStack(err, 1)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return errors.Wrap(err, 1, message)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, 1, format, args...)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
//
// nolint: errorlint, wrapcheck
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(interface{ Cause() error })
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
