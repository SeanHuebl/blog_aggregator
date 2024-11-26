package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

// RSSFeed represents the structure of an RSS feed parsed from XML.
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`       // The title of the RSS feed
		Link        string    `xml:"link"`        // The URL link of the RSS feed
		Description string    `xml:"description"` // A brief description of the RSS feed
		Item        []RSSItem `xml:"item"`        // A list of items (posts) in the RSS feed
	} `xml:"channel"`
}

// RSSItem represents an individual item (post) in an RSS feed.
type RSSItem struct {
	Title       string `xml:"title"`       // The title of the RSS item
	Link        string `xml:"link"`        // The URL link to the RSS item
	Description string `xml:"description"` // A brief description of the RSS item
	PubDate     string `xml:"pubDate"`     // The publication date of the RSS item
}

// FetchFeed retrieves and parses an RSS feed from the provided URL.
//
// Parameters:
// - ctx: A context for managing request cancellation and timeouts.
// - feedURL: The URL of the RSS feed to fetch.
//
// Returns:
// - A pointer to the RSSFeed struct containing the parsed feed data.
// - An error if the feed cannot be retrieved or parsed.
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new HTTP GET request with the provided context
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to GET feedURL: %v", err)
	}

	// Add a custom User-Agent header
	req.Header.Add("user-agent", "gator")

	// Use the default HTTP client to send the request
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get HTTP response: %v", err)
	}

	// Read the response body into memory
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot stream data: %v", err)
	}

	// Parse the XML data into an RSSFeed struct
	var RSSFeed RSSFeed
	if err := xml.Unmarshal(data, &RSSFeed); err != nil {
		return nil, fmt.Errorf("error unmarshaling XML: %v", err)
	}

	// Unescape HTML entities in the RSS feed's title and description
	RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)
	RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)

	// Unescape HTML entities in each RSS item's title and description
	for i := range RSSFeed.Channel.Item {
		RSSFeed.Channel.Item[i].Title = html.UnescapeString(RSSFeed.Channel.Item[i].Title)
		RSSFeed.Channel.Item[i].Description = html.UnescapeString(RSSFeed.Channel.Item[i].Description)
	}

	// Return the parsed RSS feed
	return &RSSFeed, nil
}
