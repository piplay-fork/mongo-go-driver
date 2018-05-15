package mongoalt

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/options"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Client struct {
	mclient *mongo.Client
}

func Connect(ctx context.Context, uri string, opts *mongo.ClientOptions) (*Client, error) {
	mclient, err := mongo.Connect(ctx, uri, opts)
	if err != nil {
		return nil, err
	}
	return &Client{mclient: mclient}, nil
}

func NewClient(uri string) (*Client, error) {
	mclient, err := mongo.NewClient(uri)
	if err != nil {
		return nil, err
	}
	return &Client{mclient: mclient}, nil
}

func (c *Client) Connect(ctx context.Context) error {
	return c.mclient.Connect(ctx)
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.mclient.Disconnect(ctx)
}

func (c *Client) ListDatabases(ctx context.Context, filter interface{}) (mongo.ListDatabasesResult, error) {
	return c.mclient.ListDatabases(ctx, filter)
}

func (c *Client) ListDatabaseNames(ctx context.Context, filter interface{}) ([]string, error) {
	return c.mclient.ListDatabaseNames(ctx, filter)
}

func (c *Client) RunCommand(ctx context.Context, db *Database, runCommand interface{}) (bson.Reader, error) {
	mdb := c.mclient.Database(db.name)
	return mdb.RunCommand(ctx, runCommand)
}

func (c *Client) DropDatabase(ctx context.Context, db *Database) error {
	mdb := c.mclient.Database(db.name)
	return mdb.Drop(ctx)
}

func (c *Client) InsertOne(ctx context.Context, coll *Collection, document interface{}, opts ...options.InsertOneOptioner) (*mongo.InsertOneResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.InsertOne(ctx, document, opts...)
}

func (c *Client) InsertMany(ctx context.Context, coll *Collection, documents []interface{}, opts ...options.InsertManyOptioner) (*mongo.InsertManyResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.InsertMany(ctx, documents, opts...)
}

func (c *Client) DeleteOne(ctx context.Context, filter interface{}, coll *Collection, opts ...options.DeleteOptioner) (*mongo.DeleteResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.DeleteOne(ctx, filter, opts...)
}

func (c *Client) DeleteMany(ctx context.Context, filter interface{}, coll *Collection, opts ...options.DeleteOptioner) (*mongo.DeleteResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.DeleteMany(ctx, filter, opts...)
}

func (c *Client) UpdateOne(ctx context.Context, coll *Collection, filter interface{}, update interface{},
	opts ...options.UpdateOptioner) (*mongo.UpdateResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.UpdateOne(ctx, filter, update, opts...)
}

func (c *Client) UpdateMany(ctx context.Context, coll *Collection, filter interface{}, update interface{},
	opts ...options.UpdateOptioner) (*mongo.UpdateResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.UpdateMany(ctx, filter, update, opts...)
}

func (c *Client) ReplaceOne(ctx context.Context, coll *Collection, filter interface{},
	replacement interface{}, opts ...options.ReplaceOptioner) (*mongo.UpdateResult, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.ReplaceOne(ctx, filter, replacement, opts...)
}

func (c *Client) Aggregate(ctx context.Context, coll *Collection, pipeline interface{},
	opts ...options.AggregateOptioner) (mongo.Cursor, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Aggregate(ctx, pipeline, opts...)
}

func (c *Client) Count(ctx context.Context, coll *Collection, filter interface{}, opts ...options.CountOptioner) (int64, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Count(ctx, filter, opts...)
}

func (c *Client) Distinct(ctx context.Context, coll *Collection, fieldName string, filter interface{},
	opts ...options.DistinctOptioner) ([]interface{}, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Distinct(ctx, fieldName, filter, opts...)
}

func (c *Client) Find(ctx context.Context, coll *Collection, filter interface{},
	opts ...options.FindOptioner) (mongo.Cursor, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Find(ctx, filter, opts...)
}

func (c *Client) FindOne(ctx context.Context, coll *Collection, filter interface{},
	opts ...options.FindOneOptioner) *mongo.DocumentResult {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.FindOne(ctx, filter, opts...)
}

func (c *Client) FindOneAndDelete(ctx context.Context, coll *Collection, filter interface{},
	opts ...options.FindOneAndDeleteOptioner) *mongo.DocumentResult {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.FindOneAndDelete(ctx, filter, opts...)
}

func (c *Client) FindOneAndReplace(ctx context.Context, coll *Collection, filter interface{},
	replacement interface{}, opts ...options.FindOneAndReplaceOptioner) *mongo.DocumentResult {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.FindOneAndReplace(ctx, filter, replacement,
		opts...)
}

func (c *Client) FindOneAndUpdate(ctx context.Context, coll *Collection, filter interface{},
	update interface{}, opts ...options.FindOneAndUpdateOptioner) *mongo.DocumentResult {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (c *Client) Watch(ctx context.Context, coll *Collection, pipeline interface{},
	opts ...options.ChangeStreamOptioner) (mongo.Cursor, error) {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Watch(ctx, pipeline, opts...)
}

func (c *Client) Indexes(coll *Collection) mongo.IndexView {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Indexes()
}

func (c *Client) DropCollection(ctx context.Context, coll *Collection) error {
	mcoll := c.mclient.Database(coll.db.name).Collection(coll.name)
	return mcoll.Drop(ctx)
}
