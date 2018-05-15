package mongoalt

import (
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/readconcern"
	"github.com/mongodb/mongo-go-driver/core/readpref"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
)

type Database struct {
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
}

func NewDatabase(name string) *Database {
	return &Database{name: name}
}

func (db *Database) Name() string {
	return db.name
}

func (db *Database) Collection(name string) *Collection {
	return &Collection{db: db, name: name}
}
