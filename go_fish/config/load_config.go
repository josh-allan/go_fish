package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongodbDatabaseName string
	MongodbCollection   string
	MongodbAtlasUri     string
	DiscordWebhookUrl   string
	DiscordUsername     string
	DotLogs             string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	return &Config{
		MongodbDatabaseName: os.Getenv("MONGODB_DATABASE_NAME"),
		MongodbCollection:   os.Getenv("MONGODB_COLLECTION_NAME"),
		MongodbAtlasUri:     os.Getenv("MONGODB_ATLAS_URI"),
		DiscordWebhookUrl:   os.Getenv("DISCORD_WEBHOOK_URL"),
		DiscordUsername:     os.Getenv("DISCORD_USERNAME"),
		DotLogs:             os.Getenv("DOT_LOGS"),
	}, nil
}
