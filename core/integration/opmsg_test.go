// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package integration

import (
	"context"
	"testing"

	"bytes"
	"errors"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/command"
	"github.com/mongodb/mongo-go-driver/core/connection"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/core/wiremessage"
	"github.com/mongodb/mongo-go-driver/internal"
	"github.com/mongodb/mongo-go-driver/internal/testutil"
)

func createServerConn(t *testing.T) (*topology.SelectedServer, connection.Connection) {
	server, err := testutil.Topology(t).SelectServer(context.Background(), description.WriteSelector())
	noerr(t, err)
	conn, err := server.Connection(context.Background())
	noerr(t, err)

	return server, conn
}

func compareDocs(t *testing.T, reader bson.Raw, doc bson.Doc) {
	marshaled, err := doc.MarshalBSON()
	if err != nil {
		t.Errorf("error marshaling document: %s", err)
	}

	if !bytes.Equal(reader, marshaled) {
		t.Errorf("documents do not match")
	}
}

func createNamespace(t *testing.T) command.Namespace {
	return command.Namespace{
		DB:         dbName,
		Collection: testutil.ColName(t),
	}
}

func compareResults(t *testing.T, channelConn *internal.ChannelConn, docs ...bson.Doc) {
	if len(channelConn.Written) != 1 {
		t.Errorf("expected 1 messages to be sent but got %d", len(channelConn.Written))
	}

	writtenMsg := (<-channelConn.Written).(wiremessage.Msg)
	if len(writtenMsg.Sections) != 2 {
		t.Errorf("expected 2 sections in message. got %d", len(writtenMsg.Sections))
	}

	docSequence := writtenMsg.Sections[1].(wiremessage.SectionDocumentSequence).Documents
	if len(docSequence) != len(docs) {
		t.Errorf("expected %d documents. got %d", len(docs), len(docSequence))
	}

	for i, doc := range docs {
		compareDocs(t, docSequence[i], doc)
	}
}

func createChannelConn() *internal.ChannelConn {
	errChan := make(chan error, 1)
	errChan <- errors.New("read error")

	return &internal.ChannelConn{
		Written:  make(chan wiremessage.WireMessage, 1),
		ReadResp: nil,
		ReadErr:  errChan,
	}
}

func TestOpMsg(t *testing.T) {
	server, _ := createServerConn(t)
	desc := server.Description()

	if desc.WireVersion.Max < wiremessage.OpmsgWireVersion {
		t.Skip("skipping op_msg for wire version < 6")
	}

	t.Run("SingleDocInsert", func(t *testing.T) {
		ctx := context.TODO()
		server, conn := createServerConn(t)
		doc := bson.Doc{{"x", bson.String("testing single doc insert")}}

		cmd := &command.Insert{
			NS: command.Namespace{
				DB:         dbName,
				Collection: testutil.ColName(t),
			},
			Docs: []bson.Doc{doc},
		}

		res, err := cmd.RoundTrip(ctx, server.SelectedDescription(), conn)
		noerr(t, err)

		if len(res.WriteErrors) != 0 {
			t.Errorf("expected no write errors. got %d", len(res.WriteErrors))
		}
	})

	t.Run("SingleDocUpdate", func(t *testing.T) {
		ctx := context.TODO()
		server, conn := createServerConn(t)
		doc := bson.Doc{
			{"$set", bson.Document(bson.Doc{
				{"x", bson.String("updated x")},
			}),
			}}

		updateDocs := []bson.Doc{
			{
				{"q", bson.Document(bson.Doc{})},
				{"u", bson.Document(doc)},
				{"multi", bson.Boolean(true)},
			},
		}

		cmd := &command.Update{
			NS:   createNamespace(t),
			Docs: updateDocs,
		}

		res, err := cmd.RoundTrip(ctx, server.SelectedDescription(), conn)
		noerr(t, err)

		if len(res.WriteErrors) != 0 {
			t.Errorf("expected no write errors. got %d", len(res.WriteErrors))
		}
	})

	t.Run("SingleDocDelete", func(t *testing.T) {
		ctx := context.TODO()
		server, conn := createServerConn(t)
		doc := bson.Doc{{"x", bson.String("testing single doc insert")}}

		deleteDocs := []bson.Doc{
			{
				{"q", bson.Document(doc)},
				{"limit", bson.Int32(0)}},
		}
		cmd := &command.Delete{
			NS:      createNamespace(t),
			Deletes: deleteDocs,
		}

		res, err := cmd.RoundTrip(ctx, server.SelectedDescription(), conn)
		noerr(t, err)

		if len(res.WriteErrors) != 0 {
			t.Errorf("expected no write errors. got %d", len(res.WriteErrors))
		}
	})

	t.Run("MultiDocInsert", func(t *testing.T) {
		ctx := context.TODO()
		server, conn := createServerConn(t)

		doc1 := bson.Doc{{"x", bson.String("testing multi doc insert")}}
		doc2 := bson.Doc{{"y", bson.Int32(50)}}

		cmd := &command.Insert{
			NS: command.Namespace{
				DB:         dbName,
				Collection: testutil.ColName(t),
			},
			Docs: []bson.Doc{doc1, doc2},
		}

		channelConn := createChannelConn()

		_, err := cmd.RoundTrip(ctx, server.SelectedDescription(), channelConn)
		if err == nil {
			t.Errorf("expected read error. got nil")
		}

		compareResults(t, channelConn, doc1, doc2)

		// write to server
		res, err := cmd.RoundTrip(ctx, server.SelectedDescription(), conn)
		noerr(t, err)

		if len(res.WriteErrors) != 0 {
			t.Errorf("expected no write errors. got %d", len(res.WriteErrors))
		}
	})

	t.Run("MultiDocUpdate", func(t *testing.T) {
		ctx := context.TODO()
		server, conn := createServerConn(t)

		doc1 := bson.Doc{
			{"$set", bson.Document(bson.Doc{
				{"x", bson.String("updated x")},
			})},
		}

		doc2 := bson.Doc{
			{"$set", bson.Document(bson.Doc{
				{"y", bson.String("updated y")},
			})},
		}

		updateDocs := []bson.Doc{
			{
				{"q", bson.Document(bson.Doc{})},
				{"u", bson.Document(doc1)},
				{"multi", bson.Boolean(true)},
			},
			{
				{"q", bson.Document(bson.Doc{})},
				{"u", bson.Document(doc2)},
				{"multi", bson.Boolean(true)},
			},
		}

		cmd := &command.Update{
			NS:   createNamespace(t),
			Docs: updateDocs,
		}

		channelConn := createChannelConn()
		_, err := cmd.RoundTrip(ctx, server.SelectedDescription(), channelConn)
		if err == nil {
			t.Errorf("expected read error. got nil")
		}

		compareResults(t, channelConn, updateDocs...)

		// write to server
		res, err := cmd.RoundTrip(ctx, server.SelectedDescription(), conn)
		noerr(t, err)

		if len(res.WriteErrors) != 0 {
			t.Errorf("expected no write errors. got %d", len(res.WriteErrors))
		}
	})

	t.Run("MultiDocDelete", func(t *testing.T) {
		ctx := context.TODO()
		server, conn := createServerConn(t)

		doc1 := bson.Doc{{"x", bson.String("x")}}
		doc2 := bson.Doc{{"y", bson.String("y")}}

		deleteDocs := []bson.Doc{
			{
				{"q", bson.Document(doc1)},
				{"limit", bson.Int32(0)},
			},
			{
				{"q", bson.Document(doc2)},
				{"limit", bson.Int32(0)},
			},
		}

		cmd := &command.Delete{
			NS:      createNamespace(t),
			Deletes: deleteDocs,
		}

		channelConn := createChannelConn()
		_, err := cmd.RoundTrip(ctx, server.SelectedDescription(), channelConn)
		if err == nil {
			t.Errorf("expected read error. got nil")
		}

		compareResults(t, channelConn, deleteDocs...)

		// write to server
		res, err := cmd.RoundTrip(ctx, server.SelectedDescription(), conn)
		noerr(t, err)

		if len(res.WriteErrors) != 0 {
			t.Errorf("expected no write errors. got %d", len(res.WriteErrors))
		}
	})
}
