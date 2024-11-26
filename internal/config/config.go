package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/seanhuebl/blog_aggregator/internal/database"
	"github.com/seanhuebl/blog_aggregator/internal/rss"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	jdata, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}
	filepath := getConfigFilePath()
	err = os.WriteFile(filepath, jdata, 0644)
	if err != nil {
		return fmt.Errorf("error writing: %v", err)
	}
	return nil
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")

	}
	err := f(s, cmd)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

type Command struct {
	Name      string
	Arguments []string
}

type State struct {
	Db        *database.Queries
	ConfigPtr *Config
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {

	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.ConfigPtr.CurrentUserName)
		if err != nil {
			return fmt.Errorf("unable to get user: %v", err)
		}
		return handler(s, cmd, user)
	}
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("login expects one argument: username")
	}
	_, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("user not found")
	}
	err = s.ConfigPtr.SetUser(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fmt.Printf("user: %v has been set\n", s.ConfigPtr.CurrentUserName)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("register expects one argument: username")
	}
	id := uuid.New()
	_, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Arguments[0]})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Printf("user already exists: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("failed to create user: %v\n", err)
		os.Exit(1)
	}
	err = s.ConfigPtr.SetUser(cmd.Arguments[0])

	if err != nil {
		return fmt.Errorf("%s", err)
	}
	var user database.User
	fmt.Printf("user: %v was created\n", s.ConfigPtr.CurrentUserName)
	user, err = s.Db.GetUser(context.Background(), s.ConfigPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Println(user)
	return nil

}

func HandlerGetUsers(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("users takes zero arguments")
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	for _, user := range users {
		if user == s.ConfigPtr.CurrentUserName {
			fmt.Printf("%v (current)\n", user)
		} else {
			fmt.Printf("%v\n", user)
		}
	}
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("reset takes zero arguments")
	}
	err := s.Db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("agg takes one argument")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("error parsing time duration: %v", err)
	}
	println("Collecting feeds every %v", cmd.Arguments[0])
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("addfeed takes exactly two arguments")
	}
	FeedID := uuid.New()
	feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{ID: FeedID, Name: cmd.Arguments[0], Url: cmd.Arguments[1], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("unable to add feed: %v", err)
	}
	followID := uuid.New()
	_, err = s.Db.CreateFeedFollow(
		context.Background(), database.CreateFeedFollowParams{ID: followID, UserID: user.ID, FeedID: feed.ID},
	)
	if err != nil {
		return fmt.Errorf("unable to follow feed: %v", err)
	}

	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("feeds takes zero arguments")
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%v\n", feed)
	}
	return nil
}
func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("follow takes one argument")
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to get feed: %v", err)
	}
	followID := uuid.New()

	_, err = s.Db.CreateFeedFollow(
		context.Background(), database.CreateFeedFollowParams{ID: followID, UserID: user.ID, FeedID: feed.ID},
	)
	if err != nil {
		return fmt.Errorf("unable to create feedfollow: %v", err)
	}
	fmt.Printf("Feed: %v\nUser: %v\n", feed.Name, user.Name)
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("unfollow takes one argument")
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to get feed: %v", err)
	}
	err = s.Db.Unfollow(context.Background(), database.UnfollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("unable to unfollow feed: %v", err)
	}
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("following takes zero arguments")
	}
	feedsFollowed, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get user's feeds: %v", err)
	}
	for _, feed := range feedsFollowed {
		fmt.Println(feed.FeedName)
	}

	return nil
}

func ScrapeFeeds(s *State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch next feed: %v", err)
	}
	s.Db.MarkFeedFetched(context.Background(), nextFeed.ID)
	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("unable to get feed: %v", err)
	}
	postID := uuid.New()
	for _, item := range feed.Channel.Item {
		s.Db.CreatePost(context.Background(), database.CreatePostParams{ID: postID, Title: item.Title, })
	}
	return nil
}

func toNullString(s string) sql.NullString

func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + configFileName
}

func Read() Config {

	filepath := getConfigFilePath()
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("error opening file at location: %v\n", filepath)
	}
	defer file.Close()
	var config Config
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file data")
	}

	if err = json.Unmarshal(data, &config); err != nil {
		fmt.Println("unable to parse data to struct")
	}
	return config
}
