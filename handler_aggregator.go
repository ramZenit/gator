package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ramZenit/gator/internal/database"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Fatalf("missing argument, syntax: %s <duration> (duration in s|m|h = seconds|minutes|hours eg: 5m is five minutes)", cmd.name)
	}
	durationString := cmd.args[0]
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		log.Fatalf("incorrect <duration> parameter: %s", err)
	}
	fmt.Printf("Collecting feeds every %s\n", durationString)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to update %w", err)
	}
	args := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), args)
	if err != nil {
		return fmt.Errorf("error in marking feed fetched %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), nextFeed.FeedUrl)
	if err != nil {
		return fmt.Errorf("failed fetching RSSfeed %w", err)
	}
	fmt.Printf("**Channel: %s\n", rssFeed.Channel.Title)
	for _, feed := range rssFeed.Channel.Item {
		fmt.Printf("**Feed title: %s\n", feed.Title)
	}
	return nil
}
