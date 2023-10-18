package errors

import (
	"fmt"
	"io"
)

// WithStack annotates err with a skipable stack trace.
// If err is nil, WithStack returns nil.
//
// It's useful when you want to create a custom error type or custom error constructor function.
//
// skip parameter indicates how many stack frames to skip when recording the stack trace.
// skip = 0 means that the stack trace should begin at your current function call.
// skip = 1 means that the stack trace should begin at the caller of your function.
func WithStack(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &withStackError{
		cause: err,
		stack: Callers(skip),
	}
}

// Wrap returns an error annotating err with a  skipable stack trace
// and the supplied message.
//
// It's useful when you want to create a custom error type or custom error constructor function.
//
// skip parameter indicates how many stack frames to skip when recording the stack trace.
// skip = 0 means that the stack trace should begin at your current function call.
// skip = 1 means that the stack trace should begin at the caller of your function.
func Wrap(err error, skip int, message string) error {
	if err == nil {
		return nil
	}
	return &withStackError{
		cause: err,
		msg:   message,
		stack: Callers(skip),
	}
}

// Wrapf returns an error annotating err with a skipable stack trace
// and the format specifier.
//
// It's useful when you want to create a custom error type or custom error constructor function.
//
// skip parameter indicates how many stack frames to skip when recording the stack trace.
// skip = 0 means that the stack trace should begin at your current function call.
// skip = 1 means that the stack trace should begin at the caller of your function.
func Wrapf(err error, skip int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withStackError{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		stack: Callers(skip),
	}
}

// withStackError is an error that contains a cause and a stack trace.
type withStackError struct {
	stack *Stacks
	cause error
	msg   string
}

func (w *withStackError) Error() string {
	if w.msg == "" {
		return w.cause.Error()
	}
	return w.msg + ": " + w.cause.Error()
}

func (w *withStackError) Cause() error {
	return w.cause
}

func (w *withStackError) Unwrap() error {
	return w.cause
}

func (w *withStackError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// this will recursively call Format method of the cause for stack printing
			fmt.Fprintf(s, "%+v", w.Cause())
			if w.msg != "" {
				_, _ = io.WriteString(s, "\n")
				_, _ = io.WriteString(s, w.msg)
			}
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
