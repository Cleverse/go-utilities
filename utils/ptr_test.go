package utils

import (
	"testing"
)

func TestRef(t *testing.T) {
	type T int

	val := T(0)
	pointer := PtrOf(val)
	if *pointer != val {
		t.Errorf("expected %d, got %d", val, *pointer)
	}

	val = T(1)
	pointer = PtrOf(val)
	if *pointer != val {
		t.Errorf("expected %d, got %d", val, *pointer)
	}
}

func TestDerefOr(t *testing.T) {
	type T int

	var val, def T = 1, 0

	out := DerefPtrOr(&val, def)
	if out != val {
		t.Errorf("expected %d, got %d", val, out)
	}

	out = DerefPtrOr(nil, def)
	if out != def {
		t.Errorf("expected %d, got %d", def, out)
	}
}

func TestEqual(t *testing.T) {
	type T int

	if !EqualPtr[T](nil, nil) {
		t.Errorf("expected true (nil == nil)")
	}
	if !EqualPtr(ToPtr(T(123)), ToPtr(T(123))) {
		t.Errorf("expected true (val == val)")
	}
	if EqualPtr(nil, ToPtr(T(123))) {
		t.Errorf("expected false (nil != val)")
	}
	if EqualPtr(ToPtr(T(123)), nil) {
		t.Errorf("expected false (val != nil)")
	}
	if EqualPtr(ToPtr(T(123)), ToPtr(T(456))) {
		t.Errorf("expected false (val != val)")
	}
}
