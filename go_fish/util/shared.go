package shared

import (
	"fmt"
	"github.com/gtuk/discordwebhook"
	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
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

func ConvertToMatchingDocuments(item *gofeed.Item) MatchingDocuments {
	return MatchingDocuments{
		ID:            primitive.NewObjectID(),
		Name:          item.Title,
		PublishedTime: primitive.NewDateTimeFromTime(*item.PublishedParsed),
		Url:           item.Link,
		GUID:          item.GUID,
	}
}
func InitLogs(logdir string) {

	LogFile := logdir + "/gofish.log"
	logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

}
func NotifyDiscord(webhookURL, username string, entry MatchingDocuments, timestamp time.Time) {
	formattedTime := timestamp.Format(time.RFC1123)
	content := fmt.Sprintf("Matching entry found in %s: %s at %s\n", entry.Url, entry.Name, formattedTime)
	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(webhookURL, message)
	if err != nil {
		log.Printf("Error sending message to Discord: %v", err)
	}
}
