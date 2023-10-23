package utils

import (
	"reflect"
)

// Must is used to simplify error handling.
// It's helpful to wraps a call to a function returning a value and an error. and panics if err is error or false.
//
// warning: this is not safe, use with caution!! (avoid to use it's in runtime)
func Must[T any](data T, err any, messageArgs ...interface{}) T {
	must(err, messageArgs...)
	return data
}

// MustNotError is used to simplify error handling.
//
// warning: this is not safe, use with caution!! (avoid to use it's in runtime)
func MustNotError[T any](data T, err error) T {
	if err != nil {
		panic(err)
	}
	return data
}

// UnsafeMust is used to simplify error/ok handling by ignoring it in runtime.
//
// warning: this is not safe, use with caution!!
// be careful when value is pointer, it may be nil. (safe in runtime, but need to check nil before use)
func UnsafeMust[T any, E any](data T, e E) T {
	return data
}

// MustOK is used to simplify ok handling.
// for case ok should be true.
//
// warning: this is not safe, use with caution!! (avoid to use it's in runtime)
func MustOK[T any](data T, ok bool) T {
	if !ok {
		panic("got not ok, but should ok")
	}
	return data
}

// MustNotOK is used to simplify ok handling.
// for case ok should be false.
//
// warning: this is not safe, use with caution!! (avoid to use it's in runtime)
func MustNotOK[T any](data T, ok bool) T {
	if ok {
		panic("got ok, but should not ok")
	}
	return data
}

// must panics if err is error or false.
func must(err any, messageArgs ...interface{}) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case bool:
		if !e {
			panic(Default[string](msgFormatter(messageArgs...), "not ok"))
		}
	case error:
		if e == nil {
			return
		}
		message := msgFormatter(messageArgs...)
		if message != "" {
			panic(message + ": " + e.Error())
		}
		panic(e)
	default:
		panic("must: invalid err type '" + reflect.TypeOf(err).Name() + "', should either be a bool or an error")
	}
}
