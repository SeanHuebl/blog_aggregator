package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to GET feedURL: %v", err)
	}
	req.Header.Add("user-agent", "gator")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("uable to get http response: %v", err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot stream data: %v", err)
	}
	var RSSFeed RSSFeed
	if err := xml.Unmarshal(data, &RSSFeed); err != nil {
		return nil, fmt.Errorf("error unmarshaling xml: %v", err)
	}
	RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)
	RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)
	for _, item := range RSSFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &RSSFeed, nil
}
