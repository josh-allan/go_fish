package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchingDocuments struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	PublishedTime primitive.DateTime `bson:"publishedtime"`
	Url           string             `bson:"url"`
	GUID          string             `bson:"guid"`
}

type MongoDatastore struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

func (self MongoDatastore) GetAllDocuments(ctx context.Context) ([]Document, error) {
	return self.GetAllDocuments(ctx)
}
