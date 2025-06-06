package nullable

type String struct {
	Nullable[string]
}

func FromString(data string) String {
	return String{
		Nullable: From(data),
	}
}

type Int struct {
	Nullable[int]
}

func FromInt(data int) Int {
	return Int{
		Nullable: From(data),
	}
}

type Int8 struct {
	Nullable[int8]
}

func FromInt8(data int8) Int8 {
	return Int8{
		Nullable: From(data),
	}
}

type Int16 struct {
	Nullable[int16]
}

func FromInt16(data int16) Int16 {
	return Int16{
		Nullable: From(data),
	}
}

type Int32 struct {
	Nullable[int32]
}

func FromInt32(data int32) Int32 {
	return Int32{
		Nullable: From(data),
	}
}

type Int64 struct {
	Nullable[int64]
}

func FromInt64(data int64) Int64 {
	return Int64{
		Nullable: From(data),
	}
}

type Uint struct {
	Nullable[uint]
}

func FromUint(data uint) Uint {
	return Uint{
		Nullable: From(data),
	}
}

type Uint8 struct {
	Nullable[uint8]
}

func FromUint8(data uint8) Uint8 {
	return Uint8{
		Nullable: From(data),
	}
}

type Uint16 struct {
	Nullable[uint16]
}

func FromUint16(data uint16) Uint16 {
	return Uint16{
		Nullable: From(data),
	}
}

type Uint32 struct {
	Nullable[uint32]
}

func FromUint32(data uint32) Uint32 {
	return Uint32{
		Nullable: From(data),
	}
}

type Uint64 struct {
	Nullable[uint64]
}

func FromUint64(data uint64) Uint64 {
	return Uint64{
		Nullable: From(data),
	}
}

type Float32 struct {
	Nullable[float32]
}

func FromFloat32(data float32) Float32 {
	return Float32{
		Nullable: From(data),
	}
}

type Float64 struct {
	Nullable[float64]
}

func FromFloat64(data float64) Float64 {
	return Float64{
		Nullable: From(data),
	}
}

type Bool struct {
	Nullable[bool]
}

func FromBool(data bool) Bool {
	return Bool{
		Nullable: From(data),
	}
}
