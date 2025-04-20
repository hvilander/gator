package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hvilander/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login expects a single argument")
	}

	if cmd.args[0] == "" {
		return fmt.Errorf("empty string user name is bad")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("current user set: %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if (len(cmd.args)) != 1 {
		return fmt.Errorf("register expects a single argument")
	}

	name := cmd.args[0]
	if name == "" {
		return fmt.Errorf("empty name argment bad")
	}

	currentNullTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentNullTime,
		UpdatedAt: currentNullTime,
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	s.config.SetUser(user.Name)
	fmt.Println("User Created: ")
	fmt.Println(user)

	return nil
}
