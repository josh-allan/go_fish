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
	"os"
)

type MatchingDocuments struct {
	Name string `bson:"name"`
	Time string `bson:"time"`
	Url  string `bson:"url"`
}

func main() {
	err := godotenv.Load()

	mongodb_database := os.Getenv("MONGODB_DB_NAME")
	mongodb_collection := os.Getenv("MONGODB_COLLECTION_NAME")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	collection := mongo.GetCollection(mongodb_database, mongodb_collection)

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
				matchingDocuments := &MatchingDocuments{Name: entry.Title,
					Time: time.Now().Format("02/01/2006, 15:04:05"),
					Url:  feedUrl,
				}
				res, err := collection.InsertOne(ctx, matchingDocuments)
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
