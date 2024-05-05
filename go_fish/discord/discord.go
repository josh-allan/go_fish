package discord

import (
	"log"
	"os"

	"github.com/ecnepsnai/discord"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webhook_URL := os.Getenv("DISCORD_WEBHOOK_URL")

	discord.WebhookURL = webhook_URL
	discord.Say("testing webhook")
}
