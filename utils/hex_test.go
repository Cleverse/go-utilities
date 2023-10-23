package utils

import (
	"testing"
)

func TestHex0XPrefix(t *testing.T) {
	type TestCase struct {
		Input       string
		Has0xPrefix bool
		Func        func(string) string
		Expected    string
	}

	testCases := []TestCase{
		{
			Input:       "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: true,
			Func:        Trim0xPrefix,
			Expected:    "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: true,
			Func:        Add0xPrefix,
			Expected:    "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: true,
			Func:        Flip0xPrefix,
			Expected:    "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: false,
			Func:        Trim0xPrefix,
			Expected:    "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: false,
			Func:        Add0xPrefix,
			Expected:    "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: false,
			Func:        Flip0xPrefix,
			Expected:    "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},

		{
			Input:       "0XEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: true,
			Func:        Trim0xPrefix,
			Expected:    "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "0XEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: true,
			Func:        Add0xPrefix,
			Expected:    "0XEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
		{
			Input:       "0XEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
			Has0xPrefix: true,
			Func:        Flip0xPrefix,
			Expected:    "EeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t *testing.T) {
			if isHasPrefix := Has0xPrefix(tc.Input); tc.Has0xPrefix != isHasPrefix {
				t.Errorf("expected Has0xPrefix to be %v, but got %v", tc.Has0xPrefix, isHasPrefix)
			}
			if actual := tc.Func(tc.Input); tc.Expected != actual {
				t.Errorf("expected result from `Func(string) string` to be %v, but got %v", tc.Expected, actual)
			}
		})
	}
}
