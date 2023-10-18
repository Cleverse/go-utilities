package errors

import (
	"errors"
	stderrors "errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

var (
	exampleStd = stderrors.New("example error")
	example    = New("example error")
)

func TestNew(t *testing.T) {
	testcases := []struct {
		err  string
		want error
	}{
		{"new error", New("new error")},
		{"with format: %v", errors.New("with format: %v")},
	}

	for _, testcase := range testcases {
		t.Run(testcase.err, func(t *testing.T) {
			got := New(testcase.err)
			if got.Error() != testcase.want.Error() {
				t.Errorf("want %q, got %q", testcase.want, got)
			}
			if Is(got, testcase.want) {
				t.Error("should not be equal errors")
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	testcases := []struct {
		format string
		args   []interface{}
		want   error
	}{{
		"new error",
		nil,
		New("new error"),
	}, {
		"with format: %v",
		[]interface{}{"value"},
		errors.New("with format: value"),
	}}

	for _, testcase := range testcases {
		t.Run(testcase.format, func(t *testing.T) {
			got := Errorf(testcase.format, testcase.args...)
			if got.Error() != testcase.want.Error() {
				t.Errorf("want %q, got %q", testcase.want, got)
			}
			if Is(got, testcase.want) {
				t.Error("should not be equal errors")
			}
		})
	}
}

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format   string
		contains []string
	}{{
		New("error"),
		"%s",
		[]string{"error"},
	}, {
		New("error"),
		"%v",
		[]string{"error"},
	}, {
		New("error"),
		"%+v",
		[]string{
			"error\n",
			"github.com/Cleverse/go-utilities/errors.TestFormatNew\n",
			"errors/errors_test.go:81",
		},
	}, {
		New("error"),
		"%q",
		[]string{`"error"`},
	}}

	for _, tt := range tests {
		testFormat(t, tt.contains, tt.format, tt.error)
	}
}

func TestCause(t *testing.T) {
	testcases := []struct {
		err  error
		want error
	}{
		{
			// nil error is nil
			err:  nil,
			want: nil,
		},
		{
			// explicit nil error is nil
			err:  (error)(nil),
			want: nil,
		},
		{
			// uncaused error is unaffected
			err:  exampleStd,
			want: exampleStd,
		},
		{
			// caused error returns cause
			err:  Wrap(exampleStd, "ignored"),
			want: exampleStd,
		},
		{
			err:  example,
			want: example,
		},
		{
			WithStack(nil),
			nil,
		},
		{
			WithStack(io.EOF),
			io.EOF,
		},
		{
			WithStack(example),
			example,
		},
	}

	for _, testcase := range testcases {
		t.Run("", func(t *testing.T) {
			got := Cause(testcase.err)
			if got != testcase.want {
				t.Errorf("want %#v, got %#v", testcase.want, got)
			}
		})
	}
}

func testFormat(t *testing.T, wantContains []string, format string, err error) {
	t.Helper()
	got := fmt.Sprintf(format, err)
	for _, contain := range wantContains {
		if !strings.Contains(got, contain) {
			t.Errorf("want %q to contain %q", got, contain)
		}
	}
}
