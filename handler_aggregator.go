package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching RSSfeed %w", err)
	}
	fmt.Printf("%+v\n", rss)
	return nil
}
