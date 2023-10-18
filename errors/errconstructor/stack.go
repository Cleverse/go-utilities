package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

// The maximum number of stackframes on any error.
var MaxStackDepth = 32

// Stacks represents a Stacks of program counters.
type Stacks []uintptr

func (s *Stacks) Format(st fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := StackFrame(pc)
				fmt.Fprintf(st, "\n%+s", f)
			}
		default:
			for _, pc := range *s {
				f := StackFrame(pc)
				fmt.Fprintf(st, "\n%s", f)
			}
		}
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := StackFrame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		default:
			for _, pc := range *s {
				f := StackFrame(pc)
				fmt.Fprintf(st, "\n%s", f)
			}
		}
	}
}

// Callers returns stack of program counters.
func Callers(skip int) *Stacks {
	pcs := make([]uintptr, MaxStackDepth)

	// Skip 3 stack frames cause:
	// 0. at runtime.Callers itself
	// 1. at runtime.Callers call
	// 2. at errors.callers call
	// 3. at errors.New, errors.Errorf, errors.WithStack and errors.Wrap functions call
	const callersSkip = 3

	n := runtime.Callers(callersSkip+skip, pcs[:])
	var st Stacks = pcs[0:n]
	return &st
}

// StackFrame represents a program counter inside a stack frame.
// For historical reasons if StackFrame is interpreted as a uintptr
// its value represents the program counter + 1.
type StackFrame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f StackFrame) pc() uintptr {
	return uintptr(f) - 1
}

// Info returns the name, file and line number for this Frame's pc.
//
// name is the name of the function, or "unknown".
// file is the full path to the file that contains the function, or "unknown".
// line is the line number of the source code of the function.
func (f StackFrame) Info() (name string, file string, line int) {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown", "unknown", 0
	}
	file, line = fn.FileLine(f.pc())
	return fn.Name(), file, line
}

// Format formats the Stackframe according to the fmt.Formatter interface.
//
//	%s	function name and source file
//	%v	function name, source file and line number
//
// Format accepts flag "+" for %s and %v verbs for printing source file with
// relative path to the compile time.
func (f StackFrame) Format(s fmt.State, verb rune) {
	name, file, line := f.Info()
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, name)
			_, _ = io.WriteString(s, "\n\t")
			_, _ = io.WriteString(s, file)
		default:
			_, _ = io.WriteString(s, path.Base(file))
		}
	case 'v':
		f.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		_, _ = io.WriteString(s, strconv.Itoa(line))
	default:
		f.Format(s, 's')
	}
}

// MarshalText formats a stacktrace Frame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (f StackFrame) MarshalText() ([]byte, error) {
	name, file, line := f.Info()
	if name == "unknown" {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, file, line)), nil
}
