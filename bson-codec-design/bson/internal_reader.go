package bson

import (
	"github.com/mongodb/mongo-go-driver/bson-codec-design/bson/decimal"
	"github.com/mongodb/mongo-go-driver/bson-codec-design/bson/objectid"
)

type ArrayReader interface {
	Next() bool
	ReadValue() (ValueReader, error)
}

type DocumentReader interface {
	ReadElement() (string, ValueReader, error)
}

type ValueReader interface {
	Type() Type
	Skip() error

	ReadArray() (ArrayReader, error)
	ReadBinary() (b []byte, btype byte, err error)
	ReadBoolean() (bool, error)
	ReadDocument() (DocumentReader, error)
	ReadCodeWithScope() (code string, dr DocumentReader, err error)
	ReadDBPointer() (ns string, oid objectid.ObjectID, err error)
	ReadDateTime() (int64, error)
	ReadDecimal128() (decimal.Decimal128, error)
	ReadDouble() (float64, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadJavascript() (code string, err error)
	ReadMaxKey() error
	ReadMinKey() error
	ReadNull() error
	ReadObjectID() (objectid.ObjectID, error)
	ReadRegex() (pattern, options string, err error)
	ReadString() (string, error)
	ReadSymbol() (symbol string, err error)
	ReadTimeStamp(t, i uint32, err error)
	ReadUndefined() error
}

type reader struct {
	b   []byte
	idx int64
}

func (r *reader) Read(p []byte) (int, error) {
	return 0, nil
}

func (r *reader) ReadAt(p []byte, off int64) (int, error) {
	return 0, nil
}
