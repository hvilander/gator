package main

import (
	"context"
	"fmt"
	"github.com/hvilander/gator/internal/database"
)

func getCurrentUser(s *state) (*database.User, error) {
	username := s.config.CurrentUserName
	user, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		return nil, fmt.Errorf("error getting cur user from db: %w", err)
	}

	return &user, nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := getCurrentUser(s)
		if err != nil {
			return err
		}
		return handler(s, cmd, *user)
	}

}
