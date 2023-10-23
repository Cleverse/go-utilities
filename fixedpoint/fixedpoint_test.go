package fixedpoint

import (
	"encoding/json"
	"math"
	"math/big"
	"os"
	"testing"

	testify "github.com/stretchr/testify/assert"
	"gotest.tools/assert"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestMarshalJSON(t *testing.T) {
	fp := NewFromInt32(100)
	bytes, err := json.Marshal(fp)
	assert.NilError(t, err)

	result := New()
	err = json.Unmarshal(bytes, &result)
	assert.NilError(t, err)
	assert.DeepEqual(t, result, NewFromInt32(100), cmpEqualFixedPoint)

	var resultNil *FixedPoint
	err = json.Unmarshal(bytes, resultNil)
	assert.ErrorContains(t, err, "json: Unmarshal(nil *fixedpoint.FixedPoint)")
}

func TestMarshalJSONFixBanked(t *testing.T) {
	fp, _ := NewFromString("100.0000000000000000004")
	bytes, err := json.Marshal(fp)
	assert.NilError(t, err)
	result := New()
	err = json.Unmarshal(bytes, &result)
	assert.NilError(t, err)
	expected, _ := NewFromString("100")
	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)

	fp, _ = NewFromString("100.0000000000000000006")
	bytes, err = json.Marshal(fp)
	assert.NilError(t, err)
	result = New()
	err = json.Unmarshal(bytes, &result)
	assert.NilError(t, err)
	expected, _ = NewFromString("100.000000000000000001")
	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)

	fp, _ = NewFromString("100.0000000000000000005")
	bytes, err = json.Marshal(fp)
	assert.NilError(t, err)
	result = New()
	err = json.Unmarshal(bytes, &result)
	assert.NilError(t, err)
	expected, _ = NewFromString("100")
	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)

	fp, _ = NewFromString("100.0000000000000000015")
	bytes, err = json.Marshal(fp)
	assert.NilError(t, err)
	result = New()
	err = json.Unmarshal(bytes, &result)
	assert.NilError(t, err)
	expected, _ = NewFromString("100.000000000000000002")
	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)
}

func TestComparison(t *testing.T) {
	t.Run("Cmp", func(t *testing.T) {
		assert.Check(t, New().Cmp(New()) == 0)
		assert.Check(t, NewFromInt32(1).Cmp(NewFromInt32(2)) == -1)
		assert.Check(t, NewFromInt32(2).Cmp(NewFromInt32(1)) == 1)
		assert.Check(t, NewFromInt32(1).Cmp(NewFromInt32(1)) == 0)
		testify.Panics(t, func() { New().Cmp(NewFromInt32(1)) })
		testify.Panics(t, func() { NewFromInt32(1).Cmp(New()) })
	})
	t.Run("Equal", func(t *testing.T) {
		assert.Check(t, New().Equal(New()))
		assert.Check(t, !NewFromInt32(1).Equal(NewFromInt32(2)))
		assert.Check(t, !NewFromInt32(2).Equal(NewFromInt32(1)))
		assert.Check(t, NewFromInt32(1).Equal(NewFromInt32(1)))
		testify.Panics(t, func() { New().Equal(NewFromInt32(1)) })
		testify.Panics(t, func() { NewFromInt32(1).Equal(New()) })
	})
}

func TestIsValid(t *testing.T) {
	assert.Check(t, !New().IsValid())
	assert.Check(t, NewFromInt32(123).IsValid())
	assert.Check(t, NewFromInt64(123).IsValid())
	assert.Check(t, NewFromBigInt(big.NewInt(123)).IsValid())
	assert.Check(t, NewFromFloat32(123.123).IsValid())
	assert.Check(t, NewFromFloat64(123.123).IsValid())
	assert.Check(t, NewFromBigFloat(big.NewFloat(123.123)).IsValid())
	assert.Check(t, NewFromInt32(32).Add(NewFromInt64(64)).IsValid())
	assert.Check(t, NewFromFloat32(32.32).Sub(NewFromInt64(64)).IsValid())
	{
		fp, err := NewFromString("123")
		assert.NilError(t, err)
		assert.Check(t, fp.IsValid())
	}
	{
		fp, err := NewFromString("abc")
		assert.ErrorContains(t, err, "can't convert")
		assert.Check(t, !fp.IsValid())
	}
}

func TestIsZero(t *testing.T) {
	assert.Check(t, NewFromInt32(0).IsZero())
	assert.Check(t, !NewFromInt32(1).IsZero())
	assert.Check(t, NewFromFloat32(0).IsZero())
	assert.Check(t, !NewFromFloat32(1).IsZero())
	assert.Check(t, NewFromInt32(5).Sub(NewFromInt32(5)).IsZero())
	testify.Panics(t, func() { New().IsZero() })
}

func TestNonMutating(t *testing.T) {
	fp := NewFromInt64(0)
	assert.Check(t, fp.IsZero())
	fp.Add(NewFromInt64(5))
	assert.Check(t, fp.IsZero())
}

func TestInteger(t *testing.T) {
	assert.Check(t, NewFromInt32(2).IsInteger())
	assert.Check(t, !NewFromFloat32(2.2).IsInteger())
	assert.Check(t, !New().IsInteger())
}

func TestNumDigits(t *testing.T) {
	assert.DeepEqual(t, NewFromInt32(123).NumDigits(), 3)
	assert.DeepEqual(t, NewFromInt32(0).NumDigits(), 1)
	assert.DeepEqual(t, NewFromInt32(1).NumDigits(), 1)
	assert.DeepEqual(t, NewFromFloat32(1.2).NumDigits(), 2)
	assert.DeepEqual(t, NewFromFloat32(0.1+0.2).NumDigits(), 1)
	assert.DeepEqual(t, NewFromFloat32(0.3).NumDigits(), 1)
	assert.DeepEqual(t, NewFromFloat32(0.34).NumDigits(), 2)
	assert.DeepEqual(t, New().NumDigits(), 0)
}

func TestNewFixedPoint(t *testing.T) {
	assert.DeepEqual(t, New(), FixedPoint{}, cmpEqualFixedPoint)
	assert.DeepEqual(t, NewFromInt32(123), FixedPoint{
		d: decimal.NullDecimal{
			Decimal: decimal.NewFromInt32(123),
			Valid:   true,
		},
	}, cmpEqualFixedPoint)
	assert.DeepEqual(t, NewFromInt64(123), FixedPoint{
		d: decimal.NullDecimal{
			Decimal: decimal.NewFromInt(123),
			Valid:   true,
		},
	}, cmpEqualFixedPoint)
	assert.DeepEqual(t, NewFromInt32(123), NewFromInt64(123), cmpEqualFixedPoint)
	assert.DeepEqual(t, NewFromFloat32(math.Pi), FixedPoint{
		d: decimal.NullDecimal{
			Decimal: decimal.NewFromFloat32(math.Pi),
			Valid:   true,
		},
	}, cmpEqualFixedPoint)
	assert.DeepEqual(t, NewFromFloat64(math.Pi), FixedPoint{
		d: decimal.NullDecimal{
			Decimal: decimal.NewFromFloat(math.Pi),
			Valid:   true,
		},
	}, cmpEqualFixedPoint)
	assert.DeepEqual(t, NewFromFloat64(math.Pi), NewFromFloat32(math.Pi), cmpApproximateFixedPoint)
}

func TestCopy(t *testing.T) {
	d1, _ := NewFromString("100.01")
	d2 := d1.Copy()
	d1 = d1.Neg()
	expectedD1, _ := NewFromString("-100.01")
	expectedD2, _ := NewFromString("100.01")
	assert.DeepEqual(t, d1, expectedD1, cmpEqualFixedPoint)
	assert.DeepEqual(t, d2, expectedD2, cmpEqualFixedPoint)

	d3 := New()
	d4 := d3.Copy()
	d3, _ = NewFromString("100.02")
	expectedD3, _ := NewFromString("100.02")
	expectedD4 := New()
	assert.DeepEqual(t, d3, expectedD3, cmpEqualFixedPoint)
	assert.DeepEqual(t, d4, expectedD4, cmpEqualFixedPoint)
}

func TestArithmetic(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(10).Add(NewFromInt32(20)), NewFromInt32(30), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromFloat64(0.1).Add(NewFromFloat64(0.2)), NewFromFloat64(0.3), cmpApproximateFixedPoint)
		testify.Panics(t, func() { New().Add(NewFromInt32(32)) })
		testify.Panics(t, func() { NewFromInt32(123).Add(New()) })
		testify.Panics(t, func() { New().Add(New()) })
	})

	t.Run("Sub", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(10).Sub(NewFromInt32(20)), NewFromInt32(-10), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromFloat64(0.1).Sub(NewFromFloat64(0.2)), NewFromFloat64(-0.1), cmpApproximateFixedPoint)
		testify.Panics(t, func() { New().Sub(NewFromInt32(32)) })
		testify.Panics(t, func() { NewFromInt32(123).Sub(New()) })
		testify.Panics(t, func() { New().Sub(New()) })
	})

	t.Run("Mul", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(10).Mul(NewFromInt32(20)), NewFromInt32(200), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromFloat64(0.1).Mul(NewFromFloat64(0.2)), NewFromFloat64(0.02), cmpApproximateFixedPoint)
		testify.Panics(t, func() { New().Mul(NewFromInt32(32)) })
		testify.Panics(t, func() { NewFromInt32(123).Mul(New()) })
		testify.Panics(t, func() { New().Mul(New()) })
	})

	t.Run("Div", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(10).Div(NewFromInt32(20)), NewFromFloat32(0.5), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromFloat64(0.1).Div(NewFromFloat64(0.2)), NewFromFloat64(0.5), cmpApproximateFixedPoint)
		testify.Panics(t, func() { New().Div(NewFromInt32(32)) })
		testify.Panics(t, func() { NewFromInt32(123).Div(New()) })
		testify.Panics(t, func() { New().Div(New()) })
		{
			expected, _ := NewFromString(".333333333333333333333333333333333333")
			assert.DeepEqual(t, NewFromInt32(1).Div(NewFromInt32(3)), expected, cmpEqualFixedPoint)
		}
		{
			expected, _ := NewFromString(".666666666666666666666666666666666667")
			assert.DeepEqual(t, NewFromInt32(2).Div(NewFromInt32(3)), expected, cmpEqualFixedPoint)
		}
	})

	t.Run("PowInt", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(2).PowInt(NewFromInt32(10)), NewFromFloat32(1024), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromInt32(2).PowInt(NewFromInt32(0)), NewFromInt32(1), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromInt32(0).PowInt(NewFromInt32(2)), NewFromInt32(0), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromFloat64(2).PowInt(NewFromInt32(-2)), NewFromFloat64(0.25), cmpApproximateFixedPoint)
		assert.DeepEqual(t, NewFromInt64(-2).PowInt(NewFromInt64(2)), NewFromInt64(4), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromInt64(-2).PowInt(NewFromInt64(3)), NewFromInt64(-8), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromInt64(-2).PowInt(NewFromInt64(-3)), NewFromFloat64(-0.125), cmpApproximateFixedPoint)
		testify.Panics(t, func() { New().PowInt(NewFromInt32(32)) })
		testify.Panics(t, func() { NewFromInt32(123).PowInt(New()) })
		testify.Panics(t, func() { New().PowInt(New()) })
		testify.Panics(t, func() { NewFromInt32(9).PowInt(NewFromFloat64(0.5)) })
		testify.Panics(t, func() { NewFromInt32(0).PowInt(NewFromInt32(0)) })
		testify.Panics(t, func() { NewFromInt32(0).PowInt(NewFromInt32(-2)) })
	})

	t.Run("Mod", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(2).Mod(NewFromInt32(10)), NewFromInt32(2), cmpEqualFixedPoint)
		testify.Panics(t, func() { New().Mod(NewFromInt32(32)) })
		testify.Panics(t, func() { NewFromInt32(123).Mod(New()) })
		testify.Panics(t, func() { New().Mod(New()) })
	})

	t.Run("Abs", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(-2).Abs(), NewFromInt32(2), cmpEqualFixedPoint)
		testify.Panics(t, func() { New().Abs() })
	})

	t.Run("Neg", func(t *testing.T) {
		assert.DeepEqual(t, NewFromInt32(-2).Neg(), NewFromInt32(2), cmpEqualFixedPoint)
		assert.DeepEqual(t, NewFromInt32(2).Neg(), NewFromInt32(-2), cmpEqualFixedPoint)
		testify.Panics(t, func() { New().Neg() })
	})
}

func TestFixedPointScanValue(t *testing.T) {
	q, _ := NewFromString("1")
	d, _ := NewFromString("3")
	value := q.Div(d)
	driverValue, _ := value.Value()
	result := New()
	result.Scan(driverValue)
	expected, _ := NewFromString("0.333333333333333333")
	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)

	q, _ = NewFromString("2")
	d, _ = NewFromString("3")
	value = q.Div(d)
	driverValue, _ = value.Value()
	result = New()
	result.Scan(driverValue)
	expected, _ = NewFromString("0.666666666666666667")
	assert.DeepEqual(t, result, expected, cmpEqualFixedPoint)
}

var cmpEqualFixedPoint = cmp.Comparer(func(a FixedPoint, b FixedPoint) bool {
	return a.Equal(b)
})

var cmpApproximateFixedPoint = cmp.Comparer(func(a FixedPoint, b FixedPoint) bool {
	return a.Sub(b).Abs().Cmp(NewFromFloat32(0.00001)) < 1
})
