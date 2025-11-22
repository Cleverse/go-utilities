package fixedpoint

import (
	"database/sql/driver"

	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

// FixedPointArray is a sql/driver compatible type for storing a slice of FixedPoint in postgres.
type FixedPointArray []FixedPoint

// Scan implements the sql.Scanner interface for database deserialization.
func (fs *FixedPointArray) Scan(value interface{}) error {
	var strArray pq.StringArray
	if err := strArray.Scan(value); err != nil {
		return errors.WithStack(err)
	}
	*fs = make([]FixedPoint, 0, len(strArray))
	for _, str := range strArray {
		f := New()
		err := f.Scan(str)
		if err != nil {
			return errors.WithStack(err)
		}
		*fs = append(*fs, f)
	}
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (f FixedPointArray) Value() (driver.Value, error) {
	var strArray pq.StringArray
	for _, v := range f {
		s, err := v.Value()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		strArray = append(strArray, s.(string))
	}
	value, err := strArray.Value()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}
