package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

// NewMongoClient creates a new MongoDB client.
//
// Returns a pointer to a MongoClient, or an error if the operation fails.
func NewMongoClient(uri string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return &MongoClient{client: client}, nil
}

// GetCollection returns a collection from the database.
//
// Returns a pointer to a mongo.Collection
func (m *MongoClient) GetCollection(database string, collection string) *mongo.Collection {
	return m.client.Database(database).Collection(collection)
}

// GetTerms retrieves all defined search terms from the database.
//
// Returns a slice of strings, or an error if the operation fails
func (m *MongoClient) GetTerms(ctx context.Context, database string, collection string) ([]string, error) {
	coll := m.GetCollection(database, collection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var terms []string
	for cursor.Next(ctx) {
		var result struct {
			Term string `bson:"term"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		terms = append(terms, result.Term)
	}
	return terms, nil
}

// GetAllDocuments retrieves all documents from the specified collection.
//
// Returns a slice of MatchingDocuments, or an error if the operation fails
func (m *MongoClient) GetAllDocuments(ctx context.Context, database string, collection string) ([]MatchingDocuments, error) {
	coll := m.GetCollection(database, collection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)
	var docs []MatchingDocuments
	for cursor.Next(ctx) {
		var result MatchingDocuments
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		docs = append(docs, result)
	}
	return docs, nil
}

// InsertDocument inserts a document into the specified collection.
//
// Returns a pointer to a mongo.InsertOneResult, or an error if the operation fails
func (m *MongoClient) InsertDocument(ctx context.Context, database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := m.GetCollection(database, collection)
	return coll.InsertOne(ctx, document)
}

// UpdateDocument updates a document in the specified collection.
//
// Returns a pointer to a mongo.UpdateResult, or an error if the operation fails
func (m *MongoClient) UpdateDocument(ctx context.Context, database string, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := m.GetCollection(database, collection)
	return coll.UpdateOne(ctx, filter, update)
}

// DeleteDocument deletes a document from the specified collection.
//
// Returns a pointer to a mongo.DeleteResult, or an error if the operation fails
func (m *MongoClient) DeleteDocument(ctx context.Context, database string, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := m.GetCollection(database, collection)
	return coll.DeleteOne(ctx, filter)
}

// Close closes the MongoDB client.
//
// Returns an error if the operation fails
func (m *MongoClient) Close() error {
	return m.client.Disconnect(context.TODO())
}
