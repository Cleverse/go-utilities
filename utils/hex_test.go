package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert := assert.New(t)
			assert.Equal(tc.Has0xPrefix, Has0xPrefix(tc.Input), "Has0xPrefix should be equal")
			assert.Equal(tc.Expected, tc.Func(tc.Input), "actual result from `Func(string) string` should equal to expected")
		})
	}
}
