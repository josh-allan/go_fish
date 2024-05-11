package parser

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

func Feed() {
	url := "https://ozbargain.com.au/feed"
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		panic(err)
	}
	// for entry
	return feed

}
