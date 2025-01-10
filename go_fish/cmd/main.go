package main

import (
	"context"
	"log"
	"time"

	"github.com/josh-allan/go_fish/config"
	"github.com/josh-allan/go_fish/db"
	"github.com/josh-allan/go_fish/parser"
	"github.com/josh-allan/go_fish/util"
)

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

	ctx := context.Background()
	dbClient, err := db.NewMongoClient(conf.MongodbAtlasUri)
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
		return
	}
	defer func(dbClient *db.MongoClient) {
		err := dbClient.Close()
		if err != nil {
			log.Fatal("Error closing MongoDB client:", err)
		}
	}(dbClient)

	searchTerms, err := dbClient.GetTerms(ctx, conf.MongodbDatabaseName, "search_terms")
	if err != nil {
		log.Fatal("Error retrieving search terms:", err)
		return
	}
	log.Println("Loading search terms:", searchTerms)

	existingDocuments, err := dbClient.GetAllDocuments(ctx, conf.MongodbDatabaseName, conf.MongodbCollection)
	if err != nil {
		log.Fatal("Error retrieving existing entries:", err)
	}

	existingIDs := make(map[string]bool)
	for _, doc := range existingDocuments {
		existingIDs[doc.GUID] = true
	}

	feedUrl := &shared.FeedUrl
	matchedIDs := make(map[string]bool)
	parser.Feed(feedUrl, &searchTerms, nil, matchedIDs)

	for {
		matchingEntries, _, newMatchedIDs, err := parser.Feed(feedUrl, &searchTerms, nil, matchedIDs)
		if err != nil {
			log.Printf("Error searching feed: %v\n", err)
			continue
		}

		if len(matchingEntries) > 0 {
			for _, entry := range matchingEntries {
				if !existingIDs[entry.GUID] {
					matchingDoc := shared.ConvertToMatchingDocuments(entry)
					shared.NotifyDiscord(conf.DiscordWebhookUrl, conf.DiscordUsername, matchingDoc, time.Now())
					_, err := dbClient.InsertDocument(ctx, conf.MongodbDatabaseName, conf.MongodbCollection, entry)
					if err != nil {
						log.Printf("Error inserting document: %v\n", err)
					} else {
						log.Printf("Matching entry found at: %s\n", matchingDoc.GUID)
					}
				}
			}
		} else {
			log.Printf("No new matching entries found in %s.\n", *feedUrl)
		}

		for _, id := range newMatchedIDs {
			matchedIDs[id] = true
		}

		time.Sleep(60 * time.Second)
	}
}
