package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MONGODB_DATABASE_NAME string
	MONGODB_COLLECTION    string
	MONGODB_ATLAS_URI     string
	DISCORD_WEBHOOK_URL   string
	DISCORD_USERNAME      string
	DOT_LOGS              string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	return &Config{
		MONGODB_DATABASE_NAME: os.Getenv("MONGODB_DATABASE_NAME"),
		MONGODB_COLLECTION:    os.Getenv("MONGODB_COLLECTION_NAME"),
		MONGODB_ATLAS_URI:     os.Getenv("MONGODB_ATLAS_URI"),
		DISCORD_WEBHOOK_URL:   os.Getenv("DISCORD_WEBHOOK_URL"),
		DISCORD_USERNAME:      os.Getenv("DISCORD_USERNAME"),
		DOT_LOGS:              os.Getenv("DOT_LOGS"),
	}, nil
}
