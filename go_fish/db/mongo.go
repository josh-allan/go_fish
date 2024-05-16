package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

func loadDotEnv() {
	godotenv.Load("../.env")
	return
}

func (self *MongoClient) connect(incoming_ctx context.Context) {
	loadDotEnv()
	if self.mongoClient != nil {
		return
	}

	ctx, cancel := context.WithTimeout(incoming_ctx, 30*time.Second)
	defer cancel()

	atlas_uri := os.Getenv("ATLAS_URI")

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

func (self MongoClient) disconnectFromMongoDB(incoming_context context.Context) {
	defer func() {
		if err := self.mongoClient.Disconnect(incoming_context); err != nil {
			log.Fatal(err)
		}
	}()
}
