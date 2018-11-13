// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/mongodb/mongo-go-driver/options"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Dialer is used to make network connections.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// BSONAppender is an interface implemented by types that can marshal a
// provided type into BSON bytes and append those bytes to the provided []byte.
// The AppendBSON can return a non-nil error and non-nil []byte. The AppendBSON
// method may also write incomplete BSON to the []byte.
type BSONAppender interface {
	AppendBSON([]byte, interface{}) ([]byte, error)
}

// BSONAppenderFunc is an adapter function that allows any function that
// satisfies the AppendBSON method signature to be used where a BSONAppender is
// used.
type BSONAppenderFunc func([]byte, interface{}) ([]byte, error)

// AppendBSON implements the BSONAppender interface
func (baf BSONAppenderFunc) AppendBSON(dst []byte, val interface{}) ([]byte, error) {
	return baf(dst, val)
}

// MarshalError is returned when attempting to transform a value into a document
// results in an error.
type MarshalError struct {
	Value interface{}
	Err   error
}

// Error implements the error interface.
func (me MarshalError) Error() string {
	return fmt.Sprintf("cannot transform type %s to a *bson.Document", reflect.TypeOf(me.Value))
}

// Pipeline is a type that makes creating aggregation pipelines easier. It is a
// helper and is intended for serializing to BSON.
//
// Example usage:
//
// 		mongo.Pipeline{{
// 			{"$group", bson.D{{"_id", "$state"}, {"totalPop", bson.D{"$sum", "$pop"}}}},
// 			{"$match": bson.D{{"totalPop", bson.D{"$gte", 10*1000*1000}}}},
// 		}}
//
type Pipeline []bson.D

func transformDocument(registry *bsoncodec.Registry, val interface{}) (bson.Doc, error) {
	if registry == nil {
		registry = bson.NewRegistryBuilder().Build()
	}
	if val == nil {
		return bson.Doc{}, nil
	}
	if doc, ok := val.(bson.Doc); ok {
		return doc.Copy(), nil
	}
	if bs, ok := val.([]byte); ok {
		// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
		val = bson.Raw(bs)
	}

	// TODO(skriptble): Use a pool of these instead.
	buf := make([]byte, 0, 256)
	b, err := bson.MarshalAppendWithRegistry(registry, buf[:0], val)
	if err != nil {
		return nil, MarshalError{Value: val, Err: err}
	}
	return bson.ReadDoc(b)
}

func ensureID(d bson.Doc) (bson.Doc, interface{}) {
	var id interface{}

	elem, err := d.LookupElementErr("_id")
	switch err.(type) {
	case nil:
		id = elem
	default:
		oid := objectid.New()
		d = append(d, bson.Elem{"_id", bson.ObjectID(oid)})
		id = oid
	}
	return d, id
}

func ensureDollarKey(doc bson.Doc) error {
	if len(doc) > 0 && !strings.HasPrefix(doc[0].Key, "$") {
		return errors.New("update document must contain key beginning with '$'")
	}
	return nil
}

func transformAggregatePipeline(registry *bsoncodec.Registry, pipeline interface{}) (bson.Arr, error) {
	pipelineArr := bson.Arr{}
	switch t := pipeline.(type) {
	case Pipeline:
		for _, d := range t {
			doc, err := transformDocument(registry, d)
			if err != nil {
				return nil, err
			}
			pipelineArr = append(pipelineArr, bson.Document(doc))
		}
	case bson.Arr:
		pipelineArr = make(bson.Arr, len(t))
		copy(pipelineArr, t)
	case []bson.Doc:
		pipelineArr = bson.Arr{}

		for _, doc := range t {
			pipelineArr = append(pipelineArr, bson.Document(doc))
		}
	case []interface{}:
		pipelineArr = bson.Arr{}

		for _, val := range t {
			doc, err := transformDocument(registry, val)
			if err != nil {
				return nil, err
			}

			pipelineArr = append(pipelineArr, bson.Document(doc))
		}
	default:
		p, err := transformDocument(registry, pipeline)
		if err != nil {
			return nil, err
		}

		for _, elem := range p {
			pipelineArr = append(pipelineArr, elem.Value)
		}
	}

	return pipelineArr, nil
}

// Build the aggregation pipeline for the CountDocument command.
func countDocumentsAggregatePipeline(registry *bsoncodec.Registry, filter interface{}, opts *options.CountOptions) (bson.Arr, error) {
	pipeline := bson.Arr{}
	filterDoc, err := transformDocument(registry, filter)

	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline, bson.Document(bson.Doc{{"$match", bson.Document(filterDoc)}}))

	if opts != nil {
		if opts.Skip != nil {
			pipeline = append(pipeline, bson.Document(bson.Doc{{"$skip", bson.Int64(*opts.Skip)}}))
		}
		if opts.Limit != nil {
			pipeline = append(pipeline, bson.Document(bson.Doc{{"$limit", bson.Int64(*opts.Limit)}}))
		}
	}

	pipeline = append(pipeline, bson.Document(bson.Doc{
		{"$group", bson.Document(bson.Doc{
			{"_id", bson.Null()},
			{"n", bson.Document(bson.Doc{{"$sum", bson.Int32(1)}})},
		})},
	},
	))

	return pipeline, nil
}
