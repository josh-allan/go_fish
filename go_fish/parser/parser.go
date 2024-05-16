package parser

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
func Feed(feedUrl *string, interestingSearches *[]string, lastUpdated *time.Time, matchedIDs map[string]bool) ([]*gofeed.Item, *time.Time, []string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(*feedUrl)
	if err != nil {
		panic(err)
	}
	var matchingEntries []*gofeed.Item
	var newMatchedIDs []string

	for _, entry := range feed.Items {
		publishedTime := entry.PublishedParsed
		// Skip over entries that were published before the last entry that was processed
		if lastUpdated != nil && publishedTime != nil && publishedTime.Before(*lastUpdated) {
			continue
		}

		// We also want to skip over entries that have already been matched
		if matchedIDs != nil && matchedIDs[entry.GUID] {
			continue
		}

		// If we match on a new string, let's append it to the appropriate list
		// ensuring to only alert on a new entry
		for _, term := range *interestingSearches {
			if term != "" && (containsIgnoreCase(entry.Title, term) || containsIgnoreCase(entry.Description, term)) {
				matchingEntries = append(matchingEntries, entry)
				newMatchedIDs = append(newMatchedIDs, entry.GUID)

				break
			}
		}
	}

	// Return the matching entries, the published time of the last entry,
	// and the IDs of the new matching entries
	return matchingEntries, feed.Items[len(feed.Items)-1].PublishedParsed, newMatchedIDs, nil
}
