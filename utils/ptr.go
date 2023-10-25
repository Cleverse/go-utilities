package utils

import "reflect"

// PtrOf returns a pointer to the given value. alias of ToPtr.
func PtrOf[T any](v T) *T {
	return ToPtr(v)
}

// ToPtr returns a pointer to the given value.
func ToPtr[T any](v T) *T {
	return &v
}

// EmptyableToPtr returns a pointer copy of value if it's nonzero.
// Otherwise, returns nil pointer.
func EmptyableToPtr[T any](x T) *T {
	isZero := reflect.ValueOf(&x).Elem().IsZero()
	if isZero {
		return nil
	}

	return &x
}

// DerefPtr dereferences ptr and returns the value it points to if no nil, or else
// returns def.
func DerefPtr[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	return Empty[T]()
}

// DerefPtrOr dereferences ptr and returns the value it points to if no nil, or else
// returns def.
func DerefPtrOr[T any](ptr *T, def T) T {
	if ptr != nil {
		return DerefPtr(ptr)
	}
	return def
}

// FromPtr alias of DerefPtr. returns the pointer value or empty.
func FromPtr[T any](x *T) T {
	return DerefPtr(x)
}

// FromPtrOr alias of DerefPtrOr. returns the pointer value or the fallback value.
func FromPtrOr[T any](x *T, fallback T) T {
	return DerefPtrOr(x, fallback)
}

// Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
func EqualPtr[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}
