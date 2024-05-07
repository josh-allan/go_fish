package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/josh-allan/terraforming-mars/discord"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webhook_URL := os.Getenv("DISCORD_WEBHOOK_URL")
	fmt.Println(webhook_URL)
}
