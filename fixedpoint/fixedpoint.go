// Package fixedpoint provides a shopspring/decimal wrapper library for fixed point arithmetic operations.
package fixedpoint

import (
	"math"
	"math/big"
	"strings"

	"github.com/Cleverse/go-utilities/errors"
	"github.com/holiman/uint256"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

const (
	DefaultPrecision    = 18
	DefaultDivPrecision = 36
)

var (
	Min = NewFromFloat64(1e-36)           // smallest possible FixedPoint (1e-36)
	Max = NewFromFloat64(math.MaxFloat64) // largest possible FixedPoint (1.7976931348623157e+308)
)

type FixedPoint struct {
	d decimal.NullDecimal
}

// New returns a new null FixedPoint.
func New() FixedPoint {
	return FixedPoint{}
}

// Zero returns a new FixedPoint with value 0.
func Zero() FixedPoint {
	return NewFromInt64(0)
}

// Empty alias for Zero.
// returns a new FixedPoint with value 0.
func Empty() FixedPoint {
	return Zero()
}

// NewFromInt32 returns a new FixedPoint from an int32.
func NewFromInt32(i32 int32) FixedPoint {
	d := decimal.NewFromInt32(i32)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromInt64 returns a new FixedPoint from an int64.
func NewFromInt64(i64 int64) FixedPoint {
	d := decimal.NewFromInt(i64)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromString safely converts a string to a fixedpoint.FixedPoint.
func NewFromString(s string) (FixedPoint, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return FixedPoint{}, errors.WithStack(err)
	}
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}, nil
}

// SafeNewFromString safe to use when input can be empty string
func SafeNewFromString(s string) (FixedPoint, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return Zero(), nil
	}
	return NewFromString(s)
}

// MustSafeNewFromString safe to use when input can be empty string.
// Will panic if error
func MustSafeFromString(s string) FixedPoint {
	d, err := SafeNewFromString(s)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return d
}

// NewFromFloat32 returns a new FixedPoint from a float32.
func NewFromFloat32(f32 float32) FixedPoint {
	d := decimal.NewFromFloat32(f32)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromFloat64 returns a new FixedPoint from a float64.
func NewFromFloat64(f64 float64) FixedPoint {
	d := decimal.NewFromFloat(f64)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromDecimal returns a new FixedPoint from a decimal.Decimal.
func NewFromDecimal(d decimal.Decimal) FixedPoint {
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromBigInt returns a new FixedPoint from a big.Int.
func NewFromBigInt(bi *big.Int) FixedPoint {
	if bi == nil {
		return New()
	}
	d := decimal.NewFromBigInt(bi, 0)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromBigIntExp returns a new FixedPoint from a big.Int with an exponent.
func NewFromBigIntExp(bi *big.Int, exp int32) FixedPoint {
	if bi == nil {
		return New()
	}
	d := decimal.NewFromBigInt(bi, exp)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromBigFloat returns a new FixedPoint from a big.Float.
func NewFromBigFloat(bf *big.Float) FixedPoint {
	if bf == nil {
		return New()
	}
	s := bf.Text('f', DefaultPrecision)
	d, _ := decimal.NewFromString(s)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromUint256 returns a new FixedPoint from a uint256.Int.
func NewFromUint256(u *uint256.Int) FixedPoint {
	if u == nil {
		return New()
	}
	d := decimal.NewFromBigInt(u.ToBig(), 0)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// NewFromNumeric returns a new FixedPoint from a pgtype.Numeric.
func NewFromNumeric(numeric pgtype.Numeric) (FixedPoint, error) {
	if numeric.Status != pgtype.Present {
		return FixedPoint{}, nil
	}
	d := decimal.NewFromBigInt(numeric.Int, numeric.Exp)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}, nil
}

// IsValid returns true if the FixedPoint is valid and can be used.
func (f FixedPoint) IsValid() bool {
	return f.d.Valid
}

// IsZero returns true if the FixedPoint is zero.
//
// Panics if FixedPoint is not valid.
func (f FixedPoint) IsZero() bool {
	if !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	return f.d.Decimal.IsZero()
}

// IsInteger returns true if the FixedPoint is an integer.
//
// Panics if FixedPoint is not valid.
func (f FixedPoint) IsInteger() bool {
	if !f.IsValid() {
		return false
	}
	return f.d.Decimal.IsInteger()
}

// NumDigits returns the number of digits of the decimal coefficient (d.Value).
//
// Panics if FixedPoint is not valid.
func (f FixedPoint) NumDigits() int {
	if !f.IsValid() {
		return 0
	}
	return f.d.Decimal.NumDigits()
}

// Uint256 returns a uint256.Int representation of the FixedPoint.
func (f FixedPoint) Uint256() *uint256.Int {
	if !f.IsValid() {
		return nil
	}
	val, overflow := uint256.FromBig(f.d.Decimal.BigInt())
	if overflow {
		return nil
	}
	return val
}

// BigInt returns a big.Int representation of the FixedPoint.
func (f FixedPoint) BigInt() *big.Int {
	if !f.IsValid() {
		return nil
	}
	return f.d.Decimal.BigInt()
}

// BigFloat returns a big.Float representation of the FixedPoint.
func (f FixedPoint) BigFloat() *big.Float {
	if !f.IsValid() {
		return nil
	}
	return f.d.Decimal.BigFloat()
}

// Float64 returns a float64 representation of the FixedPoint and a bool indicating whether f represents d exactly.
func (f FixedPoint) Float64() (value float64, exact bool) {
	if !f.IsValid() {
		return 0, false
	}
	return f.d.Decimal.Float64()
}

// InexactFloat64 returns a float64 representation of the FixedPoint.
func (f FixedPoint) InexactFloat64() float64 {
	if !f.IsValid() {
		return 0
	}
	return f.d.Decimal.InexactFloat64()
}

// Numeric returns a pgtype.Numeric representation of the FixedPoint.
func (f FixedPoint) Numeric() (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if !f.IsValid() {
		_ = numeric.Set(nil)
		return numeric, nil
	}
	err := numeric.Set(f.StringFixedBank())
	if err != nil {
		return numeric, errors.Wrap(err, "can't set numeric from string")
	}
	return numeric, nil
}

// Copy returns a copy of the FixedPoint with the same value and exponent,
// but a different pointer to value.
func (f FixedPoint) Copy() FixedPoint {
	if !f.IsValid() {
		return New()
	}
	return FixedPoint{
		d: decimal.NewNullDecimal(f.d.Decimal.Copy()),
	}
}

// NullDecimal returns the decimal.NullDecimal representation of the FixedPoint.
func (f FixedPoint) NullDecimal() decimal.NullDecimal {
	if !f.IsValid() {
		return f.d
	}
	return f.d
}

// Decimal returns the decimal.Decimal representation of the FixedPoint.
func (f FixedPoint) Decimal() decimal.Decimal {
	if !f.IsValid() {
		return decimal.Decimal{}
	}
	return f.d.Decimal
}

// String returns the string representation of the decimal
// with the fixed point.
//
// Example:
//
//	d := New(-12345, -3)
//	println(d.String())
//
// Output:
//
//	-12.345
func (f FixedPoint) String() string {
	if !f.IsValid() {
		return ""
	}
	return f.d.Decimal.String()
}

// StringFixed returns a rounded fixed-point string with places digits after
// the decimal point.
//
// Example:
//
//	NewFromFloat64(5.45).StringFixed() // output: "5.450000000000000000"
//	NewFromFloat64(5.5555555555555555555).StringFixed() // output: "5.555555555555555556"
func (f FixedPoint) StringFixed() string {
	if !f.IsValid() {
		return ""
	}
	return f.d.Decimal.StringFixed(DefaultPrecision)
}

// StringFixedBank returns a banker rounded fixed-point string with places digits
// after the decimal point.
//
// Example:
//
//	NewFromFloat64(5.45).StringFixed() // output: "5.450000000000000000"
//	NewFromFloat64(5.5555555555555555555).StringFixed() // output: "5.555555555555555555"
func (f FixedPoint) StringFixedBank() string {
	if !f.IsValid() {
		return ""
	}
	return f.d.Decimal.StringFixedBank(DefaultPrecision)
}

// StringWithPrecision returns a rounded fixed-point string with given precision digits after
func (f FixedPoint) StringWithPrecision(precision int32) string {
	if !f.IsValid() {
		return ""
	}
	return f.d.Decimal.StringFixed(precision)
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
func (f FixedPoint) Sign() int {
	if !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	return f.d.Decimal.Sign()
}

// IsOverPrecision return true if first significant digit is lower than given precision digit
func (f FixedPoint) IsOverPrecision(precision int) bool {
	precisionCutoff := NewFromBigIntExp(big.NewInt(10), int32(-precision))
	return f.Abs().Cmp(precisionCutoff) < 0
}
