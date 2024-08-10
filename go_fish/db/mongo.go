package db

import (
	"context"
	"log"
	"time"

	"github.com/josh-allan/go_fish/config"
	shared "github.com/josh-allan/go_fish/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	mongoClient *mongo.Client
}

func (self *MongoClient) connect(incoming_ctx context.Context) {
	config, err := config.LoadConfig()
	if self.mongoClient != nil {
		return
	}

	ctx, cancel := context.WithTimeout(incoming_ctx, 30*time.Second)
	defer cancel()

	atlas_uri := config.MONGODB_ATLAS_URI

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(atlas_uri))
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	self.mongoClient = client

	return
}

func (self *MongoClient) GetCollection(databaseName, collectionName string) *mongo.Collection {
	self.connect(context.TODO())
	return self.mongoClient.Database(databaseName).Collection(collectionName)
}

func (self *MongoClient) GetTerms(incoming_ctx context.Context, databaseName, collectionName string) ([]shared.SearchTerms, error) {
	coll := self.GetCollection(databaseName, collectionName)

	filter := bson.D{{Key: "term", Value: bson.D{{Key: "$exists", Value: true}}}}
	cursor, err := coll.Find(incoming_ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(incoming_ctx)

	var SearchTerms []shared.SearchTerms
	for cursor.Next(incoming_ctx) {
		var term shared.SearchTerms
		if err := cursor.Decode(&term); err != nil {
			return nil, err
		}
		SearchTerms = append(SearchTerms, term)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return SearchTerms, nil
}

func (self *MongoClient) GetAllDocuments(incoming_ctx context.Context, databaseName, collectionName string) ([]shared.MatchingDocuments, error) {
	coll := self.GetCollection(databaseName, collectionName)

	filter := bson.D{}

	cursor, err := coll.Find(incoming_ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(incoming_ctx)

	var documents []shared.MatchingDocuments
	for cursor.Next(incoming_ctx) {
		var doc shared.MatchingDocuments
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (self MongoClient) disconnectFromMongoDB(incoming_context context.Context) {
	defer func() {
		if err := self.mongoClient.Disconnect(incoming_context); err != nil {
			log.Fatal(err)
		}
	}()
}
