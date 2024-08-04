package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gtuk/discordwebhook"
	"github.com/josh-allan/go_fish/config"
	"github.com/josh-allan/go_fish/db"
	"github.com/josh-allan/go_fish/parser"
	shared "github.com/josh-allan/go_fish/util"
)

var lastUpdated *time.Time

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
		return
	}
	log.Println("Config loaded successfully:", config)

	// var store db.Datastore
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// store.GetAllDocuments(ctx)
	db_client := &db.MongoClient{}
	collection := db_client.GetCollection(config.MONGODB_DATABASE_NAME, config.MONGODB_COLLECTION)
	existingDocuments, err := db_client.GetAllDocuments(ctx, config.MONGODB_DATABASE_NAME, config.MONGODB_COLLECTION)
	if err != nil {
		log.Fatal("Error retrieving existing entries:", err)
	}

	existingIDs := make(map[string]bool)
	for _, doc := range existingDocuments {
		existingIDs[doc.GUID] = true
	}

	interestingSearches := &shared.SearchTerms
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

					matchingDocuments := &db.MatchingDocuments{
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
					fmt.Printf("Document inserted with ID: %s\n", res.InsertedID)
				}
			}
		} else {
			fmt.Printf("No new matching entries found in %s.\n", *feedUrl)
		}

		// Add the new matched IDs to the matched IDs map
		for _, id := range newMatchedIDs {
			matchedIDs[id] = true
		}

		time.Sleep(60 * time.Second) // Wait 60 seconds for the next run
	}
}
