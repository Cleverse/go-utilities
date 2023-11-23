package fixedpoint

import (
	"database/sql/driver"

	"github.com/Cleverse/go-utilities/errors"
)

// UnmarshalJSON implements the json.Unmarshaler interface for json deserialization.
func (f *FixedPoint) UnmarshalJSON(decimalBytes []byte) error {
	if err := f.d.UnmarshalJSON(decimalBytes); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for json serialization.
func (f FixedPoint) MarshalJSON() ([]byte, error) {
	if !f.d.Valid {
		return []byte("null"), nil
	}
	str := "\"" + f.d.Decimal.RoundBank(Precision).String() + "\""
	return []byte(str), nil
}

// MarshalBinary implements the encoding.TextMarshaler interface for text serialization.
func (f FixedPoint) MarshalBinary() ([]byte, error) {
	if !f.d.Valid {
		return nil, nil
	}

	b, err := f.d.Decimal.RoundBank(Precision).MarshalBinary()
	return b, errors.WithStack(err)
}

// UnmarshalBinary implements the encoding.TextUnmarshaler interface for text deserialization.
func (f *FixedPoint) UnmarshalBinary(data []byte) error {
	if len(data) == 0 || data == nil {
		return nil
	}

	if err := f.d.Decimal.UnmarshalBinary(data); err != nil {
		f.d.Valid = false
		return errors.WithStack(err)
	}

	f.d.Valid = true
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for text serialization.
func (f *FixedPoint) MarshalText() ([]byte, error) {
	if !f.d.Valid {
		return []byte{}, nil
	}

	b, err := f.d.Decimal.RoundBank(Precision).MarshalText()
	return b, errors.WithStack(err)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for text deserialization.
func (f *FixedPoint) UnmarshalText(text []byte) error {
	if err := f.d.UnmarshalText(text); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (f *FixedPoint) Scan(value interface{}) error {
	err := f.d.Scan(value)
	return errors.WithStack(err)
}

// Value implements the driver.Valuer interface for database serialization.
func (f FixedPoint) Value() (driver.Value, error) {
	if !f.d.Valid {
		return nil, nil
	}
	return f.d.Decimal.RoundBank(Precision).String(), nil
}
