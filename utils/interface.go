package utils

import "reflect"

// IsNil checks if the given interface value is literally nil.
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	return reflect.ValueOf(i).IsNil()
}
