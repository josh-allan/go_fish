package main

import (
	"log"

	"github.com/josh-allan/go_fish/scraper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go_fish",
	Short: "Go Fish searches RSS feeds",
	Long:  "Go Fish is a tool to search RSS feeds for specific terms",
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape RSS feeds for search terms",
	Long:  "Initialise the main scraper loop",
	Run: func(cmd *cobra.Command, args []string) {
		scraper.GoFish()
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

}
