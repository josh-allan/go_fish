package main

import (
	//"fmt"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	//"github.com/josh-allan/terraforming-mars/discord"
	"log"

	"github.com/josh-allan/terraforming-mars/parser"
	//"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// webhook_URL := os.Getenv("DISCORD_WEBHOOK_URL")
	// fmt.Println(webhook_URL)
	interestingSearches := []string{"Nvidia"}
	feedUrl := "https://ozbargain.com.au/feed"
	var lastUpdated *time.Time
	matchedIDs := make(map[string]bool)
	parser.Feed(feedUrl, interestingSearches, lastUpdated, matchedIDs)
	for {
		matchingEntries, lastUpdated, newMatchedIDs, err := parser.Feed(feedUrl, interestingSearches, lastUpdated, matchedIDs)
		if err != nil {
			log.Printf("Error searching feed: %v\n", err)
			continue
		}

		if len(matchingEntries) > 0 {
			for _, entry := range matchingEntries {
				fmt.Printf("Matching entry found in %s: %s at %s\n", feedUrl, entry.Title, time.Now().Format("02/01/2006, 15:04:05"))
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
