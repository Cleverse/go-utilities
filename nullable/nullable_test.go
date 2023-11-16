package nullable

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	type args[T Primitive] struct {
		A Nullable[T]
		B Nullable[T]
	}
	type specs[T Primitive] struct {
		Name     string
		Args     args[T]
		Expected bool
	}
	stringTests := []specs[string]{
		{
			Name: "(string) both valid and equal",
			Args: args[string]{
				A: From[string]("Hello"),
				B: From[string]("Hello"),
			},
			Expected: true,
		},
		{
			Name: "(string) both invalid",
			Args: args[string]{
				A: New[string](),
				B: New[string](),
			},
			Expected: true,
		},
		{
			Name: "(string) one valid, one invalid",
			Args: args[string]{
				A: From[string]("Hello"),
				B: New[string](),
			},
			Expected: false,
		},
		{
			Name: "(string) one valid, one invalid",
			Args: args[string]{
				A: New[string](),
				B: From[string]("Hello"),
			},
			Expected: false,
		},
		{
			Name: "(string) both valid, not equal",
			Args: args[string]{
				A: From[string]("Hello"),
				B: From[string]("World"),
			},
			Expected: false,
		},
		{
			Name: "(string) both invalid, not equal data, should be equal",
			Args: args[string]{
				A: Nullable[string]{valid: false, data: "Hello"},
				B: Nullable[string]{valid: false, data: "World"},
			},
			Expected: true,
		},
	}
	intTests := []specs[int]{
		{
			Name: "(int) both valid and equal",
			Args: args[int]{
				A: From[int](1),
				B: From[int](1),
			},
			Expected: true,
		},
		{
			Name: "(int) both invalid",
			Args: args[int]{
				A: New[int](),
				B: New[int](),
			},
			Expected: true,
		},
		{
			Name: "(int) one valid, one invalid",
			Args: args[int]{
				A: From[int](1),
				B: New[int](),
			},
			Expected: false,
		},
		{
			Name: "(int) one valid, one invalid",
			Args: args[int]{
				A: New[int](),
				B: From[int](1),
			},
			Expected: false,
		},
		{
			Name: "(int) both valid, not equal",
			Args: args[int]{
				A: From[int](1),
				B: From[int](2),
			},
			Expected: false,
		},
		{
			Name: "(int) both invalid, not equal data, should be equal",
			Args: args[int]{
				A: Nullable[int]{valid: false, data: 1},
				B: Nullable[int]{valid: false, data: 2},
			},
			Expected: true,
		},
	}
	boolTests := []specs[bool]{
		{
			Name: "(bool) both valid and equal",
			Args: args[bool]{
				A: From[bool](true),
				B: From[bool](true),
			},
			Expected: true,
		},
		{
			Name: "(bool) both invalid",
			Args: args[bool]{
				A: New[bool](),
				B: New[bool](),
			},
			Expected: true,
		},
		{
			Name: "(bool) one valid, one invalid",
			Args: args[bool]{
				A: From[bool](true),
				B: New[bool](),
			},
			Expected: false,
		},
		{
			Name: "(bool) one valid, one invalid",
			Args: args[bool]{
				A: New[bool](),
				B: From[bool](true),
			},
			Expected: false,
		},
		{
			Name: "(bool) both valid, not equal",
			Args: args[bool]{
				A: From[bool](true),
				B: From[bool](false),
			},
			Expected: false,
		},
		{
			Name: "(bool) both invalid, not equal data, should be equal",
			Args: args[bool]{
				A: Nullable[bool]{valid: false, data: true},
				B: Nullable[bool]{valid: false, data: false},
			},
			Expected: true,
		},
	}
	float64Tests := []specs[float64]{
		{
			Name: "(float64) both valid and equal",
			Args: args[float64]{
				A: From[float64](1.23),
				B: From[float64](1.23),
			},
			Expected: true,
		},
		{
			Name: "(float64) both invalid",
			Args: args[float64]{
				A: New[float64](),
				B: New[float64](),
			},
			Expected: true,
		},
		{
			Name: "(float64) one valid, one invalid",
			Args: args[float64]{
				A: From[float64](1.23),
				B: New[float64](),
			},
			Expected: false,
		},
		{
			Name: "(float64) one valid, one invalid",
			Args: args[float64]{
				A: New[float64](),
				B: From[float64](1.23),
			},
			Expected: false,
		},
		{
			Name: "(float64) both valid, not equal",
			Args: args[float64]{
				A: From[float64](1.23),
				B: From[float64](2.45),
			},
			Expected: false,
		},
		{
			Name: "(float64) both invalid, not equal data, should be equal",
			Args: args[float64]{
				A: Nullable[float64]{valid: false, data: 1.23},
				B: Nullable[float64]{valid: false, data: 2.45},
			},
			Expected: true,
		},
	}
	for _, test := range stringTests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Args.A.Equal(test.Args.B))
		})
	}
	for _, test := range intTests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Args.A.Equal(test.Args.B))
		})
	}
	for _, test := range boolTests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Args.A.Equal(test.Args.B))
		})
	}
	for _, test := range float64Tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Args.A.Equal(test.Args.B))
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	type args struct {
		Str   Nullable[string]
		Int   Nullable[int]
		Bool  Nullable[bool]
		Float Nullable[float64]
	}
	type specs struct {
		Name     string
		Args     args
		Expected []byte
	}

	tests := []specs{
		{
			Name:     "normal string",
			Expected: []byte(`{"Str":"Hello","Int":null,"Bool":null,"Float":null}`),
			Args: args{
				Str: From[string]("Hello"),
			},
		},
		{
			Name:     "zero-value string",
			Expected: []byte(`{"Str":"","Int":null,"Bool":null,"Float":null}`),
			Args: args{
				Str: Zero[string](),
			},
		},
		{
			Name:     "null string",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":null,"Float":null}`),
			Args: args{
				Str: New[string](),
			},
		},
		{
			Name:     "normal int",
			Expected: []byte(`{"Str":null,"Int":5,"Bool":null,"Float":null}`),
			Args: args{
				Int: From[int](5),
			},
		},
		{
			Name:     "zero-value int",
			Expected: []byte(`{"Str":null,"Int":0,"Bool":null,"Float":null}`),
			Args: args{
				Int: Zero[int](),
			},
		},
		{
			Name:     "null int",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":null,"Float":null}`),
			Args: args{
				Int: New[int](),
			},
		},
		{
			Name:     "normal bool",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":true,"Float":null}`),
			Args: args{
				Bool: From[bool](true),
			},
		},
		{
			Name:     "zero-value bool",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":false,"Float":null}`),
			Args: args{
				Bool: Zero[bool](),
			},
		},
		{
			Name:     "null bool",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":null,"Float":null}`),
			Args: args{
				Bool: New[bool](),
			},
		},
		{
			Name:     "normal float",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":null,"Float":0.15235}`),
			Args: args{
				Float: From[float64](0.15235),
			},
		},
		{
			Name:     "zero-value float",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":null,"Float":0}`),
			Args: args{
				Float: Zero[float64](),
			},
		},
		{
			Name:     "null float",
			Expected: []byte(`{"Str":null,"Int":null,"Bool":null,"Float":null}`),
			Args: args{
				Float: New[float64](),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := json.Marshal(test.Args)
			assert.NoError(t, err)
			assert.Truef(t, bytes.EqualFold(got, test.Expected), "MarshalJSON() got = %s, want %s", got, test.Expected)
		})
	}

	// test error
}

func TestUnmarshalJSON(t *testing.T) {
	type args struct {
		Str   Nullable[string]
		Int   Nullable[int]
		Bool  Nullable[bool]
		Float Nullable[float64]
	}
	type specs struct {
		Name     string
		JSON     string
		Expected args
	}
	isEqualArgs := func(a, b args) bool {
		return a.Str.Equal(b.Str) && a.Int.Equal(b.Int) && a.Bool.Equal(b.Bool) && a.Float.Equal(b.Float)
	}

	tests := []specs{
		{
			Name: "normal string",
			JSON: `{"Str":"Hello"}`,
			Expected: args{
				Str: From[string]("Hello"),
			},
		},
		{
			Name: "zero-value string",
			JSON: `{"Str":""}`,
			Expected: args{
				Str: Zero[string](),
			},
		},
		{
			Name: "null string",
			JSON: `{"Str":null}`,
			Expected: args{
				Str: New[string](),
			},
		},
		{
			Name: "normal int",
			JSON: `{"Int":5}`,
			Expected: args{
				Int: From[int](5),
			},
		},
		{
			Name: "zero-value int",
			JSON: `{"Int":0}`,
			Expected: args{
				Int: Zero[int](),
			},
		},
		{
			Name: "null int",
			JSON: `{"Int":null}`,
			Expected: args{
				Int: New[int](),
			},
		},
		{
			Name: "normal bool",
			JSON: `{"Bool":true}`,
			Expected: args{
				Bool: From[bool](true),
			},
		},
		{
			Name: "zero-value bool",
			JSON: `{"Bool":false}`,
			Expected: args{
				Bool: Zero[bool](),
			},
		},
		{
			Name: "null bool",
			JSON: `{"Bool":null}`,
			Expected: args{
				Bool: New[bool](),
			},
		},
		{
			Name: "normal float",
			JSON: `{"Float":0.15235}`,
			Expected: args{
				Float: From[float64](0.15235),
			},
		},
		{
			Name: "zero-value float",
			JSON: `{"Float":0}`,
			Expected: args{
				Float: Zero[float64](),
			},
		},
		{
			Name: "null float",
			JSON: `{"Float":null}`,
			Expected: args{
				Float: New[float64](),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var got args
			assert.NoError(t, json.Unmarshal([]byte(test.JSON), &got))
			assert.Truef(t, isEqualArgs(got, test.Expected), "UnmarshalJSON() got = %v, want %v", got, test.Expected)
		})
	}

	// test error
	t.Run("unmarshal error", func(t *testing.T) {
		payload := []byte(`{"Str":true}`)
		var got args
		assert.Error(t, json.Unmarshal(payload, &got))
	})
}
