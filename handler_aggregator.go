package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	argsFeed := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), argsFeed)
	if err != nil {
		return fmt.Errorf("error in marking feed fetched %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), nextFeed.FeedUrl)
	if err != nil {
		return fmt.Errorf("failed fetching RSSfeed %w", err)
	}

	var argsPost = database.CreatePostParams{}

	fmt.Printf("**Channel: %s\n", rssFeed.Channel.Title)
	for _, feed := range rssFeed.Channel.Item {
		pubDate := convertTime(feed.PubDate)
		argsPost.ID = uuid.New()
		argsPost.CreatedAt = time.Now()
		argsPost.UpdatedAt = time.Now()
		argsPost.Title = feed.Title
		argsPost.Url = feed.Link
		argsPost.Description = feed.Description
		argsPost.PublishedAt = pubDate
		argsPost.FeedID = nextFeed.ID
		_, err := s.db.CreatePost(context.Background(), argsPost)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			} else {
				return fmt.Errorf("failed to create post %w", err)
			}
		}
	}
	return nil
}

func convertTime(input string) time.Time {
	t, err := time.Parse(time.RFC1123Z, input)
	if err == nil {
		return t
	}

	t, err = time.Parse(time.RFC822, input)
	if err == nil {
		return t
	}
	// add more attemps in future (**not really**)
	fmt.Printf("unable to parse time: %v", err)
	return time.Now()
}
