package adapter

import (
	"context"
	"time"

	"github.com/arsura/gourney/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
)

type BlogDatabase struct {
	PostCollection MongoCollectionProvider
}

type LogDatabase struct {
	PostLogCollection MongoCollectionProvider
}

type MongoCollections struct {
	*BlogDatabase
	*LogDatabase
}

type MongoClient struct {
	Client *mongo.Client
	Logger *zap.SugaredLogger
	Config *config.Config
}

func NewMongoClient(logger *zap.SugaredLogger, config *config.Config) *MongoClient {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		logger.With("error", err).Panic("failed to new mongodb client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.With("error", err).Panic("failed to connect to mongodb")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.With("error", err).Panic("failed to ping to mongodb")
	}

	return &MongoClient{client, logger, config}
}

func (m *MongoClient) GetMongoCollections() *MongoCollections {
	var (
		blogDatabase = m.Client.Database(m.Config.MongoDB.BlogDatabase.Name)
		logDatabase  = m.Client.Database(m.Config.MongoDB.LogDatabase.Name)
	)

	var (
		postCollection    = blogDatabase.Collection(m.Config.MongoDB.BlogDatabase.Collections.Posts)
		postLogCollection = logDatabase.Collection(m.Config.MongoDB.LogDatabase.Collections.PostLogs)
	)

	return &MongoCollections{
		BlogDatabase: &BlogDatabase{
			PostCollection: postCollection,
		},
		LogDatabase: &LogDatabase{
			PostLogCollection: postLogCollection,
		},
	}
}

func (m *MongoCollections) CreateIndexes() error {
	return nil
}

type MongoCollectionProvider interface {
	// Clone creates a copy of the Collection configured with the given CollectionOptions.
	// The specified options are merged with the existing options on the collection, with the specified options taking
	// precedence.
	Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error)
	// Name returns the name of the collection.
	Name() string
	// Database returns the Database that was used to create the Collection.
	Database() *mongo.Database
	// BulkWrite performs a bulk write operation (https://docs.mongodb.com/manual/core/bulk-write-operations/).
	//
	// The models parameter must be a slice of operations to be executed in this bulk write. It cannot be nil or empty.
	// All of the models must be non-nil. See the mongo.WriteModel documentation for a list of valid model types and
	// examples of how they should be used.
	//
	// The opts parameter can be used to specify options for the operation (see the options.BulkWriteOptions documentation.)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	// InsertOne executes an insert command to insert a single document into the collection.
	//
	// The document parameter must be the document to be inserted. It cannot be nil. If the document does not have an _id
	// field when transformed into BSON, one will be added automatically to the marshalled document. The original document
	// will not be modified. The _id can be retrieved from the InsertedID field of the returned InsertOneResult.
	//
	// The opts parameter can be used to specify options for the operation (see the options.InsertOneOptions documentation.)
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/insert/.
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	// InsertMany executes an insert command to insert multiple documents into the collection. If write errors occur
	// during the operation (e.g. duplicate key error), this method returns a BulkWriteException error.
	//
	// The documents parameter must be a slice of documents to insert. The slice cannot be nil or empty. The elements must
	// all be non-nil. For any document that does not have an _id field when transformed into BSON, one will be added
	// automatically to the marshalled document. The original document will not be modified. The _id values for the inserted
	// documents can be retrieved from the InsertedIDs field of the returned InsertManyResult.
	//
	// The opts parameter can be used to specify options for the operation (see the options.InsertManyOptions documentation.)
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/insert/.
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	// DeleteOne executes a delete command to delete at most one document from the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// deleted. It cannot be nil. If the filter does not match any documents, the operation will succeed and a DeleteResult
	// with a DeletedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
	// matched set.
	//
	// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/delete/.
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	// DeleteMany executes a delete command to delete documents from the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the documents to
	// be deleted. It cannot be nil. An empty document (e.g. bson.D{}) should be used to delete all documents in the
	// collection. If the filter does not match any documents, the operation will succeed and a DeleteResult with a
	// DeletedCount of 0 will be returned.
	//
	// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/delete/.
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	// UpdateByID executes an update command to update the document whose _id value matches the provided ID in the collection.
	// This is equivalent to running UpdateOne(ctx, bson.D{{"_id", id}}, update, opts...).
	//
	// The id parameter is the _id of the document to be updated. It cannot be nil. If the ID does not match any documents,
	// the operation will succeed and an UpdateResult with a MatchedCount of 0 will be returned.
	//
	// The update parameter must be a document containing update operators
	// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be
	// made to the selected document. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// UpdateOne executes an update command to update at most one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
	// with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
	// matched set and MatchedCount will equal 1.
	//
	// The update parameter must be a document containing update operators
	// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be
	// made to the selected document. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// UpdateMany executes an update command to update documents in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the documents to be
	// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
	// with a MatchedCount of 0 will be returned.
	//
	// The update parameter must be a document containing update operators
	// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be made
	// to the selected documents. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// ReplaceOne executes an update command to replace at most one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// replaced. It cannot be nil. If the filter does not match any documents, the operation will succeed and an
	// UpdateResult with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be
	// selected from the matched set and MatchedCount will equal 1.
	//
	// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
	// and cannot contain any update operators (https://docs.mongodb.com/manual/reference/operator/update/).
	//
	// The opts parameter can be used to specify options for the operation (see the options.ReplaceOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
	//
	// The pipeline parameter must be an array of documents, each representing an aggregation stage. The pipeline cannot
	// be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
	// mongo.Pipeline type can be used. See
	// https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/#db-collection-aggregate-stages for a list of
	// valid stages in aggregations.
	//
	// The opts parameter can be used to specify options for the operation (see the options.AggregateOptions documentation.)
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/aggregate/.
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	// CountDocuments returns the number of documents in the collection. For a fast count of the documents in the
	// collection, see the EstimatedDocumentCount method.
	//
	// The filter parameter must be a document and can be used to select which documents contribute to the count. It
	// cannot be nil. An empty document (e.g. bson.D{}) should be used to count all documents in the collection. This will
	// result in a full collection scan.
	//
	// The opts parameter can be used to specify options for the operation (see the options.CountOptions documentation).
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	// EstimatedDocumentCount executes a count command and returns an estimate of the number of documents in the collection
	// using collection metadata.
	//
	// The opts parameter can be used to specify options for the operation (see the options.EstimatedDocumentCountOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/count/.
	EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	// Distinct executes a distinct command to find the unique values for a specified field in the collection.
	//
	// The fieldName parameter specifies the field name for which distinct values should be returned.
	//
	// The filter parameter must be a document containing query operators and can be used to select which documents are
	// considered. It cannot be nil. An empty document (e.g. bson.D{}) should be used to select all documents.
	//
	// The opts parameter can be used to specify options for the operation (see the options.DistinctOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/distinct/.
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	// Find executes a find command and returns a Cursor over the matching documents in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select which documents are
	// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/find/.
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	// FindOne executes a find command and returns a SingleResult for one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// returned. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments will be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The opts parameter can be used to specify options for this operation (see the options.FindOneOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/find/.
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	// FindOneAndDelete executes a findAndModify command to delete at most one document in the collection. and returns the
	// document as it appeared before deletion.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// deleted. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOneAndDeleteOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/findAndModify/.
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	// FindOneAndReplace executes a findAndModify command to replace at most one document in the collection
	// and returns the document as it appeared before replacement.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// replaced. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
	// and cannot contain any update operators (https://docs.mongodb.com/manual/reference/operator/update/).
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOneAndReplaceOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/findAndModify/.
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult
	// FindOneAndUpdate executes a findAndModify command to update at most one document in the collection and returns the
	// document as it appeared before updating.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// updated. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
	// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
	//
	// The update parameter must be a document containing update operators
	// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be made
	// to the selected document. It cannot be nil or empty.
	//
	// The opts parameter can be used to specify options for the operation (see the options.FindOneAndUpdateOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/findAndModify/.
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	// Watch returns a change stream for all changes on the corresponding collection. See
	// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
	//
	// The Collection must be configured with read concern majority or no read concern for a change stream to be created
	// successfully.
	//
	// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
	// nil but can be empty. The stage documents must all be non-nil. See https://docs.mongodb.com/manual/changeStreams/ for
	// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
	// mongo.Pipeline{} type can be used.
	//
	// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
	// documentation).
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	// Indexes returns an IndexView instance that can be used to perform operations on the indexes for the collection.
	Indexes() mongo.IndexView
	// Drop drops the collection on the server. This method ignores "namespace not found" errors so it is safe to drop
	// a collection that does not exist on the server.
	Drop(ctx context.Context) error
}

type MongoClientProvider interface {
	// Connect initializes the Client by starting background monitoring goroutines.
	// If the Client was created using the NewClient function, this method must be called before a Client can be used.
	//
	// Connect starts background goroutines to monitor the state of the deployment and does not do any I/O in the main
	// goroutine. The Client.Ping method can be used to verify that the connection was created successfully.
	Connect(ctx context.Context) error
	// Disconnect closes sockets to the topology referenced by this Client. It will
	// shut down any monitoring goroutines, close the idle connection pool, and will
	// wait until all the in use connections have been returned to the connection
	// pool and closed before returning. If the context expires via cancellation,
	// deadline, or timeout before the in use connections have returned, the in use
	// connections will be closed, resulting in the failure of any in flight read
	// or write operations. If this method returns with no errors, all connections
	// associated with this Client have been closed.
	Disconnect(ctx context.Context) error
	// Ping sends a ping command to verify that the client can connect to the deployment.
	//
	// The rp parameter is used to determine which server is selected for the operation.
	// If it is nil, the client's read preference is used.
	//
	// If the server is down, Ping will try to select a server until the client's server selection timeout expires.
	// This can be configured through the ClientOptions.SetServerSelectionTimeout option when creating a new Client.
	// After the timeout expires, a server selection error is returned.
	//
	// Using Ping reduces application resilience because applications starting up will error if the server is temporarily
	// unavailable or is failing over (e.g. during autoscaling due to a load spike).
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	// StartSession starts a new session configured with the given options.
	//
	// StartSession does not actually communicate with the server and will not error if the client is
	// disconnected.
	//
	// If the DefaultReadConcern, DefaultWriteConcern, or DefaultReadPreference options are not set, the client's read
	// concern, write concern, or read preference will be used, respectively.
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	// ListDatabases executes a listDatabases command and returns the result.
	//
	// The filter parameter must be a document containing query operators and can be used to select which
	// databases are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include
	// all databases.
	//
	// The opts parameter can be used to specify options for this operation (see the options.ListDatabasesOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listDatabases/.
	ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error)
	// ListDatabaseNames executes a listDatabases command and returns a slice containing the names of all of the databases
	// on the server.
	//
	// The filter parameter must be a document containing query operators and can be used to select which databases
	// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
	// databases.
	//
	// The opts parameter can be used to specify options for this operation (see the options.ListDatabasesOptions
	// documentation.)
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listDatabases/.
	ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error)
	// UseSession creates a new Session and uses it to create a new SessionContext, which is used to call the fn callback.
	// The SessionContext parameter must be used as the Context parameter for any operations in the fn callback that should
	// be executed under a session. After the callback returns, the created Session is ended, meaning that any in-progress
	// transactions started by fn will be aborted even if fn returns an error.
	//
	// If the ctx parameter already contains a Session, that Session will be replaced with the newly created one.
	//
	// Any error returned by the fn callback will be returned without any modifications.
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	// UseSessionWithOptions operates like UseSession but uses the given SessionOptions to create the Session.
	UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error
	// Watch returns a change stream for all changes on the deployment. See
	// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
	//
	// The client must be configured with read concern majority or no read concern for a change stream to be created
	// successfully.
	//
	// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
	// nil or empty. The stage documents must all be non-nil. See https://docs.mongodb.com/manual/changeStreams/ for a list
	// of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the mongo.Pipeline{}
	// type can be used.
	//
	// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
	// documentation).
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	// NumberSessionsInProgress returns the number of sessions that have been started for this client but have not been
	// closed (i.e. EndSession has not been called).
	NumberSessionsInProgress() int
}

type MongoDatabaseProvider interface {
	// Client returns the Client the Database was created from.
	Client() *mongo.Client
	// Name returns the name of the database.
	Name() string
	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	// Aggregate executes an aggregate command the database. This requires MongoDB version >= 3.6 and driver version >=
	// 1.1.0.
	//
	// The pipeline parameter must be a slice of documents, each representing an aggregation stage. The pipeline
	// cannot be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
	// mongo.Pipeline type can be used. See
	// https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/#db-aggregate-stages for a list of valid
	// stages in database-level aggregations.
	//
	// The opts parameter can be used to specify options for this operation (see the options.AggregateOptions documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/aggregate/.
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	// RunCommand executes the given command against the database. This function does not obey the Database's read
	// preference. To specify a read preference, the RunCmdOptions.ReadPreference option must be used.
	//
	// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
	// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
	// If the command document contains a session ID or any transaction-specific fields, the behavior is undefined.
	// Specifying API versioning options in the command document and declaring an API version on the client is not supported.
	// The behavior of RunCommand is undefined in this case.
	//
	// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
	RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult
	// RunCommandCursor executes the given command against the database and parses the response as a cursor. If the command
	// being executed does not return a cursor (e.g. insert), the command will be executed on the server and an error will
	// be returned because the server response cannot be parsed as a cursor. This function does not obey the Database's read
	// preference. To specify a read preference, the RunCmdOptions.ReadPreference option must be used.
	//
	// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
	// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
	// If the command document contains a session ID or any transaction-specific fields, the behavior is undefined.
	//
	// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
	RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*mongo.Cursor, error)
	// Drop drops the database on the server. This method ignores "namespace not found" errors so it is safe to drop
	// a database that does not exist on the server.
	Drop(ctx context.Context) error
	// ListCollectionSpecifications executes a listCollections command and returns a slice of CollectionSpecification
	// instances representing the collections in the database.
	//
	// The filter parameter must be a document containing query operators and can be used to select which collections
	// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
	// collections.
	//
	// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listCollections/.
	//
	// BUG(benjirewis): ListCollectionSpecifications prevents listing more than 100 collections per database when running
	// against MongoDB version 2.6.
	ListCollectionSpecifications(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]*mongo.CollectionSpecification, error)
	// ListCollections executes a listCollections command and returns a cursor over the collections in the database.
	//
	// The filter parameter must be a document containing query operators and can be used to select which collections
	// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
	// collections.
	//
	// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listCollections/.
	//
	// BUG(benjirewis): ListCollections prevents listing more than 100 collections per database when running against
	// MongoDB version 2.6.
	ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*mongo.Cursor, error)
	// ListCollectionNames executes a listCollections command and returns a slice containing the names of the collections
	// in the database. This method requires driver version >= 1.1.0.
	//
	// The filter parameter must be a document containing query operators and can be used to select which collections
	// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
	// collections.
	//
	// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listCollections/.
	//
	// BUG(benjirewis): ListCollectionNames prevents listing more than 100 collections per database when running against
	// MongoDB version 2.6.
	ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error)
	// ReadConcern returns the read concern used to configure the Database object.
	ReadConcern() *readconcern.ReadConcern
	// ReadPreference returns the read preference used to configure the Database object.
	ReadPreference() *readpref.ReadPref
	// WriteConcern returns the write concern used to configure the Database object.
	WriteConcern() *writeconcern.WriteConcern
	// Watch returns a change stream for all changes to the corresponding database. See
	// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
	//
	// The Database must be configured with read concern majority or no read concern for a change stream to be created
	// successfully.
	//
	// The pipeline parameter must be a slice of documents, each representing a pipeline stage. The pipeline cannot be
	// nil but can be empty. The stage documents must all be non-nil. See https://docs.mongodb.com/manual/changeStreams/ for
	// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
	// mongo.Pipeline{} type can be used.
	//
	// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
	// documentation).
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	// CreateCollection executes a create command to explicitly create a new collection with the specified name on the
	// server. If the collection being created already exists, this method will return a mongo.CommandError. This method
	// requires driver version 1.4.0 or higher.
	//
	// The opts parameter can be used to specify options for the operation (see the options.CreateCollectionOptions
	// documentation).
	//
	// For more information about the command, see https://docs.mongodb.com/manual/reference/command/create/.
	CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error
	// CreateView executes a create command to explicitly create a view on the server. See
	// https://docs.mongodb.com/manual/core/views/ for more information about views. This method requires driver version >=
	// 1.4.0 and MongoDB version >= 3.4.
	//
	// The viewName parameter specifies the name of the view to create.
	//
	// The viewOn parameter specifies the name of the collection or view on which this view will be created
	//
	// The pipeline parameter specifies an aggregation pipeline that will be exececuted against the source collection or
	// view to create this view.
	//
	// The opts parameter can be used to specify options for the operation (see the options.CreateViewOptions
	// documentation).
	CreateView(ctx context.Context, viewName, viewOn string, pipeline interface{}, opts ...*options.CreateViewOptions) error
}
