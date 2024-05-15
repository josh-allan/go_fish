package main

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"

	//"github.com/josh-allan/terraforming-mars/discord"
	"log"

	"github.com/josh-allan/terraforming-mars/db"

	"github.com/josh-allan/terraforming-mars/parser"
	//"os"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	collection := mongo.GetCollection("MONGODB_DB_NAME", "MONGODB_COLLECTION_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	interestingSearches := []string{"Samsung", "Steam"}
	feedUrl := "https://ozbargain.com.au/feed"
	var lastUpdated *time.Time
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
				fmt.Printf("Matching entry found in %s: %s at %s\n", feedUrl, entry.Title, time.Now().Format("02/01/2006, 15:04:05"))
				res, err := collection.InsertOne(ctx, bson.M{"Title": entry.Title})
				if err != nil {
					log.Fatalf("Error inserting document: %v", err)
				}
				fmt.Printf("Document inserted with ID: %s\n", res.InsertedID)

			}
		} else {
			fmt.Printf("No new matching entries found in %s.\n", feedUrl)
		}

		// Add the new matched IDs to the matched IDs map
		for _, id := range newMatchedIDs {
			matchedIDs[id] = true
		}

		time.Sleep(60 * time.Second) // Wait 60 seconds for the next run
	}
}
