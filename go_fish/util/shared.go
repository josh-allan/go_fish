package shared

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchingDocuments struct {
	Name string             `bson:"name"`
	Time primitive.DateTime `bson:"time"`
	Url  string             `bson:"url"`
}

var SearchTerms = []string{
	"Samsung",
	"Steam",
	"Credit Card",
	"NVME",
	"RTX",
	"Lenovo",
}
