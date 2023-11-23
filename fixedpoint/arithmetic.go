package fixedpoint

import "github.com/shopspring/decimal"

func (f FixedPoint) Add(a FixedPoint) FixedPoint {
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.Add(a.d.Decimal)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

func (f FixedPoint) Sub(a FixedPoint) FixedPoint {
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.Sub(a.d.Decimal)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

func (f FixedPoint) Mul(a FixedPoint) FixedPoint {
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.Mul(a.d.Decimal)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

func (f FixedPoint) Div(a FixedPoint) FixedPoint {
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.DivRound(a.d.Decimal, DivPrecision)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

// PowerInt returns f^a, where a is an Integer only.
// Panics if a is float/decimal.
//
// Warning: Power with negative exponent is not work normally.
// E.g. 10^-17, 10^-19, 10^-21 will return 0.
//
// Why we not support PowDecimal. ref: https://github.com/shopspring/decimal/issues/201
func (f FixedPoint) PowInt(a FixedPoint) FixedPoint {
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	if !a.IsInteger() {
		panic("FixedPoint is not integer")
	}
	if f.IsZero() && a.IsZero() {
		panic("0^0 is undefined")
	}

	d := f.d.Decimal.Pow(a.d.Decimal)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

func (f FixedPoint) Mod(a FixedPoint) FixedPoint {
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.Mod(a.d.Decimal)
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

func (f FixedPoint) Abs() FixedPoint {
	if !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.Abs()
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}

func (f FixedPoint) Neg() FixedPoint {
	if !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	d := f.d.Decimal.Neg()
	return FixedPoint{
		d: decimal.NewNullDecimal(d),
	}
}
