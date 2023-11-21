package utils

import "testing"

func TestIsNil(t *testing.T) {
	oldIsNil := func(i interface{}) bool {
		return i == nil
	}

	// assign nil value to interface
	var i interface{} = (*int)(nil)

	// oldIsNil() will return false, cause i is not nil (but *int is nil)
	if oldIsNil(i) {
		t.Error("expect oldIsNil() to return false, cause can't use `==` to check nil")
	} else {
		// this is incorrect, cause i is literally nil
	}

	// IsNil() will return true, cause i is not nil (but *int is nil)
	if IsNil(i) {
		// this is correct, cause i is literally nil
	} else {
		t.Errorf("expect IsNil() to return true, cause `i` is %v", i)
	}
}
