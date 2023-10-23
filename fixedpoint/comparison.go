package fixedpoint

// Cmp compares the numbers represented by d and d2 and returns:
//
//	-1 if d <  d2
//	 0 if d == d2
//	+1 if d >  d2
func (f FixedPoint) Cmp(a FixedPoint) int {
	if !a.IsValid() && !f.IsValid() {
		return 0
	}
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	return f.d.Decimal.Cmp(a.d.Decimal)
}

func (f FixedPoint) Equal(a FixedPoint) bool {
	if !a.IsValid() && !f.IsValid() {
		return true
	}
	if !a.IsValid() || !f.IsValid() {
		panic("FixedPoint is not valid")
	}
	return f.d.Decimal.Equal(a.d.Decimal)
}
