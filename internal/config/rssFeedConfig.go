package config

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/user14556454/project2/internal/database"
)

type RSSFeed struct {
	Channel ChannelRSS `xml:"channel"`
}

type ChannelRSS struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Items       []ItemsRSS `xml:"item"`
}

type ItemsRSS struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error during request making: %w", err)
	}
	newClient := http.Client{}
	res, err := newClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error during response receive: %w", err)
	}
	defer res.Body.Close()

	resSlice, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error during ioReadAll: %w", err)
	}
	var jsonData RSSFeed
	if err := xml.Unmarshal(resSlice, &jsonData); err != nil {
		return nil, fmt.Errorf("Error during unmarshal: %w", err)
	}

	jsonData.Channel.Title = html.UnescapeString(jsonData.Channel.Title)
	jsonData.Channel.Description = html.UnescapeString(jsonData.Channel.Description)
	for i := range jsonData.Channel.Items {
		jsonData.Channel.Items[i].Title = html.UnescapeString(jsonData.Channel.Items[i].Title)
		jsonData.Channel.Items[i].Description = html.UnescapeString(jsonData.Channel.Items[i].Description)
	}
	return &jsonData, nil
}

func scrapeFeeds(s *State) error {
	ctx := context.Background()
	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	if err := s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:        feed.UserID,
	}); err != nil {
		return err
	}

	rssFeed, err := FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}
	for i := range rssFeed.Channel.Items {
		fmt.Println(rssFeed.Channel.Items[i].Title)
	}
	return nil
}
