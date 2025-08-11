package main

import (
	"context"
	"encoding/xml"
	"errors"
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

var ErrNotFound = errors.New("not found")

const ErrCodeUniqueViolation = "23505"

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating the request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %w", err)
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading the body: %w", err)
		}

		RSS := &RSSFeed{}
		err = xml.Unmarshal(data, RSS)
		if err != nil {
			return nil, fmt.Errorf("error decoding the body: %w", err)
		}
		RSS.unescapeAll()
		return RSS, nil

	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}

func (rss *RSSFeed) unescapeAll() {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}
}
