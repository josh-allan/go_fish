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

func (c *MongoClient) connect(incomingCtx context.Context) {
	dbConfig, err := config.LoadConfig()
	if c.mongoClient != nil {
		return
	}

	ctx, cancel := context.WithTimeout(incomingCtx, 30*time.Second)
	defer cancel()

	atlasUri := dbConfig.MongodbAtlasUri

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(atlasUri))
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	c.mongoClient = client

	return
}

func (c *MongoClient) GetCollection(databaseName, collectionName string) *mongo.Collection {
	c.connect(context.TODO())
	return c.mongoClient.Database(databaseName).Collection(collectionName)
}

func (c *MongoClient) GetTerms(incomingCtx context.Context, databaseName, collectionName string) ([]shared.SearchTerms, error) {
	coll := c.GetCollection(databaseName, collectionName)

	filter := bson.D{{Key: "term", Value: bson.D{{Key: "$exists", Value: true}}}}
	cursor, err := coll.Find(incomingCtx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, incomingCtx)

	var SearchTerms []shared.SearchTerms
	for cursor.Next(incomingCtx) {
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

func (c *MongoClient) GetAllDocuments(incomingCtx context.Context, databaseName, collectionName string) ([]shared.MatchingDocuments, error) {
	coll := c.GetCollection(databaseName, collectionName)

	filter := bson.D{}

	cursor, err := coll.Find(incomingCtx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, incomingCtx)

	var documents []shared.MatchingDocuments
	for cursor.Next(incomingCtx) {
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

func (c MongoClient) disconnectFromMongoDB(incomingContext context.Context) {
	defer func() {
		if err := c.mongoClient.Disconnect(incomingContext); err != nil {
			log.Fatal(err)
		}
	}()
}
