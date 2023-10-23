package utils

import (
	"fmt"
	"testing"
)

func TestZeroFieldsEqualStructHasZero(t *testing.T) {
	type C struct {
		C string
	}

	type BC struct {
		B string
		C *C
	}

	type ABC struct {
		A string
		B BC
	}

	{
		val := ABC{
			A: "A",
			B: BC{},
		}
		isZero_01 := len(StructZeroFields(val)) > 0
		isZero_02 := StructHasZero(val)
		assertEqual(t, isZero_01, isZero_02)
		assertTrue(t, isZero_01)
	}
	{
		val := ABC{
			A: "A",
			B: BC{
				B: "B",
			},
		}

		{
			isZero_01 := len(StructZeroFields(val)) > 0
			isZero_02 := StructHasZero(val)
			assertEqual(t, isZero_01, isZero_02)
			assertFalse(t, isZero_01)
		}
		{
			isZero_01 := len(StructZeroFields(val, true)) > 0
			isZero_02 := StructHasZero(val, true)
			assertEqual(t, isZero_01, isZero_02)
			assertTrue(t, isZero_01)
		}
	}
	{
		val := ABC{
			A: "A",
			B: BC{
				B: "B",
				C: &C{},
			},
		}

		{
			isZero_01 := len(StructZeroFields(val)) > 0
			isZero_02 := StructHasZero(val)
			assertEqual(t, isZero_01, isZero_02)
			assertFalse(t, isZero_01)
		}
		{
			isZero_01 := len(StructZeroFields(val, true)) > 0
			isZero_02 := StructHasZero(val, true)
			assertEqual(t, isZero_01, isZero_02)
			assertTrue(t, isZero_01)
		}
	}
}

func TestStructZeroFields(t *testing.T) {
	type TestSpec struct {
		ExpectedFields []string
		CheckNested    bool
		Struct         interface{}
	}

	testSpecs := []TestSpec{
		{
			ExpectedFields: []string{"C", "D"},
			Struct: struct {
				A string
				B int
				C bool
				D []string
			}{
				A: "A",
				B: 2,
			},
		},
		{
			ExpectedFields: []string{"B"},
			Struct: struct {
				A string
				B int
				c bool
			}{
				A: "A",
			},
		},
		{
			ExpectedFields: []string{},
			Struct: struct {
				A string
				b int
			}{
				A: "A",
			},
		},
		{
			ExpectedFields: []string{"F", "X"},
			Struct: struct {
				A string
				F *bool
				X struct {
					Y string
					Z *bool
				}
			}{
				A: "A",
			},
		},
		{
			ExpectedFields: []string{"B"},
			CheckNested:    true,
			Struct: struct {
				A string
				B interface{}
			}{
				A: "A",
				B: &struct {
					C []string
				}{
					C: nil,
				},
			},
		},
	}

	for i, testSpec := range testSpecs {
		t.Run(fmt.Sprint("#", i+1), func(t *testing.T) {
			assertEqualAny(t, testSpec.ExpectedFields, StructZeroFields(testSpec.Struct, testSpec.CheckNested))
		})
	}
}

func TestStructHasZero(t *testing.T) {
	type TestSpec struct {
		ExpectedZero bool
		Struct       interface{}
	}

	testSpecs := []TestSpec{
		{
			ExpectedZero: true,
			Struct: struct {
				A string
				B int
				C bool
				D []string
			}{
				A: "A",
				B: 2,
			},
		},
		{
			ExpectedZero: true,
			Struct: struct {
				A string
				F *bool
			}{
				A: "A",
			},
		},
		{
			ExpectedZero: false,
			Struct: struct {
				a string
				b string
			}{
				b: "Lowercase field B",
			},
		},
		{
			ExpectedZero: false,
			Struct: struct {
				a string
				B string
			}{
				B: "Uppercase field B",
			},
		},
		{
			ExpectedZero: false,
			Struct: struct {
				A string
				B interface{}
			}{
				A: "A",
				B: struct {
					C string
				}{
					C: "C",
				},
			},
		},
		{
			ExpectedZero: true,
			Struct: struct {
				A string
				B struct {
					C []string
					D string
				}
			}{
				A: "A",
			},
		},
	}

	for i, testSpec := range testSpecs {
		t.Run(fmt.Sprint("#", i+1), func(t *testing.T) {
			assertEqual(t, testSpec.ExpectedZero, StructHasZero(testSpec.Struct))
		})
	}
}

func TestMerge(t *testing.T) {
	type Data struct {
		A string
		B int
		C *Data
		D *Data
		E *Data
	}

	org := Data{
		A: "A",
		B: 1,
		C: nil,
		D: &Data{
			A: "D",
		},
		E: &Data{
			A: "E",
		},
	}

	to := org
	from := &Data{
		A: "AAAA",
		B: 0,
		C: &Data{
			A: "C",
		},
		D: nil,
		E: &Data{
			A: "EEEE",
		},
	}

	result := Merge(&to, from)

	assertEqual(t, from.A, result.A)
	assertNotEqual(t, org.A, result.A)
	assertEqual(t, org.B, result.B)
	assertNotEqual(t, from.B, result.B)
	assertEqual(t, from.C, result.C)
	assertNotEqual(t, org.C, result.C)
	assertEqual(t, org.D, result.D)
	assertNotEqual(t, from.D, result.D)
	assertEqual(t, from.E, result.E)
	assertNotEqual(t, org.E, result.E)
}
