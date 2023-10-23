package fixedpoint

import (
	"testing"

	"gotest.tools/assert"
)

func TestFixedPointArrayScanValue(t *testing.T) {
	q, _ := NewFromString("1")
	d, _ := NewFromString("3")
	value1 := q.Div(d)

	q, _ = NewFromString("2")
	d, _ = NewFromString("3")
	value2 := q.Div(d)

	q, _ = NewFromString("3")
	d, _ = NewFromString("3")
	value3 := q.Div(d)

	fs := FixedPointArray{value1, value2, value3}

	expected1, _ := NewFromString("0.333333333333333333")
	expected2, _ := NewFromString("0.666666666666666667")
	expected3, _ := NewFromString("1")
	expected := FixedPointArray{expected1, expected2, expected3}

	value, err := fs.Value()
	assert.NilError(t, err)

	result := FixedPointArray{}
	result.Scan(value)

	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)
}
