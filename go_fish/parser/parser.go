package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

func feed() {
	url := "https://ozbargain.com.au/feed"
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		panic(err)
	}
	// for entry
	fmt.Println(feed.PublishedParsed)
}

func main() {
	feed()
}
