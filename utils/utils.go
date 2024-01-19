// Package utils provides Minimalist and zero-dependency utility functions.
package utils

import "fmt"

// Default inspired by Nullish coalescing operator (??) in JavaScript
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Nullish_coalescing
func Default[T comparable](value T, defaultValue T) (result T) {
	// if value != zero
	if value != result {
		return value
	}
	// if value == zero
	return defaultValue
}

// Optional extract optional parameter from variadic function parameter. If parameter is not provided, return zero value of type T and false.
// If parameter is provided, return parameter and true.
//
// Example: func Get(key string, option ...Option) (string, error) { ... }
//
// It's useful to reduce boilerplate code when implementing optional parameter. You won't need to check if parameter is provided or not.
func Optional[T any](opt []T) (optional T, ok bool) {
	if len(opt) > 0 {
		return opt[0], true
	}

	return optional, false
}

// DefaultOptional extract optional parameter from variadic function parameter.
// If parameter is not provided or zero value of type T, return defaultValue.
//
// It's useful to reduce boilerplate code when implementing optional parameter. You won't need to check if parameter is provided or not.
func DefaultOptional[T comparable](opt []T, defaultValue T) (result T) {
	val, ok := Optional[T](opt)
	if ok && val != Empty[T]() {
		return val
	}
	return defaultValue
}

// ToZero alias of Zero. returns zero value of the type of the given value.
func ToZero[T any](value T) (result T) {
	return Zero[T](value)
}

// Zero returns zero value of the type of the given value.
func Zero[T any](value T) T {
	return Empty[T]()
}

// Empty returns an empty value of the given type.
func Empty[T any]() T {
	var zero T
	return zero
}

func msgFormatter(msgAndArgs ...interface{}) string {
	switch len(msgAndArgs) {
	case 0:
		return ""
	case 1:
		if msgAsStr, ok := msgAndArgs[0].(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msgAndArgs[0])
	default:
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
}
