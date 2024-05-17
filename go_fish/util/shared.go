package shared

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchingDocuments struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	PublishedTime primitive.DateTime `bson:"publishedtime"`
	Url           string             `bson:"url"`
	GUID          string             `bson:"guid"`
}

var SearchTerms = []string{
	"Samsung",
	"Steam",
	"NVME",
	"RTX",
	"Lenovo",
	"Ryzen",
	"Sonos",
	"Epic",
}

var FeedUrl = "https://ozbargain.com.au/feed"
