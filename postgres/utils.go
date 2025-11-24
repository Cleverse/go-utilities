package postgres

import (
	"errors"
	"math/big"
	"time"

	"github.com/Cleverse/go-utilities/fixedpoint"
	"github.com/Cleverse/go-utilities/nullable"
	"github.com/google/uuid"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func IsUniqueViolationErr(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

func UUIDToPgUUID(src uuid.UUID) pgtype.UUID {
	if src == uuid.Nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{
		Bytes: src,
		Valid: true,
	}
}

func PgUUIDToUUID(src pgtype.UUID) uuid.UUID {
	if !src.Valid {
		return uuid.Nil
	}
	return src.Bytes
}

func TimeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

func PgTimestamptzToTime(tstz pgtype.Timestamptz) time.Time {
	if !tstz.Valid {
		return time.Time{}
	}
	return tstz.Time
}

func StringToPgText(s string) pgtype.Text {
	return NullableStringToPgText(nullable.FromString(s))
}

func NullableStringToPgText(s nullable.String) pgtype.Text {
	data, ok := s.Get()
	return pgtype.Text{
		String: data,
		Valid:  ok,
	}
}

func PgTextToString(t pgtype.Text) string {
	return PgTextToNullableString(t).Data()
}

func PgTextToNullableString(t pgtype.Text) nullable.String {
	if !t.Valid {
		return nullable.String{}
	}
	return nullable.FromString(t.String)
}

func IntToPgNumeric(n int) pgtype.Numeric {
	return NullableIntToPgNumeric(nullable.FromInt(n))
}

func PgNumericToInt(src pgtype.Numeric) int {
	return PgNumericToNullableInt(src).Data()
}

func NullableIntToPgNumeric(src nullable.Int) pgtype.Numeric {
	if !src.IsValid() {
		return pgtype.Numeric{}
	}
	return pgtype.Numeric{
		Int:   big.NewInt(int64(src.Data())),
		Valid: true,
	}
}

func PgNumericToNullableInt(src pgtype.Numeric) nullable.Int {
	if !src.Valid {
		return nullable.Int{}
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(src.Exp)), nil)

	return nullable.FromInt(int(src.Int.Mul(src.Int, scale).Int64()))
}

func Uint256ToPgNumeric(n *uint256.Int) pgtype.Numeric {
	if n == nil {
		return pgtype.Numeric{}
	}
	return pgtype.Numeric{
		Int:   n.ToBig(),
		Valid: true,
	}
}

func PgNumericToUint256(src pgtype.Numeric) *uint256.Int {
	if !src.Valid {
		return nil
	}
	val := new(uint256.Int)
	val.SetFromBig(src.Int)
	return val
}

func BigIntToPgNumeric(n *big.Int) pgtype.Numeric {
	if n == nil {
		return pgtype.Numeric{}
	}
	return pgtype.Numeric{
		Int:   n,
		Valid: true,
	}
}

func PgNumericToBigInt(src pgtype.Numeric) *big.Int {
	if !src.Valid {
		return nil
	}

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(src.Exp)), nil)
	return src.Int.Mul(src.Int, scale)
}

func FixedPointToPgNumeric(n fixedpoint.FixedPoint) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   n.Decimal().Coefficient(),
		Exp:   n.Decimal().Exponent(),
		Valid: true,
	}
}

func PgNumericToFixedPoint(src pgtype.Numeric) fixedpoint.FixedPoint {
	return fixedpoint.NewFromBigIntExp(src.Int, src.Exp)
}

func Int32ToPgInt4(n int32) pgtype.Int4 {
	return NullableInt32ToPgInt4(nullable.FromInt32(n))
}

func PgInt4ToInt32(src pgtype.Int4) int32 {
	return PgInt4ToNullableInt32(src).Data()
}

func NullableInt32ToPgInt4(n nullable.Int32) pgtype.Int4 {
	data, ok := n.Get()
	return pgtype.Int4{
		Int32: data,
		Valid: ok,
	}
}

func PgInt4ToNullableInt32(src pgtype.Int4) nullable.Int32 {
	if !src.Valid {
		return nullable.Int32{}
	}
	return nullable.FromInt32(src.Int32)
}
