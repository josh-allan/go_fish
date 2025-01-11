package list_terms

import (
	"context"
	"fmt"
	"github.com/josh-allan/go_fish/internal/config"
	"github.com/josh-allan/go_fish/internal/db"
)

// List terms from database
func ListTerms(config *config.Config) {

	ctx := context.Background()
	dbClient, err := db.NewMongoClient(config.MongoDBAtlasUri)
	if err != nil {
		fmt.Println("Error creating MongoDB client:", err)
		return
	}
	defer func(dbClient *db.MongoClient) {
		err := dbClient.Close()
		if err != nil {
			fmt.Println("Error closing MongoDB client:", err)
		}
	}(dbClient)

	terms, err := dbClient.GetTerms(ctx, config.MongoDBDatabaseName, config.SearchTerms)
	if err != nil {
		fmt.Println("Error retrieving search terms:", err)
		return
	}

	for _, term := range terms {
		fmt.Println(term)
	}
}
