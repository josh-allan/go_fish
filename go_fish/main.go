package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gtuk/discordwebhook"
	"github.com/josh-allan/go_fish/config"
	"github.com/josh-allan/go_fish/db"
	"github.com/josh-allan/go_fish/parser"
	shared "github.com/josh-allan/go_fish/util"
)

var lastUpdated *time.Time

func initLogs() {

	config, err := config.LoadConfig()

	LOG_FILE := config.DOT_LOGS + "/gofish.log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

}

func main() {

	initLogs()

	log.Println("Initialising Go Fish")
	log.Println("Logger started")

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
		return
	}
	log.Println("Config loaded successfully")

	// var store db.Datastore
	ctx := context.Background()

	// store.GetAllDocuments(ctx)
	db_client := &db.MongoClient{}

	SearchTerms, err := db_client.GetTerms(ctx, config.MONGODB_DATABASE_NAME, "search_terms")
	log.Println("Loading searching terms:", SearchTerms)
	if err != nil {
		log.Fatal("Error retrieving search terms:", err)
		return
	}

	collection := db_client.GetCollection(config.MONGODB_DATABASE_NAME, config.MONGODB_COLLECTION)
	existingDocuments, err := db_client.GetAllDocuments(ctx, config.MONGODB_DATABASE_NAME, config.MONGODB_COLLECTION)

	if err != nil {
		log.Fatal("Error retrieving existing entries:", err)
	}

	existingIDs := make(map[string]bool)
	for _, doc := range existingDocuments {
		existingIDs[doc.GUID] = true
	}

	interestingSearches := &SearchTerms
	feedUrl := &shared.FeedUrl
	matchedIDs := make(map[string]bool)
	parser.Feed(feedUrl, interestingSearches, lastUpdated, matchedIDs)
	formattedTime := time.Now().Format("02/01/2006, 15:04:05")
	var webhookURL = config.DISCORD_WEBHOOK_URL
	var username = config.DISCORD_USERNAME

	for {
		matchingEntries, _, newMatchedIDs, err := parser.Feed(feedUrl, interestingSearches, nil, matchedIDs)
		if err != nil {
			log.Printf("Error searching feed: %v\n", err)
			continue
		}

		if len(matchingEntries) > 0 {
			for _, entry := range matchingEntries {
				if !existingIDs[entry.GUID] {
					content := fmt.Sprintf("Matching entry found in %s: %s at %s\n", entry.Link, entry.Title, formattedTime)
					message := discordwebhook.Message{
						Username: &username,
						Content:  &content,
					}

					err := discordwebhook.SendMessage(webhookURL, message)
					if err != nil {
						log.Fatal(err)
					}

					matchingDocuments := &shared.MatchingDocuments{
						ID:            primitive.NewObjectID(),
						Name:          entry.Title,
						PublishedTime: primitive.NewDateTimeFromTime(time.Time(*entry.PublishedParsed)),
						Url:           entry.Link,
						GUID:          entry.GUID,
					}
					res, err := collection.InsertOne(ctx, matchingDocuments)
					if err != nil {
						log.Fatalf("Error inserting document: %v", err)
					}
					log.Printf("Document inserted with ID: %s\n", res.InsertedID)
				}
			}
		} else {
			log.Printf("No new matching entries found in %s.\n", *feedUrl)
		}

		// Add the new matched IDs to the matched IDs map
		for _, id := range newMatchedIDs {
			matchedIDs[id] = true
		}

		time.Sleep(60 * time.Second) // Wait 60 seconds for the next run
	}
}
