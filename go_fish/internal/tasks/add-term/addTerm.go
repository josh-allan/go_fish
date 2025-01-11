package add_term

import (
	"context"
	"github.com/josh-allan/go_fish/internal/config"
	"github.com/josh-allan/go_fish/internal/db"
)

// Add term to database
func AddTerm(config *config.Config, term string) error {

	ctx := context.Background()
	dbClient, err := db.NewMongoClient(config.MongoDBAtlasUri)
	if err != nil {
		return err
	}
	defer func(dbClient *db.MongoClient) {
		err := dbClient.Close()
		if err != nil {
			return
		}
	}(dbClient)

	_, err = dbClient.InsertDocument(ctx, config.MongoDBDatabaseName, config.SearchTerms, map[string]string{"term": term})
	if err != nil {
		return err
	}

	return nil
}
