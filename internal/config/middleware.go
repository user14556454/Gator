package config

import (
	"context"
	"fmt"

	"github.com/user14556454/project2/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.SharedConfig.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error fetching the user: %w", err)
		}
		if err := handler(s, cmd, user); err != nil {
			return fmt.Errorf("Error in the handler: %w", err)
		}
		return nil
	}
}
