package extjson_test

import (
	"testing"

	"bytes"

	"github.com/10gen/mongo-go-driver/bson/extjson"
)

func Test_ParseObjectToBuilder(t *testing.T) {
	bson := []byte{
		// length
		0x3f, 0x0, 0x0, 0x0,

		// type - null
		0xa,
		// key - "a"
		0x61, 0x0,

		// type - boolean
		0x8,
		// key - "b"
		0x62, 0x0,
		// value - false
		0x0,

		//
		// ----- begin subarray -----
		//

		// type - array
		0x4,
		// key - "c"
		0x63, 0x0,

		// length
		0x1b, 0x0, 0x0, 0x0,

		// type - string
		0x2,
		// key - "0"
		0x30, 0x0,
		// value - "foo" length
		0x4, 0x0, 0x0, 0x0,
		// value - "foo"
		0x66, 0x6f, 0x6f, 0x0,

		// type - double
		0x1,
		// key - "1"
		0x31, 0x0,
		// value - -2.7
		0x9a, 0x99, 0x99, 0x99, 0x99, 0x99, 0x5, 0xc0,

		// null terminator
		0x0,

		//
		// ----- end subarray -----
		//

		//
		// ----- subdocument -----
		//

		// type - document
		0x3,
		// key - "d"
		0x64, 0x0,

		// length
		0x12, 0x0, 0x0, 0x0,

		// type - int64
		0x12,
		// key - "efg"
		0x65, 0x66, 0x67, 0x0,
		// value - 4,
		0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,

		// null terminator
		0x0,

		//
		// ----- end subdocument ----
		//

		// null terminator
		0x0,
	}

	json := `{ "a": null, "b": false, "c": [ "foo", -2.7 ], "d": { "efg": 4 } }`
	builder, err := extjson.ParseObjectToBuilder(json)
	if err != nil {
		t.Fatalf("unable to parse json string: %s", err)
	}

	buf := make([]byte, len(bson))
	_, err = builder.WriteDocument(buf)
	if err != nil {
		t.Fatalf("unable to write document to bytes buffer: %s", err)
	}

	if !bytes.Equal(bson, buf) {
		t.Fatalf("bytes are unequal:\n%#v\n%#v", bson, buf)
	}
}

func Test_ParseObjectToBuilderSubDoc(t *testing.T) {
	bson := []byte{
		// length
		0x21, 0x0, 0x0, 0x0,

		// type - null
		0xa,
		// key - "a"
		0x61, 0x0,

		// type - boolean
		0x8,
		// key - "b"
		0x62, 0x0,
		// value - false
		0x0,

		//
		// ----- subdocument -----
		//

		// type - document
		0x3,
		// key - "d"
		0x64, 0x0,

		// length
		0x12, 0x0, 0x0, 0x0,

		// type - int64
		0x12,
		// key - "efg"
		0x65, 0x66, 0x67, 0x0,
		// value - 4,
		0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,

		// null terminator
		0x0,

		//
		// ----- end subdocument ----
		//

		// null terminator
		0x0,
	}

	json := `{ "a": null, "b": false, "d": { "efg": 4 } }`
	builder, err := extjson.ParseObjectToBuilder(json)
	if err != nil {
		t.Fatalf("unable to parse json string: %s", err)
	}

	buf := make([]byte, len(bson))
	_, err = builder.WriteDocument(buf)
	if err != nil {
		t.Fatalf("unable to write document to bytes buffer: %s", err)
	}

	if !bytes.Equal(bson, buf) {
		t.Fatalf("bytes are unequal:\n%#v\n%#v", bson, buf)
	}
}

func Test_ParseObjectToBuilderFlat(t *testing.T) {
	bson := []byte{
		// length
		0x2f, 0x0, 0x0, 0x0,

		// type - null
		0xa,
		// key - "a"
		0x61, 0x0,

		// type - boolean
		0x8,
		// key - "b"
		0x62, 0x0,
		// value - false
		0x0,

		// type - string
		0x2,
		//
		// key - "c"
		0x63, 0x0,
		// value - "foo" length
		0x4, 0x0, 0x0, 0x0,
		// value - "foo"
		0x66, 0x6f, 0x6f, 0x0,

		// type - double
		0x1,
		// key - "d"
		0x64, 0x0,
		// value - -2.7
		0x9a, 0x99, 0x99, 0x99, 0x99, 0x99, 0x5, 0xc0,

		// type - int64
		0x12,
		// key - "efg"
		0x65, 0x66, 0x67, 0x0,
		// value - 4,
		0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,

		// null terminator
		0x0,
	}

	json := `{ "a": null, "b": false, "c": "foo", "d": -2.7, "efg": 4 }`
	builder, err := extjson.ParseObjectToBuilder(json)
	if err != nil {
		t.Fatalf("unable to parse json string: %s", err)
	}

	buf := make([]byte, len(bson))
	_, err = builder.WriteDocument(buf)
	if err != nil {
		t.Fatalf("unable to write document to bytes buffer: %s", err)
	}

	if !bytes.Equal(bson, buf) {
		t.Fatalf("bytes are unequal:\n%#v\n%#v", bson, buf)
	}
}
