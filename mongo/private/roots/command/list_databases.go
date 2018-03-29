package command

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/connection"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/topology"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/wiremessage"
)

// ListDatabases represents the listDatabases command.
//
// The listDatabases command lists the databases in a MongoDB deployment.
type ListDatabases struct {
	Opts []options.ListDatabasesOptioner
}

// Encode will encode this command into a wire message for the given server description.
func (ld *ListDatabases) Encode(topology.ServerDescription) (wiremessage.WireMessage, error) {
	return nil, nil
}

// Decode will decode the wire message using the provided server description. Errors during decoding
// are deferred until either the Result or Err methods are called.
func (ld *ListDatabases) Decode(topology.ServerDescription, wiremessage.WireMessage) *ListDatabases {
	return nil
}

// Result returns the result of a decoded wire message and server description.
func (ld *ListDatabases) Result() (Cursor, error) { return nil, nil }

// Err returns the error set on this command.
func (ld *ListDatabases) Err() error { return nil }

// Dispatch handles the full cycle dispatch and execution of this command against the provided
// topology.
func (ld *ListDatabases) Dispatch(context.Context, topology.Topology) (Cursor, error) {
	return nil, nil
}

// RoundTrip handles the execution of this command using the provided connection.
func (ld *ListDatabases) RoundTrip(context.Context, topology.ServerDescription, connection.Connection) (Cursor, error) {
	return nil, nil
}
