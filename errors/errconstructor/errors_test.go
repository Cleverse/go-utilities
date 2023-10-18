package errors

import (
	stderrors "errors"
	"fmt"
	"strings"
	"testing"
)

var exampleStd = stderrors.New("exampleStd error")

func TestNil(t *testing.T) {
	t.Run("WithStack", func(t *testing.T) {
		got := WithStack(nil, 0)
		if got != nil {
			t.Errorf("want nil, got %q", got)
		}
	})

	t.Run("Wrap", func(t *testing.T) {
		got := Wrap(nil, 0, "no error")
		if got != nil {
			t.Errorf("want nil, got %q", got)
		}
	})

	t.Run("Wrapf", func(t *testing.T) {
		got := Wrap(nil, 0, "no error")
		if got != nil {
			t.Errorf("want nil, got %q", got)
		}
	})
}

func TestWrap(t *testing.T) {
	testcases := []struct {
		err     error
		message string
		want    string
	}{
		{exampleStd, "add message", "add message: " + exampleStd.Error()},
		{WithStack(exampleStd, 0), "add message", "add message: " + exampleStd.Error()},
		{Wrap(exampleStd, 0, "wrapped err"), "add message", "add message: wrapped err: " + exampleStd.Error()},
		{Wrapf(exampleStd, 0, "wrapped err %s", "formatted"), "add message", "add message: wrapped err formatted: " + exampleStd.Error()},
		{Wrapf(exampleStd, 0, "wrapped err %q", "formatted"), "add message", "add message: wrapped err \"formatted\": " + exampleStd.Error()},
	}

	t.Run("Wrap", func(t *testing.T) {
		for i, testcase := range testcases {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				got := Wrap(testcase.err, 0, testcase.message).Error()
				if got != testcase.want {
					t.Errorf("want %q, got %q", testcase.want, got)
				}
			})
		}
	})

	t.Run("Wrapf", func(t *testing.T) {
		for i, testcase := range testcases {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				want := fmt.Sprintf("formatted %s", testcase.want)
				got := Wrapf(testcase.err, 0, "formatted %s", testcase.message).Error()
				if got != want {
					t.Errorf("want %q, got %q", want, got)
				}
			})
		}
	})
}

func TestFormat(t *testing.T) {
	tests := []struct {
		error
		format      string
		contains    []string
		notContains []string
	}{
		{
			WithStack(exampleStd, 0),
			"%s",
			[]string{exampleStd.Error()},
			[]string{},
		},
		{
			WithStack(exampleStd, 0),
			"%s",
			[]string{exampleStd.Error()},
			[]string{},
		},
		{
			WithStack(exampleStd, 0),
			"%q",
			[]string{`"` + exampleStd.Error() + `"`},
			[]string{},
		},
		{
			WithStack(exampleStd, 0),
			"%+v",
			[]string{
				"exampleStd error\n",
				"github.com/Cleverse/go-utilities/errors/errconstructor.TestFormat\n",
				"errconstructor/errors_test.go:98",
			},
			[]string{},
		},
		{
			Wrapf(exampleStd, 0, "wrapped err %s", "formatted"),
			"%+v",
			[]string{
				"wrapped err formatted\n",
				"exampleStd error\n",
				"github.com/Cleverse/go-utilities/errors/errconstructor.TestFormat\n",
				"errconstructor/errors_test.go:108",
			},
			[]string{},
		},
		{
			func() error {
				return WithStack(exampleStd, 0)
			}(),
			"%+v",
			[]string{
				"exampleStd error\n",
				"github.com/Cleverse/go-utilities/errors/errconstructor.TestFormat\n",
				"errconstructor/errors_test.go:120",
				"errconstructor/errors_test.go:121",
			},
			[]string{},
		},
		{
			func() error {
				return WithStack(customErrorWithStack1("call me maybe"), 0)
			}(),
			"%+v",
			[]string{
				"call me maybe\n",
				"github.com/Cleverse/go-utilities/errors/errconstructor.TestFormat\n",
				"errconstructor/errors_test.go:177",
				"errconstructor/errors_test.go:181",
				"errconstructor/errors_test.go:134",
				"errconstructor/errors_test.go:133",
			},
			[]string{
				"errconstructor/errors_test.go:173",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			testFormat(t, tt.contains, tt.notContains, tt.format, tt.error)
		})
	}
}

func testFormat(t *testing.T, wantContains, dontContains []string, format string, err error) {
	t.Helper()
	got := fmt.Sprintf(format, err)
	for _, contain := range wantContains {
		if !strings.Contains(got, contain) {
			t.Errorf("want %q to contain %q", got, contain)
		}
	}
	for _, contain := range dontContains {
		if strings.Contains(got, contain) {
			t.Errorf("%q shold not contain %q", got, contain)
		}
	}
}

func customError(format string, a ...any) error {
	return WithStack(fmt.Errorf(format, a...), 1)
}

func customErrorWithStack0(format string, a ...any) error {
	return WithStack(customError(format, a...), 0)
}

func customErrorWithStack1(format string, a ...any) error {
	return WithStack(customErrorWithStack0(format, a...), 0)
}
