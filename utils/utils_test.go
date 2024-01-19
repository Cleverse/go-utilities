package utils

import (
	"bytes"
	"reflect"
	"testing"
)

// assertEqual asserts that expected and actual are equal.
func assertEqual[T comparable](t *testing.T, expected, actual T, msgAndArgs ...interface{}) bool {
	if expected != actual {
		if msg := msgFormatter(msgAndArgs...); msg != "" {
			t.Error(msg)
		} else {
			t.Errorf("expected: %v, got: %v", expected, actual)
		}
		return false
	}
	return true
}

func assertNotEqual[T comparable](t *testing.T, expected, actual T, msgAndArgs ...interface{}) bool {
	if expected == actual {
		if msg := msgFormatter(msgAndArgs...); msg != "" {
			t.Error(msg)
		} else {
			t.Errorf("expected: %v not equal to: %v", expected, actual)
		}
		return false
	}
	return true
}

// assertNotEqual asserts that expected and actual are not equal.
func assertTrue(t *testing.T, actual bool, msgAndArgs ...interface{}) bool {
	return assertEqual(t, true, actual, msgAndArgs...)
}

// assertNotEqual asserts that expected and actual are not equal.
func assertFalse(t *testing.T, actual bool, msgAndArgs ...interface{}) bool {
	return assertEqual(t, false, actual, msgAndArgs...)
}

func assertEqualAny(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) (equal bool) {
	// check nil
	switch {
	case expected == nil && actual == nil:
		return true
	case expected == nil && actual != nil:
		fallthrough
	case expected != nil && actual == nil:
		t.Errorf("expected: %v, got: %v", expected, actual)
		return false
	}

	// if types are different, they are not equal.
	t1 := reflect.TypeOf(expected)
	t2 := reflect.TypeOf(actual)
	if t1 != t2 {
		t.Errorf("type mismatch: expected: %v, got: %v", t1, t2)
		return false
	}

	defer func() {
		if !equal {
			if msg := msgFormatter(msgAndArgs...); msg != "" {
				t.Error(msg)
			} else {
				t.Errorf("expected: %v, got: %v", expected, actual)
			}
		}
	}()

	// default compare
	switch expected := expected.(type) {
	case []byte:
		actual := actual.([]byte)
		return bytes.Equal(expected, actual)
	case []string:
		actual := actual.([]string)
		if len(expected) != len(actual) {
			return false
		}
		for i := range expected {
			if expected[i] != actual[i] {
				return false
			}
		}
		return true
	default:
		return expected == actual
	}
}

func TestDefaultOptional(t *testing.T) {
	f := func(a int, b ...int) int {
		_b := DefaultOptional(b, 1)
		return a + _b
	}

	assertEqual(t, 3, f(2))
	assertEqual(t, 3, f(2, 0))
	assertEqual(t, 3, f(2, 1))
	assertEqual(t, 4, f(2, 2))
}
