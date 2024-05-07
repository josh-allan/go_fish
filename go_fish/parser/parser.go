package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

type MyCustomTranslator struct {
	defaultTranslator *gofeed.DefaultRSSTranslator
}

func NewMyCustomTranslator() *MyCustomTranslator {
	t := &MyCustomTranslator{}
	t.defaultTranslator = &gofeed.DefaultRSSTranslator{}
	return t
}

/*
func (ct *MyCustomTranslator) Translate(feed interface{}) (*gofeed.Feed, error) {
	rss, found := feed.(*rss.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *rss.Feed")
	}

	f, err := ct.defaultTranslator.Translate(rss)
	if err != nil {
		return nil, err
	}
	if rss.ITunesExt != nil && rss.ITunesExt.Author != "" {
		f.Author = rss.ITunesExt.Author
	} else {
		f.Author = rss.ManagingEditor
	}
	return f, nil
}
*/

func main() {
	feedData := "https://ozbargain.com.au/feed"
	fp := gofeed.NewParser()
	fp.RSSTranslator = NewMyCustomTranslator()
	feed, _ := fp.ParseString(feedData)
	fmt.Println(feed.Items)
}
