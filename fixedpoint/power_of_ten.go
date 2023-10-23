package fixedpoint

import "golang.org/x/exp/constraints"

// max precision is 36
const (
	minPowerOfTen = -36
	maxPowerOfTen = 36
)

var powerOfTen = map[int64]FixedPoint{
	minPowerOfTen: MustSafeFromString("0.000000000000000000000000000000000001"),
	-35:           MustSafeFromString("0.00000000000000000000000000000000001"),
	-34:           MustSafeFromString("0.0000000000000000000000000000000001"),
	-33:           MustSafeFromString("0.000000000000000000000000000000001"),
	-32:           MustSafeFromString("0.00000000000000000000000000000001"),
	-31:           MustSafeFromString("0.0000000000000000000000000000001"),
	-30:           MustSafeFromString("0.000000000000000000000000000001"),
	-29:           MustSafeFromString("0.00000000000000000000000000001"),
	-28:           MustSafeFromString("0.0000000000000000000000000001"),
	-27:           MustSafeFromString("0.000000000000000000000000001"),
	-26:           MustSafeFromString("0.00000000000000000000000001"),
	-25:           MustSafeFromString("0.0000000000000000000000001"),
	-24:           MustSafeFromString("0.000000000000000000000001"),
	-23:           MustSafeFromString("0.00000000000000000000001"),
	-22:           MustSafeFromString("0.0000000000000000000001"),
	-21:           MustSafeFromString("0.000000000000000000001"),
	-20:           MustSafeFromString("0.00000000000000000001"),
	-19:           MustSafeFromString("0.0000000000000000001"),
	-18:           MustSafeFromString("0.000000000000000001"),
	-17:           MustSafeFromString("0.00000000000000001"),
	-16:           MustSafeFromString("0.0000000000000001"),
	-15:           MustSafeFromString("0.000000000000001"),
	-14:           MustSafeFromString("0.00000000000001"),
	-13:           MustSafeFromString("0.0000000000001"),
	-12:           MustSafeFromString("0.000000000001"),
	-11:           MustSafeFromString("0.00000000001"),
	-10:           MustSafeFromString("0.0000000001"),
	-9:            MustSafeFromString("0.000000001"),
	-8:            MustSafeFromString("0.00000001"),
	-7:            MustSafeFromString("0.0000001"),
	-6:            MustSafeFromString("0.000001"),
	-5:            MustSafeFromString("0.00001"),
	-4:            MustSafeFromString("0.0001"),
	-3:            MustSafeFromString("0.001"),
	-2:            MustSafeFromString("0.01"),
	-1:            MustSafeFromString("0.1"),
	0:             MustSafeFromString("1"),
	1:             MustSafeFromString("10"),
	2:             MustSafeFromString("100"),
	3:             MustSafeFromString("1000"),
	4:             MustSafeFromString("10000"),
	5:             MustSafeFromString("100000"),
	6:             MustSafeFromString("1000000"),
	7:             MustSafeFromString("10000000"),
	8:             MustSafeFromString("100000000"),
	9:             MustSafeFromString("1000000000"),
	10:            MustSafeFromString("10000000000"),
	11:            MustSafeFromString("100000000000"),
	12:            MustSafeFromString("1000000000000"),
	13:            MustSafeFromString("10000000000000"),
	14:            MustSafeFromString("100000000000000"),
	15:            MustSafeFromString("1000000000000000"),
	16:            MustSafeFromString("10000000000000000"),
	17:            MustSafeFromString("100000000000000000"),
	18:            MustSafeFromString("1000000000000000000"),
	19:            MustSafeFromString("10000000000000000000"),
	20:            MustSafeFromString("100000000000000000000"),
	21:            MustSafeFromString("1000000000000000000000"),
	22:            MustSafeFromString("10000000000000000000000"),
	23:            MustSafeFromString("100000000000000000000000"),
	24:            MustSafeFromString("1000000000000000000000000"),
	25:            MustSafeFromString("10000000000000000000000000"),
	26:            MustSafeFromString("100000000000000000000000000"),
	27:            MustSafeFromString("1000000000000000000000000000"),
	28:            MustSafeFromString("10000000000000000000000000000"),
	29:            MustSafeFromString("100000000000000000000000000000"),
	30:            MustSafeFromString("1000000000000000000000000000000"),
	31:            MustSafeFromString("10000000000000000000000000000000"),
	32:            MustSafeFromString("100000000000000000000000000000000"),
	33:            MustSafeFromString("1000000000000000000000000000000000"),
	34:            MustSafeFromString("10000000000000000000000000000000000"),
	35:            MustSafeFromString("100000000000000000000000000000000000"),
	maxPowerOfTen: MustSafeFromString("1000000000000000000000000000000000000"),
}

// PowerOfTen optimized arithmetic performance for 10^n.
func PowerOfTen[T constraints.Signed](n T) FixedPoint {
	nInt64 := int64(n)
	if val, ok := powerOfTen[nInt64]; ok {
		return val
	}
	return powerOfTen[1].PowInt(NewFromInt64(nInt64))
}
