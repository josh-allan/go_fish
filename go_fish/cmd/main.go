package main

import (
	"github.com/josh-allan/go_fish/internal/config"
	addterm "github.com/josh-allan/go_fish/internal/tasks/add-term"
	listterm "github.com/josh-allan/go_fish/internal/tasks/list-terms"
	"github.com/josh-allan/go_fish/internal/tasks/scraper"
	"log"

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

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert RSS feed search terms",
	Long:  "Insert new terms for the RSS scraper to search on",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		term, _ := cmd.Flags().GetString("term")
		err = addterm.AddTerm(config, term)
		if err != nil {
			log.Fatalf("Error adding term: %v", err)
		}
		log.Printf("Added term: %s", term)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List RSS feed search terms",
	Long:  "List the current search terms",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		listterm.ListTerms(config)
	},
}

func init() {
	insertCmd.Flags().StringP("term", "t", "", "The term to add to the search list")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(scrapeCmd)
	rootCmd.AddCommand(insertCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
