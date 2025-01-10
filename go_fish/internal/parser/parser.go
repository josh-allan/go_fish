package parser

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// containsIgnoreCase performs a case-insensitive check for if a string contains a substring returning true if it does, false otherwise
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Feed takes a feed URL, a list of search terms, the last time the feed was updated and a map of matched IDs.
//
// Returns a list of matching entries, the time of the last entry, and the IDs of the new matching entries
func Feed(feedUrl *string, searches *[]string, lastUpdated *time.Time, matchedIDs map[string]bool) ([]*gofeed.Item, *time.Time, []string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(*feedUrl)
	if err != nil {
		return nil, nil, nil, err
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
		for _, searchTerm := range *searches {
			if searchTerm != "" && (containsIgnoreCase(entry.Title, searchTerm) || containsIgnoreCase(entry.Description, searchTerm)) {
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
