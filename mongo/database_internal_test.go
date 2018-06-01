// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"testing"

	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/internal/testutil"
	"github.com/stretchr/testify/require"
)

func createTestDatabase(t *testing.T, name *string) *Database {
	if name == nil {
		db := testutil.DBName(t)
		name = &db
	}

	client := createTestClient(t)
	return client.Database(*name)
}

func TestDatabase_initialize(t *testing.T) {
	t.Parallel()

	name := "foo"

	db := createTestDatabase(t, &name)
	require.Equal(t, db.name, name)
	require.NotNil(t, db.client)
}

func TestDatabase_RunCommand(t *testing.T) {
	t.Parallel()

	db := createTestDatabase(t, nil)

	result, err := db.RunCommand(context.Background(), bson.NewDocument(bson.EC.Int32("ismaster", 1)))
	require.NoError(t, err)

	isMaster, err := result.Lookup("ismaster")
	require.NoError(t, err)
	require.Equal(t, isMaster.Value().Type(), bson.TypeBoolean)
	require.Equal(t, isMaster.Value().Boolean(), true)

	ok, err := result.Lookup("ok")
	require.NoError(t, err)
	require.Equal(t, ok.Value().Type(), bson.TypeDouble)
	require.Equal(t, ok.Value().Double(), 1.0)
}

func TestDatabase_Drop(t *testing.T) {
	t.Parallel()

	name := "TestDatabase_Drop"

	db := createTestDatabase(t, &name)

	client := createTestClient(t)
	err := db.Drop(context.Background())
	require.NoError(t, err)
	list, err := client.ListDatabaseNames(context.Background(), nil)

	require.NoError(t, err)
	require.NotContains(t, list, name)

}

// creates 1 normal collection and 1 capped collection of size 64*1024
func setupListCollectionsDb(db *Database) (uncappedName string, cappedName string, err error) {
	uncappedName, cappedName = "listcoll_uncapped", "listcoll_capped"
	uncappedColl := db.Collection(uncappedName)

	_, err = db.RunCommand(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("create", cappedName),
			bson.EC.Boolean("capped", true),
			bson.EC.Int32("size", 64*1024),
		),
	)
	if err != nil {
		return "", "", err
	}
	cappedColl := db.Collection(cappedName)

	id := objectid.New()
	want := bson.EC.ObjectID("_id", id)
	doc := bson.NewDocument(want, bson.EC.Int32("x", 1))

	_, err = uncappedColl.InsertOne(context.Background(), doc)
	if err != nil {
		return "", "", err
	}

	_, err = cappedColl.InsertOne(context.Background(), doc)
	if err != nil {
		return "", "", err
	}

	return uncappedName, cappedName, nil
}

// verifies both collection names are found in cursor, cursor does not have extra collections, and cursor has no
// duplicates
func verifyListCollections(cursor Cursor, uncappedName string, cappedName string, cappedOnly bool) (err error) {
	var uncappedFound bool
	var cappedFound bool

	for cursor.Next(context.Background()) {
		next := bson.NewDocument()
		err = cursor.Decode(next)
		if err != nil {
			return err
		}

		elem, err := next.LookupErr("name")
		if err != nil {
			return err
		}

		if elem.Type() != bson.TypeString {
			return fmt.Errorf("incorrect type for 'name'. got %v. want %v", elem.Type(), bson.TypeString)
		}

		elemName := elem.StringValue()

		if elemName != uncappedName && elemName != cappedName {
			return fmt.Errorf("incorrect collection name. got: %s. wanted: %s or %s", elemName, uncappedName,
				cappedName)
		}

		if elemName == uncappedName && !uncappedFound {
			if cappedOnly {
				return fmt.Errorf("found uncapped collection %s. expected only capped collections", uncappedName)
			}

			uncappedFound = true
			continue
		}

		if elemName == cappedName && !cappedFound {
			cappedFound = true
			continue
		}

		// duplicate found
		return fmt.Errorf("found duplicate collection %s", elemName)
	}

	if !cappedFound {
		return fmt.Errorf("did not find collection %s", cappedName)
	}

	if !cappedOnly && !uncappedFound {
		return fmt.Errorf("did not find collection %s", uncappedName)
	}

	return nil
}

func listCollectionsTest(db *Database, cappedOnly bool) error {
	uncappedName, cappedName, err := setupListCollectionsDb(db)
	if err != nil {
		return err
	}

	var filter *bson.Document
	if cappedOnly {
		filter = bson.NewDocument(
			bson.EC.Boolean("options.capped", true),
		)
	}

	cursor, err := db.ListCollections(context.Background(), filter)
	if err != nil {
		return err
	}

	return verifyListCollections(cursor, uncappedName, cappedName, cappedOnly)
}

func TestDatabase_ListCollections(t *testing.T) {
	// TODO(GODRIVER-272): Add tests for the replica_set topology using the secondary read preference
	t.Parallel()

	var listCollectionsTable = []struct {
		name             string
		expectedTopology string
		cappedOnly       bool
	}{
		{"standalone_nofilter", "server", false},
		{"standalone_filter", "server", true},
		{"replicaset_nofilter", "replica_set", false},
		{"replicaset_filter", "replica_set", true},
		{"sharded_nofilter", "sharded_cluster", false},
		{"sharded_filter", "sharded_cluster", true},
	}

	for _, tt := range listCollectionsTable {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("topology") != tt.expectedTopology {
				t.Skip()
			}
			dbName := "db_list_collections"
			db := createTestDatabase(t, &dbName)

			defer func() {
				err := db.Drop(context.Background())
				require.NoError(t, err)
			}()

			err := listCollectionsTest(db, tt.cappedOnly)
			require.NoError(t, err)
		})
	}
}
