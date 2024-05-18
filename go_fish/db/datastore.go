package db

import "context"

type Document struct {
	ID   string
	Name string
}
type Datastore interface {
	GetAllDocuments(ctx context.Context) ([]Document, error)
	InsertOne(ctx context.Context, doc Document) (string, error)
}
