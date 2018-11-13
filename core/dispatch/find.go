// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package dispatch

import (
	"context"

	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/core/command"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/session"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/core/uuid"
	"github.com/mongodb/mongo-go-driver/options"
)

// Find handles the full cycle dispatch and execution of a find command against the provided
// topology.
func Find(
	ctx context.Context,
	cmd command.Find,
	topo *topology.Topology,
	selector description.ServerSelector,
	clientID uuid.UUID,
	pool *session.Pool,
	registry *bsoncodec.Registry,
	opts ...*options.FindOptions,
) (command.Cursor, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return nil, err
	}

	desc := ss.Description()
	conn, err := ss.Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rp, err := getReadPrefBasedOnTransaction(cmd.ReadPref, cmd.Session)
	if err != nil {
		return nil, err
	}
	cmd.ReadPref = rp

	// If no explicit session and deployment supports sessions, start implicit session.
	if cmd.Session == nil && topo.SupportsSessions() {
		cmd.Session, err = session.NewClientSession(pool, clientID, session.Implicit)
		if err != nil {
			return nil, err
		}
	}

	fo := options.MergeFindOptions(opts...)
	if fo.AllowPartialResults != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"allowPartialResults", bson.Boolean(*fo.AllowPartialResults)})
	}
	if fo.BatchSize != nil {
		elem := bson.Elem{"batchSize", bson.Int32(*fo.BatchSize)}
		cmd.Opts = append(cmd.Opts, elem)
		cmd.CursorOpts = append(cmd.CursorOpts, elem)

		if fo.Limit != nil && *fo.BatchSize != 0 && *fo.Limit <= int64(*fo.BatchSize) {
			cmd.Opts = append(cmd.Opts, bson.Elem{"singleBatch", bson.Boolean(true)})
		}
	}
	if fo.Collation != nil {
		if desc.WireVersion.Max < 5 {
			return nil, ErrCollation
		}
		cmd.Opts = append(cmd.Opts, bson.Elem{"collation", bson.Document(fo.Collation.ToDocument())})
	}
	if fo.Comment != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"comment", bson.String(*fo.Comment)})
	}
	if fo.CursorType != nil {
		switch *fo.CursorType {
		case options.Tailable:
			cmd.Opts = append(cmd.Opts, bson.Elem{"tailable", bson.Boolean(true)})
		case options.TailableAwait:
			cmd.Opts = append(cmd.Opts, bson.Elem{"tailable", bson.Boolean(true)}, bson.Elem{"awaitData", bson.Boolean(true)})
		}
	}
	if fo.Hint != nil {
		hintElem, err := interfaceToElement("hint", fo.Hint, registry)
		if err != nil {
			return nil, err
		}

		cmd.Opts = append(cmd.Opts, hintElem)
	}
	if fo.Limit != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"limit", bson.Int64(*fo.Limit)})
	}
	if fo.Max != nil {
		maxElem, err := interfaceToElement("max", fo.Max, registry)
		if err != nil {
			return nil, err
		}

		cmd.Opts = append(cmd.Opts, maxElem)
	}
	if fo.MaxAwaitTime != nil {
		// Specified as maxTimeMS on the in the getMore command and not given in initial find command.
		cmd.CursorOpts = append(cmd.CursorOpts, bson.Elem{"maxTimeMS", bson.Int64(int64(*fo.MaxAwaitTime / time.Millisecond))})
	}
	if fo.MaxTime != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"maxTimeMS", bson.Int64(int64(*fo.MaxTime / time.Millisecond))})
	}
	if fo.Min != nil {
		minElem, err := interfaceToElement("min", fo.Min, registry)
		if err != nil {
			return nil, err
		}

		cmd.Opts = append(cmd.Opts, minElem)
	}
	if fo.NoCursorTimeout != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"noCursorTimeout", bson.Boolean(*fo.NoCursorTimeout)})
	}
	if fo.OplogReplay != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"oplogReplay", bson.Boolean(*fo.OplogReplay)})
	}
	if fo.Projection != nil {
		projElem, err := interfaceToElement("projection", fo.Projection, registry)
		if err != nil {
			return nil, err
		}

		cmd.Opts = append(cmd.Opts, projElem)
	}
	if fo.ReturnKey != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"returnKey", bson.Boolean(*fo.ReturnKey)})
	}
	if fo.ShowRecordID != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"showRecordId", bson.Boolean(*fo.ShowRecordID)})
	}
	if fo.Skip != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"skip", bson.Int64(*fo.Skip)})
	}
	if fo.Snapshot != nil {
		cmd.Opts = append(cmd.Opts, bson.Elem{"snapshot", bson.Boolean(*fo.Snapshot)})
	}
	if fo.Sort != nil {
		sortElem, err := interfaceToElement("sort", fo.Sort, registry)
		if err != nil {
			return nil, err
		}

		cmd.Opts = append(cmd.Opts, sortElem)
	}

	return cmd.RoundTrip(ctx, desc, ss, conn)
}
