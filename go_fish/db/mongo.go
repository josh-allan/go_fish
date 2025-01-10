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

func NewMongoClient(uri string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return &MongoClient{client: client}, nil
}

func (m *MongoClient) GetCollection(database, collection string) *mongo.Collection {
	return m.client.Database(database).Collection(collection)
}

func (m *MongoClient) GetTerms(ctx context.Context, database, collection string) ([]string, error) {
	coll := m.GetCollection(database, collection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

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

func (m *MongoClient) GetAllDocuments(ctx context.Context, database, collection string) ([]MatchingDocuments, error) {
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

func (m *MongoClient) InsertDocument(ctx context.Context, database, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := m.GetCollection(database, collection)
	return coll.InsertOne(ctx, document)
}

func (m *MongoClient) UpdateDocument(ctx context.Context, database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := m.GetCollection(database, collection)
	return coll.UpdateOne(ctx, filter, update)
}

func (m *MongoClient) DeleteDocument(ctx context.Context, database, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := m.GetCollection(database, collection)
	return coll.DeleteOne(ctx, filter)
}

func (m *MongoClient) Close() error {
	return m.client.Disconnect(context.TODO())
}
