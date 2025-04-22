package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hvilander/gator/internal/database"
)

// get an sql friendly null time thingy
func getNullTime() sql.NullTime {
	return sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
}

func getNullUUID(id uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		UUID:  id,
		Valid: true,
	}
}

func getNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login expects a single argument")
	}

	if cmd.args[0] == "" {
		return fmt.Errorf("empty string user name is bad")
	}

	username := cmd.args[0]

	user, err := s.db.GetUserByName(context.Background(), username)
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

	currentNullTime := getNullTime()

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

func handlerReset(s *state, cmd command) error {
	fmt.Println("deleteing all users")
	return s.db.ResetUsers(context.Background())
}

// list all users
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUserName := s.config.CurrentUserName

	for _, user := range users {
		fmt.Printf(" * %s", user.Name)
		if currentUserName == user.Name {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Printf("\n")
		}
	}
	return nil

}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to fetch feed: %w", err))
	}

	fmt.Println(rssFeed)

	return nil

}

func handlerAddFeed(s *state, cmd command) error {
	// check args need name and url
	if len(cmd.args) != 2 {
		return fmt.Errorf("add feed needs to arguments, a name an url")
	}

	name := getNullString(cmd.args[0])
	url := getNullString(cmd.args[1])

	//get current user from db
	username := s.config.CurrentUserName
	fmt.Println(username)
	user, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting cur user from db: %w", err)
	}

	currentNullTime := getNullTime()
	userID := getNullUUID(user.ID)

	// create the feed record
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentNullTime,
		UpdatedAt: currentNullTime,
		Name:      name,
		Url:       url,
		UserID:    userID,
	}
	s.db.CreateFeed(context.Background(), feedParams)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retreving feeds: %w", err)
	}

	for _, feed := range feeds {

		fmt.Println("Feed Name         URL        User")
		fmt.Println("---------------------------------")
		fmt.Printf("%s           %s            %s", feed.FeedName.String, feed.Url.String, feed.UserName)
	}
	return nil

}
