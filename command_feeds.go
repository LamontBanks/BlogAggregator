package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// Create a feed in the system, attributed to the user
// Fails if the feed already exists
func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Args: feedName, feedUrl
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: %v <Name> <RSS Feed URL>", cmd.name)
	}
	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	// Insert feed info
	addFeedResult, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add: %v (%v) for %v - possible duplicate feed?", feedName, feedUrl, user.Name)
	}
	fmt.Printf("Saved \"%v\" (%v) for user %v\n", addFeedResult.Name, addFeedResult.Url, user.Name)

	// Also follow the added feed
	// Save userId -> feedId mapping
	queryResult, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    addFeedResult.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %v", queryResult.FeedName)
	}
	fmt.Printf("%v followed %v\n", queryResult.UserName, queryResult.FeedName)

	return nil
}

// Lists all feeds from all users
func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%v - added by %v\n", feed.FeedName, feed.UserName.String)
	}

	return nil
}

// TESTING
func handlerUpdateFeed(s *state, cmd command) error {
	// Args: <feed url>
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <feed url>", cmd.name)
	}
	feedUrl := cmd.args[0]

	// Feed info from url
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return fmt.Errorf("feed url %v has not been added yet", feedUrl)
	}
	if err != nil {
		return err
	}

	// Update timestamps
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Fetched %v (%v)\n", feed.Name, feedUrl)

	return nil
}
