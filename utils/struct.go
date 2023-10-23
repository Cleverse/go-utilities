package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// PrettyStruct returns pretty json string of the given struct.
func PrettyStruct(s interface{}) string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

// StructHasZero returns true if a field in a struct is zero value(has not initialized)
// if includeNested is true, it will check nested struct (default is false)
func StructHasZero(s any, includeNested ...bool) bool {
	checkNested := false
	if len(includeNested) > 0 {
		checkNested = includeNested[0]
	}
	return len(StructZeroFields(s, checkNested)) > 0
}

// StructZeroFields returns name of fields if that's field in a struct is zero value(has not initialized)
// if checkNested is true, it will check nested struct (default is false)
func StructZeroFields(s any, checkNested ...bool) []string {
	checknested := false
	if len(checkNested) > 0 {
		checknested = checkNested[0]
	}
	value := reflect.ValueOf(s)

	// if it's pointe, then get the underlying element≤
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// support only struct
	if value.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Invalid input type %q, must be struct", value.Kind()))
	}

	t := value.Type()
	tLen := t.NumField()

	// get all fields in the struct
	fields := make([]reflect.StructField, 0, tLen)
	for i := 0; i < tLen; i++ {
		field := t.Field(i)

		// check is not public fields (can exported field)
		if field.PkgPath == "" {
			fields = append(fields, field)
		}
	}

	// find zero fields
	zeroFields := make([]string, 0)
	for _, field := range fields {
		fieldValue := value.FieldByName(field.Name)

		zeroValue := reflect.Zero(fieldValue.Type()).Interface()
		if reflect.DeepEqual(fieldValue.Interface(), zeroValue) {
			zeroFields = append(zeroFields, field.Name)
		}

		// Check nested struct
		if checknested && IsStruct(fieldValue.Interface()) {
			if nestedZero := StructZeroFields(fieldValue.Interface(), checknested); len(nestedZero) > 0 {
				zeroFields = append(zeroFields, field.Name)
			}
			continue
		}
	}

	return zeroFields
}

// IsStruct returns true if the given variable is a struct or *struct.
func IsStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Kind() == reflect.Struct
}

// Merge merges two structs with the same type. merge from -> to.
//
// warning: this function will modify the first struct (to)
func Merge[T comparable](to, from T) T {
	if !IsStruct(to) {
		return to
	}

	t := reflect.ValueOf(to).Elem()
	f := reflect.ValueOf(from).Elem()

	for i := 0; i < t.NumField(); i++ {
		defaultField := t.Field(i)
		newField := f.Field(i)
		if newField.Interface() != reflect.Zero(defaultField.Type()).Interface() {
			defaultField.Set(reflect.ValueOf(newField.Interface()))
		}
	}
	return to
}
