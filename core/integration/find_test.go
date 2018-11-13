// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/command"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/event"
	"github.com/mongodb/mongo-go-driver/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func initMonitor() (chan *event.CommandStartedEvent, chan *event.CommandSucceededEvent, chan *event.CommandFailedEvent, *event.CommandMonitor) {
	startedChan := make(chan *event.CommandStartedEvent, 100)
	succeededChan := make(chan *event.CommandSucceededEvent, 100)
	failedChan := make(chan *event.CommandFailedEvent, 100)
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, cse *event.CommandStartedEvent) {
			startedChan <- cse
		},
		Succeeded: func(ctx context.Context, cse *event.CommandSucceededEvent) {
			succeededChan <- cse
		},
		Failed: func(ctx context.Context, cfe *event.CommandFailedEvent) {
			failedChan <- cfe
		},
	}

	return startedChan, succeededChan, failedChan, monitor
}

func TestFindPassesMaxAwaitTimeMSThroughToGetMore(t *testing.T) {
	startedChan, succeededChan, failedChan, monitor := initMonitor()

	dbName := fmt.Sprintf("mongo-go-driver-%d-find", os.Getpid())
	colName := testutil.ColName(t)

	server, err := testutil.MonitoredTopology(t, dbName, monitor).SelectServer(context.Background(), description.WriteSelector())
	noerr(t, err)

	// create capped collection
	createCmd := bson.Doc{
		{"create", bson.String(colName)},
		{"capped", bson.Boolean(true)},
		{"size", bson.Int32(1000)}}
	_, err = testutil.RunCommand(t, server.Server, dbName, createCmd)
	noerr(t, err)

	// insert some documents
	insertCmd := bson.Doc{
		{"insert", bson.String(colName)},
		{"documents", bson.Array(bson.Arr{
			bson.Document(bson.Doc{{"_id", bson.Int32(1)}}),
			bson.Document(bson.Doc{{"_id", bson.Int32(2)}}),
			bson.Document(bson.Doc{{"_id", bson.Int32(3)}}),
			bson.Document(bson.Doc{{"_id", bson.Int32(4)}}),
			bson.Document(bson.Doc{{"_id", bson.Int32(5)}})})}}
	_, err = testutil.RunCommand(t, server.Server, dbName, insertCmd)

	conn, err := server.Connection(context.Background())
	noerr(t, err)

	// find those documents, setting cursor type to TAILABLEAWAIT
	cursor, err := (&command.Find{
		NS:     command.Namespace{DB: dbName, Collection: colName},
		Filter: bson.Doc{{"_id", bson.Document(bson.Doc{{"$gte", bson.Int32(1)}})}},
		Opts: []bson.Elem{
			{"batchSize", bson.Int32(3)},
			{"tailable", bson.Boolean(true)},
			{"awaitData", bson.Boolean(true)},
		},
		CursorOpts: []bson.Elem{
			{"batchSize", bson.Int32(3)},
			{"maxTimeMS", bson.Int64(250)},
		},
	}).RoundTrip(context.Background(), server.SelectedDescription(), server, conn)
	noerr(t, err)

	// exhaust the cursor, triggering getMore commands
	for i := 0; i < 4; i++ {
		cursor.Next(context.Background())
	}

	// allow for iteration over range chan
	close(startedChan)
	close(succeededChan)
	close(failedChan)

	// no commands should have failed
	if len(failedChan) != 0 {
		t.Errorf("%d command(s) failed", len(failedChan))
	}

	// check that the expected commands were started
	for started := range startedChan {
		switch started.CommandName {
		case "find":
			assert.Equal(t, 3, int(started.Command.Lookup("batchSize").Int32()))
			assert.True(t, started.Command.Lookup("tailable").Boolean())
			assert.True(t, started.Command.Lookup("awaitData").Boolean())
			assert.Equal(t, started.Command.Lookup("maxAwaitTimeMS"), bson.Val{},
				"Should not have sent maxAwaitTimeMS in find command")
		case "getMore":
			assert.Equal(t, 3, int(started.Command.Lookup("batchSize").Int32()))
			assert.Equal(t, 250, int(started.Command.Lookup("maxTimeMS").Int64()),
				"Should have sent maxTimeMS in getMore command")
		default:
			continue
		}
	}

	// to keep track of seen documents
	id := 1

	// check expected commands succeeded
	for succeeded := range succeededChan {
		switch succeeded.CommandName {
		case "find":
			assert.Equal(t, 1, int(succeeded.Reply.Lookup("ok").Double()))

			actual := succeeded.Reply.Lookup("cursor", "firstBatch").Array()

			for _, v := range actual {
				assert.Equal(t, id, int(v.Document().Lookup("_id").Int32()))
				id++
			}
		case "getMore":
			assert.Equal(t, "getMore", succeeded.CommandName)
			assert.Equal(t, 1, int(succeeded.Reply.Lookup("ok").Double()))

			actual := succeeded.Reply.Lookup("cursor", "nextBatch").Array()

			for _, v := range actual {
				assert.Equal(t, id, int(v.Document().Lookup("_id").Int32()))
				id++
			}
		default:
			continue
		}
	}

	if id <= 5 {
		t.Errorf("not all documents returned; last seen id = %d", id-1)
	}
}
