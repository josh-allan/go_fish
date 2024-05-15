package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	collection  *mongo.Collection
	ctx         context.Context
)

func loadDotEnv() {
	godotenv.Load("../.env")
	return
}

func ConnectToMongoDB() {
	loadDotEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	atlas_uri := os.Getenv("ATLAS_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(atlas_uri))
	if err != nil {
		log.Fatal(err)
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pinged your deployment. You successfully connected to MongoDB!")

	mongoClient = client
}

func GetMongoDBClient() *mongo.Client {
	ConnectToMongoDB()
	return mongoClient
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return GetMongoDBClient().Database(databaseName).Collection(collectionName)
}
