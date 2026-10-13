package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/user14556454/project2/internal/database"
)

type State struct {
	Db           *database.Queries
	SharedConfig *Config
}

type Command struct {
	Name      string
	Arguments []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Username is also required")
	}

	ctx := context.Background()
	_, err := s.Db.GetUser(ctx, cmd.Arguments[1])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("No user exists with username %v", cmd.Arguments[1])
		}
		return fmt.Errorf("Error executing the query: %w", err)
	}

	s.SharedConfig.CurrentUserName = cmd.Arguments[1]
	s.SharedConfig.SetUser()
	fmt.Printf("User %v has been set\n", cmd.Arguments[1])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("No username provided")
	}

	ctx := context.Background()
	_, err := s.Db.GetUser(ctx, cmd.Arguments[1])
	if err == nil {
		return fmt.Errorf("User already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Error executing the query: %w", err)
	}

	newUser, err := s.Db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[1],
	})
	if err != nil {
		return fmt.Errorf("Error creating new user: %w", err)
	}
	s.SharedConfig.CurrentUserName = cmd.Arguments[1]
	s.SharedConfig.SetUser()
	fmt.Println("User created successfully")
	fmt.Println("ID: ", newUser.ID)
	fmt.Println("CreatedAt: ", newUser.CreatedAt)
	fmt.Println("UpdatedAt: ", newUser.UpdatedAt)
	fmt.Println("Name: ", newUser.Name)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if err := s.Db.DeleteFeedFollows(context.Background()); err != nil {
		return fmt.Errorf("Error deleting feed follows")
	}
	if err := s.Db.DeleteFeeds(context.Background()); err != nil {
		return fmt.Errorf("Error deleting feeds")
	}
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("Error deleting the users")
	}
	s.SharedConfig.CurrentUserName = ""
	s.SharedConfig.SetUser()
	return nil
}

func HandlerGetAllUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("No users are registered")
		}
		return fmt.Errorf("Error executing the query: %w", err)
	}
	for _, user := range users {
		if user.Name == s.SharedConfig.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}
	return nil
}

func Handleragg(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Not enough arguments")
	}
	fmt.Println("Collecting feeds every", cmd.Arguments[1])
	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[1])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func Handleraddfeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 3 {
		return fmt.Errorf("Not enough arguments")
	}
	ctx := context.Background()
	feedID := uuid.New()
	userFeed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[1],
		Url:       cmd.Arguments[2],
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	if _, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	}); err != nil {
		return fmt.Errorf("Error during feed follow: %w", err)
	}

	fmt.Println(userFeed)
	return nil
}

func Handlerfeeds(s *State, cmd Command) error {
	ctx := context.Background()
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %w", err)
	}
	for i := range feeds {
		userName, err := s.Db.GetUserID(ctx, feeds[i].UserID)
		if err != nil {
			return fmt.Errorf("Error fetching user with id %v: %w", feeds[i].UserID, err)
		}
		fmt.Println(feeds[i].Name)
		fmt.Println(feeds[i].Url)
		fmt.Println(userName)
	}
	return nil
}

func Handlerfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Not enough arguments")
	}
	ctx := context.Background()
	feedURL := cmd.Arguments[1]
	feed, err := s.Db.GetFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("Error during feed load: %w", err)
	}

	if _, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}); err != nil {
		return fmt.Errorf("Error during feed_follow query: %w", err)
	}

	fmt.Println(feed.Name)
	fmt.Println(user.Name)

	return nil
}

func Handlerfollowing(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	feedFollowRows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("Error during following load: %w", err)
	}
	for i := range feedFollowRows {
		fmt.Println(feedFollowRows[i].Name)
	}
	return nil
}

func Handlerunfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Not enough arguments")
	}
	ctx := context.Background()
	feed, err := s.Db.GetFeed(ctx, cmd.Arguments[1])
	if err != nil {
		return fmt.Errorf("Error loading feed: %w", err)
	}
	if err := s.Db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("Error deleting the feed: %w", err)
	}
	return nil
}

type Commands struct {
	CommandsMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	for cmdName := range c.CommandsMap {
		if cmdName == cmd.Arguments[0] {
			if err := c.CommandsMap[cmdName](s, cmd); err != nil {
				return fmt.Errorf("Error during execution of command: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Command not found")
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandsMap[name] = f
}
