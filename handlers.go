package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

func getNullTimeFromStr(s string) sql.NullTime {
	time, error := time.Parse(time.RFC1123, s)
	var isValid = error != nil

	return sql.NullTime{
		Time:  time,
		Valid: isValid,
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
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg takes one argument, time between reqs")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s...\n", duration.String())

	ticker := time.NewTicker(duration)

	// starts immediately and tickes every time the ticker hits the channel
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println(fmt.Errorf("scrape error: %w", err))
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// check args need name and url
	if len(cmd.args) != 2 {
		return fmt.Errorf("add feed needs two arguments, a name an url")
	}

	name := getNullString(cmd.args[0])
	url := getNullString(cmd.args[1])

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
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	fmt.Println("userID")
	fmt.Println(userID)

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentNullTime,
		UpdatedAt: currentNullTime,
		UserID:    userID,
		FeedID:    getNullUUID(feed.ID),
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("added feed follow: %s for user: %s\n", feedFollow.FeedName.String, feedFollow.UserName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retreving feeds: %w", err)
	}

	fmt.Println("Feed Name                   URL                  User")
	fmt.Println("-----------------------------------------------------")
	for _, feed := range feeds {

		fmt.Printf("%s           %s            %s\n", feed.FeedName.String, feed.Url.String, feed.UserName)
		fmt.Println("-----------------------------------------------------")
	}
	return nil

}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow expects a single argument, a url")
	}
	urlForQuery := getNullString(cmd.args[0])

	feed, err := s.db.GetFeedByURL(context.Background(), urlForQuery)
	if err != nil {
		return fmt.Errorf("error finding feed by URL: %w", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: getNullTime(),
		UpdatedAt: getNullTime(),
		UserID:    getNullUUID(user.ID),
		FeedID:    getNullUUID(feed.ID),
	}

	result, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	fmt.Printf("FeedName: %s, User: %s", result.FeedName.String, result.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	follows, err := s.db.GetFollowsByUserID(context.Background(), getNullUUID(user.ID))
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, follow := range follows {
		fmt.Println(follow.FeedName)

	}

	return nil

}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow expects only one argument")
	}

	queryArgs := database.DeleteFeedFollowByUserAndURLParams{
		UserID: getNullUUID(user.ID),
		Url:    getNullString(cmd.args[0]),
	}

	return s.db.DeleteFeedFollowByUserAndURL(context.Background(), queryArgs)
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next to fetch: %w", err)
	}

	params := database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: getNullTime(),
	}

	s.db.MarkFeedFetched(context.Background(), params)

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url.String)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("Scraping %s\n", rssFeed.Channel.Title)
	for _, item := range rssFeed.Channel.Item {
		currentNullTime := getNullTime()

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   currentNullTime,
			UpdatedAt:   currentNullTime,
			Title:       getNullString(item.Title),
			Url:         getNullString(item.Link),
			Description: getNullString(item.Description),
			PublishedAt: getNullTimeFromStr(item.PubDate), // todo this may need some normalization
			FeedID:      getNullUUID(nextFeed.ID),
		}

		post, err := s.db.CreatePost(context.Background(), postParams)
		// error are common since it will try to create dups
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				// ignore the error
				fmt.Println("skipping, post already exists")
				continue
			}
			return fmt.Errorf("error creating post: %w", err)
		}

		fmt.Printf("created post: %s\n", post.Title.String)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	limit = 2
	if len(cmd.args) == 1 {
		arglimit, err := strconv.Atoi(cmd.args[0])
		limit = int32(arglimit)
		if err != nil {
			return fmt.Errorf("error parsing your int argument: %w", err)
		}
	}

	params := database.GetPostsByUserIDParams{
		UserID: getNullUUID(user.ID),
		Limit:  limit,
	}

	posts, err := s.db.GetPostsByUserID(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error getting posts: %w", err)
	}

	for idx, post := range posts {
		fmt.Printf("Post #%d *  %s\n", idx+1, post.Title.String)
		fmt.Printf("   pubDate: %s\n", post.PublishedAt.Time)
		fmt.Printf("   * %s\n", post.Description.String)
	}

	return nil

}
