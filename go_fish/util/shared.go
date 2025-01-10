package shared

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
)

type MatchingDocuments struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	PublishedTime primitive.DateTime `bson:"publishedtime"`
	Url           string             `bson:"url"`
	GUID          string             `bson:"guid"`
}

type SearchTerms struct {
	ID   primitive.ObjectID `bson:"_id"`
	Term string             `bson:"term"`
}

var FeedUrl = "https://ozbargain.com.au/feed"

func InitLogs(logdir string) {

	LogFile := logdir + "/gofish.log"
	logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

}
