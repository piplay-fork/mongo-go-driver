package bson

import (
	"reflect"
	"strings"
)

// node is a compact representation of an element within a BSON document.
// The first 4 bytes are where the element starts in an underlying []byte. The
// last 4 bytes are where the value for that element begins.
//
// The type of the element can be accessed as `data[n[0]]`. The key of the
// element can be accessed as `data[n[0]+1:n[1]-1]`. This will account for the
// null byte at the end of the c-style string. The value can be accessed as
// `data[n[1]:]`. Since there is no value end byte, an unvalidated document
// could result in parsing errors.
type node [2]uint32

// TODO: This needs to be flushed out much more
type StructTagParser interface {
	ParseStructTags(reflect.StructField) (StructTags, error)
}

type StructTagParserFunc func(reflect.StructField) (StructTags, error)

func (stpf StructTagParserFunc) ParseStructTags(sf reflect.StructField) (StructTags, error) {
	return stpf(sf)
}

// StructTags represents the struct tag fields that the StructCodec uses during
// the encoding and decoding process.
//
// In the case of a struct, the lowercased field name is used as the key for each exported
// field but this behavior may be changed using a struct tag. The tag may also contain flags to
// adjust the marshalling behavior for the field.
//
// The properties are defined below:
//
//     OmitEmpty  Only include the field if it's not set to the zero value for the type or to
//                empty slices or maps.
//
//     MinSize    Marshal an integer of a type larger than 32 bits value as an int32, if that's
//                feasible while preserving the numeric value.
//
//     Truncate   When unmarshaling a BSON double, it is permitted to lose precision to fit within
//                a float32.
//
//     Inline     Inline the field, which must be a struct or a map, causing all of its fields
//                or keys to be processed as if they were part of the outer struct. For maps,
//                keys must not conflict with the bson keys of other struct fields.
//
//     Skip       This struct field should be skipped. This is usually denoted by parsing a "-"
//                for the name.
//
type StructTags struct {
	Name      string
	OmitEmpty bool
	MinSize   bool
	Truncate  bool
	Inline    bool
	Skip      bool
}

// DefaultStructTagParser is the StructTagParser used by the StructCodec by default.
// It will handle the bson struct tag. See the documentation for StructTags to see
// what each of the returned fields means.
//
// If there is no name in the struct tag fields, the struct field name is lowercased.
// The tag formats accepted are:
//
//     "[<key>][,<flag1>[,<flag2>]]"
//
//     `(...) bson:"[<key>][,<flag1>[,<flag2>]]" (...)`
//
var DefaultStructTagParser StructTagParserFunc = func(sf reflect.StructField) (StructTags, error) {
	key := strings.ToLower(sf.Name)
	tag, ok := sf.Tag.Lookup("bson")
	var st StructTags
	switch {
	case ok:
		if tag == "-" || sf.Tag == "-" {
			st.Skip = true
		}
		for idx, str := range strings.Split(tag, ",") {
			if idx == 0 && str != "" {
				key = str
			}
			switch str {
			case "omitempty":
				st.OmitEmpty = true
			case "minsize":
				st.MinSize = true
			case "truncate":
				st.Truncate = true
			case "inline":
				st.Inline = true
			}
		}
	case !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0:
		key = string(sf.Tag)
	}

	st.Name = key

	return st, nil
}
