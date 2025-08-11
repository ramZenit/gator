package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ramZenit/gator/internal/database"
)

func handlerCreateFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("missing argument, syntax: follow <url>")
	}
	feedURL := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("unable to retrieve feed info: %w", err)
	}
	paramsFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), paramsFeedFollow)
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %w", err)
	}
	fmt.Printf("%s just followed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowsPerUser(s *state, cmd command, user database.User) error {
	feedList, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("unable to retrieve feed follows: %w", err)
	}
	if len(feedList) == 0 {
		fmt.Println("no feeds followed")
		return nil
	}
	for _, feed := range feedList {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerUnfollowPerUser(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		log.Fatalf("missing argument, syntax: %s <URL>", cmd.name)
		return errors.New("missing argument, syntax: follow <url>")
	}
	feedURL := cmd.args[0]

	args := database.DeleteFeedFollowParams{
		Name: user.Name,
		Url:  feedURL,
	}
	err := s.db.DeleteFeedFollow(context.Background(), args)
	if err != nil {
		return fmt.Errorf("unable to unfollow feed: %w", err)
	}
	fmt.Println("Feed unfollowed")
	return nil
}
