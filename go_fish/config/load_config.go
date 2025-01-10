package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBDatabaseName string
	MongoDBCollection   string
	MongoDBAtlasUri     string
	DiscordWebhookUrl   string
	DiscordUsername     string
	DotLogs             string
}

// LoadConfig loads the configuration from a .env file
func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	return &Config{
		MongoDBDatabaseName: os.Getenv("MONGODB_DATABASE_NAME"),
		MongoDBCollection:   os.Getenv("MONGODB_COLLECTION_NAME"),
		MongoDBAtlasUri:     os.Getenv("MONGODB_ATLAS_URI"),
		DiscordWebhookUrl:   os.Getenv("DISCORD_WEBHOOK_URL"),
		DiscordUsername:     os.Getenv("DISCORD_USERNAME"),
		DotLogs:             os.Getenv("DOT_LOGS"),
	}, nil
}
