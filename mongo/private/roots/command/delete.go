package command

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/connection"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/result"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/topology"
	"github.com/mongodb/mongo-go-driver/mongo/private/roots/wiremessage"
)

// Delete represents the delete command.
//
// The delete command executes a delete with a given set of delete documents
// and options.
type Delete struct {
	NS    Namespace
	Query *bson.Document
	Opts  []options.DeleteOptioner
}

// Encode will encode this command into a wire message for the given server description.
func (d *Delete) Encode(topology.ServerDescription) (wiremessage.WireMessage, error) { return nil, nil }

// Decode will decode the wire message using the provided server description. Errors during decoding
// are deferred until either the Result or Err methods are called.
func (d *Delete) Decode(topology.ServerDescription, wiremessage.WireMessage) *Delete {
	return nil
}

// Result returns the result of a decoded wire message and server description.
func (d *Delete) Result() (result.Delete, error) { return result.Delete{}, nil }

// Err returns the error set on this command.
func (d *Delete) Err() error { return nil }

// Dispatch handles the full cycle dispatch and execution of this command against the provided
// topology.
func (d *Delete) Dispatch(context.Context, topology.Topology) (result.Delete, error) {
	return result.Delete{}, nil
}

// RoundTrip handles the execution of this command using the provided connection.
func (d *Delete) RoundTrip(context.Context, topology.ServerDescription, connection.Connection) (result.Delete, error) {
	return result.Delete{}, nil
}
