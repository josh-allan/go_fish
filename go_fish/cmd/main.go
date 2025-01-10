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
	"github.com/josh-allan/go_fish/util"
)

var lastUpdated *time.Time

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config:", err)
		return
	}

	shared.InitLogs(conf.DotLogs)
	log.Println("Initialising Go Fish")
	log.Println("Logger started")
	log.Println("Config loaded successfully")

	// var store db.Datastore
	ctx := context.Background()

	// store.GetAllDocuments(ctx)
	dbClient := &db.MongoClient{}

	SearchTerms, err := dbClient.GetTerms(ctx, conf.MongodbDatabaseName, "search_terms")
	log.Println("Loading searching terms:", SearchTerms)
	if err != nil {
		log.Fatal("Error retrieving search terms:", err)
		return
	}

	collection := dbClient.GetCollection(conf.MongodbDatabaseName, conf.MongodbCollection)
	existingDocuments, err := dbClient.GetAllDocuments(ctx, conf.MongodbDatabaseName, conf.MongodbCollection)

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
	var webhookURL = conf.DiscordWebhookUrl
	var username = conf.DiscordUsername

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
