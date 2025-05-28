// Package nullable provides a safe way to represent nullable primitive values in Go.
package nullable

import (
	"bytes"
	"encoding/json"

	"github.com/Cleverse/go-utilities/errors"
)

// Primitive is a type constraint for all primitive types, except pointers, slices, maps, channels and structs.
type Primitive interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~bool | ~string
}

var nullBytes = []byte("null")

// Nullable is a generic type that can be used to represent a nullable value. If valid is true, then data is considered non-null.
// If valid is false, then data is considered null. Nullable supports all primitive types, except pointers, slices, maps, channels and structs.
// Nullable supports
type Nullable[T Primitive] struct {
	valid bool
	data  T
}

type String struct {
	Nullable[string]
}

type Int struct {
	Nullable[int]
}

type Int8 struct {
	Nullable[int8]
}

type Int16 struct {
	Nullable[int16]
}

type Int32 struct {
	Nullable[int32]
}

type Int64 struct {
	Nullable[int64]
}

type Uint struct {
	Nullable[uint]
}

type Uint8 struct {
	Nullable[uint8]
}

type Uint16 struct {
	Nullable[uint16]
}

type Uint32 struct {
	Nullable[uint32]
}

type Uint64 struct {
	Nullable[uint64]
}

type Float32 struct {
	Nullable[float32]
}

type Float64 struct {
	Nullable[float64]
}

type Bool struct {
	Nullable[bool]
}

// New returns a new null Nullable.
func New[T Primitive]() Nullable[T] {
	return Nullable[T]{}
}

// Null is an alias for New.
func Null[T Primitive]() Nullable[T] {
	return New[T]()
}

// From returns a non-null Nullable with the given data.
func From[T Primitive](data T) Nullable[T] {
	return Nullable[T]{
		valid: true,
		data:  data,
	}
}

// Zero returns a non-null Nullable with the zero value of the given type.
func Zero[T Primitive]() Nullable[T] {
	return Nullable[T]{
		valid: true,
	}
}

// Get returns the data and a boolean indicating whether the Nullable is considered null or non-null.
// If boolean is false, then Nullable is considered null. If boolean is true, then Nullable is considered non-null.
func (n Nullable[T]) Get() (T, bool) {
	return n.data, n.valid
}

// Data returns the without checking if Nullable is considered null. Only use this if you are sure that Nullable is non-null.
func (n Nullable[T]) Data() T {
	return n.data
}

// Set sets the data and marks it as non-null.
func (n *Nullable[T]) Set(data T) {
	n.valid = true
	n.data = data
}

// SetNull marks the data as null.
func (n *Nullable[T]) SetNull() {
	var zero T
	n.valid = false
	n.data = zero
}

// SetZero sets the data to the zero value of the given type and marks it as non-null.
func (n *Nullable[T]) SetZero() {
	var zero T
	n.valid = true
	n.data = zero
}

// IsValid returns true if the Nullable is non-null.
func (n Nullable[T]) IsValid() bool {
	return n.valid
}

// IsZero returns true if the Nullable is non-null and is the zero value of the given type.
func (n Nullable[T]) IsZero() bool {
	if !n.valid {
		return false
	}
	var zero T
	return n.data == zero
}

// Ptr returns a pointer to the data. If the Nullable is null, then nil is returned.
func (n Nullable[T]) Ptr() *T {
	if !n.valid {
		return nil
	}
	return &n.data
}

// Equal returns true if both Nullable are null or if both Nullable are non-null and have the same data.
func (n Nullable[T]) Equal(other Nullable[T]) bool {
	if !n.valid && !other.valid {
		return true
	}
	if n.valid && other.valid && n.data == other.data {
		return true
	}
	return false
}

// MarshalJSON implements json.Marshaler interface. If the Nullable is considered null, then "null" is returned.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return json.Marshal(nil)
	}
	data, err := json.Marshal(n.data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler interface. If "null" is passed, then the Nullable is marked as null.
// Otherwise, the data is marked as non-null and the data is unmarshalled.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.EqualFold(data, nullBytes) {
		n.valid = false
		return nil
	}
	if err := json.Unmarshal(data, &n.data); err != nil {
		n.valid = false
		return errors.WithStack(err)
	}

	n.valid = true
	return nil
}
