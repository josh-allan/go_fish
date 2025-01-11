package delete_term

import (
	"context"
	"github.com/josh-allan/go_fish/internal/config"
	"github.com/josh-allan/go_fish/internal/db"
)

func DeleteTerm(config *config.Config, term string) error {
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

	_, err = dbClient.DeleteDocument(ctx, config.MongoDBDatabaseName, config.SearchTerms, map[string]string{"term": term})
	if err != nil {
		return err
	}

	return nil
}
