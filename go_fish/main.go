package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/josh-allan/go_fish/db"
	//"github.com/josh-allan/go_fish/discord"
	"github.com/josh-allan/go_fish/parser"
	"github.com/josh-allan/go_fish/util"
)

var lastUpdated *time.Time

func main() {
	err := godotenv.Load()

	mongodb_database := os.Getenv("MONGODB_DB_NAME")
	mongodb_collection := os.Getenv("MONGODB_COLLECTION_NAME")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db_client := &db.MongoClient{}
	collection := db_client.GetCollection(mongodb_database, mongodb_collection)
	checkDb := db_client.GetAllDocuments(mongodb_database, mongodb_collection)

	checkDb()
	defer cancel()

	interestingSearches := &shared.SearchTerms
	feedUrl := &shared.FeedUrl
	matchedIDs := make(map[string]bool)
	parser.Feed(feedUrl, interestingSearches, lastUpdated, matchedIDs)

	for {
		matchingEntries, _, newMatchedIDs, err := parser.Feed(feedUrl, interestingSearches, nil, matchedIDs)
		if err != nil {
			log.Printf("Error searching feed: %v\n", err)
			continue
		}

		if len(matchingEntries) > 0 {
			for _, entry := range matchingEntries {

				fmt.Printf("Matching entry found in %s: %s at %s\n", entry.Link, entry.Title, time.Now().Format("02/01/2006, 15:04:05"))
				matchingDocuments := &shared.MatchingDocuments{
					ID:            primitive.NewObjectID(),
					Name:          entry.Title,
					PublishedTime: primitive.NewDateTimeFromTime(time.Time(*entry.PublishedParsed)),
					Url:           entry.Link,
				}
				res, err := collection.InsertOne(ctx, matchingDocuments)
				if err != nil {
					log.Fatalf("Error inserting document: %v", err)
				}
				fmt.Printf("Document inserted with ID: %s\n", res.InsertedID)

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
