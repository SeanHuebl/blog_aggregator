package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/seanhuebl/blog_aggregator/internal/database"
	"github.com/seanhuebl/blog_aggregator/internal/rss"
)

// configFileName defines the location of the configuration file within the user's home directory.
const configFileName = "/.gatorconfig.json"

// Config represents the application's configuration, including database URL and the current user.
type Config struct {
	DbUrl           string `json:"db_url"`            // Database connection URL
	CurrentUserName string `json:"current_user_name"` // The currently logged-in user's username
}

// SetUser updates the current user in the configuration file.
//
// Parameters:
// - username: The new username to set.
//
// Returns:
// - An error if the configuration cannot be updated or saved.
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	// Marshal the configuration struct to a formatted JSON string.
	jdata, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}
	// Write the updated configuration to the file.
	filepath := getConfigFilePath()
	err = os.WriteFile(filepath, jdata, 0644)
	if err != nil {
		return fmt.Errorf("error writing: %v", err)
	}
	return nil
}

// Commands stores a registry of available commands and their associated handlers.
type Commands struct {
	Commands map[string]func(*State, Command) error // Maps command names to their handler functions
}

// Register adds a new command to the registry.
//
// Parameters:
// - name: The name of the command.
// - f: The function to handle the command.
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

// Run executes the specified command.
//
// Parameters:
// - s: The current application state.
// - cmd: The command to execute.
//
// Returns:
// - An error if the command cannot be found or its handler fails.
func (c *Commands) Run(s *State, cmd Command) error {
	// Check if the command exists in the registry.
	f, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	// Execute the command's handler.
	err := f(s, cmd)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// Command represents a user-specified command and its arguments.
type Command struct {
	Name      string   // The name of the command
	Arguments []string // A list of arguments passed to the command
}

// State holds the application state, including the database and configuration.
type State struct {
	Db        *database.Queries // A pointer to the database queries interface
	ConfigPtr *Config           // A pointer to the application's configuration
}

// MiddlewareLoggedIn ensures that a user is logged in before executing a command.
//
// Parameters:
// - handler: The function to execute if the user is authenticated.
//
// Returns:
// - A wrapper function that first validates the logged-in user, then executes the handler.
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		// Retrieve the currently logged-in user from the database.
		user, err := s.Db.GetUser(context.Background(), s.ConfigPtr.CurrentUserName)
		if err != nil {
			return fmt.Errorf("unable to get user: %v", err)
		}
		// Pass control to the original handler with the validated user.
		return handler(s, cmd, user)
	}
}

// HandlerLogin validates the provided username and sets it as the current user.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing the username as an argument.
//
// Returns:
// - An error if the username is invalid or cannot be set.
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

// HandlerRegister creates a new user and sets them as the current user.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing the new username as an argument.
//
// Returns:
// - An error if the username already exists or cannot be created.
func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("register expects one argument: username")
	}
	id := uuid.New()
	_, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Arguments[0],
	})
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
	user, err := s.Db.GetUser(context.Background(), s.ConfigPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Printf("user: %v was created\n", s.ConfigPtr.CurrentUserName)
	fmt.Println(user)
	return nil
}

// HandlerGetUsers retrieves and prints all users in the database.
//
// Parameters:
// - s: The current application state.
// - cmd: The command with no arguments.
//
// Returns:
// - An error if users cannot be retrieved.
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

// HandlerReset deletes all users in the database.
//
// Parameters:
// - s: The current application state.
// - cmd: The command with no arguments.
//
// Returns:
// - An error if the reset operation fails.
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

// HandlerAgg periodically scrapes feeds based on a provided interval.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing the time interval as an argument.
//
// Returns:
// - An error if the interval is invalid or the scraping fails.
func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("agg takes one argument")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("error parsing time duration: %v", err)
	}
	fmt.Printf("Collecting feeds every %v\n", cmd.Arguments[0])
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

// HandlerAddFeed adds a new feed and subscribes the current user to it.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing the feed name and URL as arguments.
// - user: The currently logged-in user.
//
// Returns:
// - An error if the feed cannot be added or followed.
func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("addfeed takes exactly two arguments")
	}
	feedID := uuid.New()
	feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID: feedID, Name: cmd.Arguments[0], Url: cmd.Arguments[1], UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to add feed: %v", err)
	}
	followID := uuid.New()
	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: followID, UserID: user.ID, FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to follow feed: %v", err)
	}
	fmt.Println(feed)
	return nil
}

// HandlerFeeds retrieves and prints all feeds in the database.
//
// Parameters:
// - s: The current application state.
// - cmd: The command with no arguments.
//
// Returns:
// - An error if feeds cannot be retrieved.
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

// HandlerFollow subscribes the current user to an existing feed by its URL.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing the feed URL as an argument.
// - user: The currently logged-in user.
//
// Returns:
// - An error if the feed cannot be found or the follow operation fails.
func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("follow takes one argument")
	}
	// Retrieve the feed using the provided URL.
	feed, err := s.Db.GetFeed(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to get feed: %v", err)
	}
	// Create a new follow record for the user.
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

// HandlerUnfollow unsubscribes the current user from a feed by its URL.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing the feed URL as an argument.
// - user: The currently logged-in user.
//
// Returns:
// - An error if the feed cannot be found or the unfollow operation fails.
func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("unfollow takes one argument")
	}
	// Retrieve the feed using the provided URL.
	feed, err := s.Db.GetFeed(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("unable to get feed: %v", err)
	}
	// Delete the follow record for the user and feed.
	err = s.Db.Unfollow(context.Background(), database.UnfollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("unable to unfollow feed: %v", err)
	}
	return nil
}

// HandlerFollowing lists all feeds that the current user is following.
//
// Parameters:
// - s: The current application state.
// - cmd: The command with no arguments.
// - user: The currently logged-in user.
//
// Returns:
// - An error if the feeds cannot be retrieved.
func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("following takes zero arguments")
	}
	// Retrieve the list of feeds the user is following.
	feedsFollowed, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get user's feeds: %v", err)
	}
	// Print each feed's name.
	for _, feed := range feedsFollowed {
		fmt.Println(feed.FeedName)
	}
	return nil
}

// HandlerBrowse retrieves and displays posts from feeds that the current user follows.
//
// Parameters:
// - s: The current application state.
// - cmd: The command containing an optional limit argument.
// - user: The currently logged-in user.
//
// Returns:
// - An error if posts cannot be retrieved or if the limit argument is invalid.
func HandlerBrowse(s *State, cmd Command, user database.User) error {
	var limit int
	var err error

	// Validate the number of arguments and parse the limit if provided.
	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("browse takes up to one argument")
	} else if len(cmd.Arguments) == 1 {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("post-browse argument must be an integer: %v", err)
		}
	} else {
		limit = 2 // Default limit if no argument is provided.
	}

	// Retrieve posts from the user's followed feeds with the specified limit.
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID, Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts: %v", err)
	}

	// Display the retrieved posts.
	for _, post := range posts {
		fmt.Printf("Title:\n%v\n\nURL:\n%v\n\n", post.Title.String, post.Url.String)
		fmt.Printf("Content:\n%v\n\n", post.Description.String)
		fmt.Printf("Published on:\n%v\n\n", post.PublishedAt.Time)
	}
	return nil
}

// ScrapeFeeds fetches the next feed to be processed and stores its posts in the database.
//
// Parameters:
// - s: The current application state.
//
// Returns:
// - An error if the feed cannot be fetched or posts cannot be stored.
func ScrapeFeeds(s *State) error {
	// Retrieve the next feed to fetch, prioritized by last fetched time.
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch next feed: %v", err)
	}

	// Mark the feed as fetched with the current timestamp.
	s.Db.MarkFeedFetched(context.Background(), nextFeed.ID)

	// Fetch the RSS feed from the given URL.
	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("unable to get feed: %v", err)
	}

	// Process each item in the feed and store it as a post in the database.
	for _, item := range feed.Channel.Item {
		postID := uuid.New() // Generate a unique ID for the post.
		err := s.Db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          postID,
				Title:       parseToNullString(item.Title),
				Url:         parseToNullString(item.Link),
				Description: parseToNullString(item.Description),
				PublishedAt: parseToNullTime(item.PubDate),
				FeedID:      nextFeed.ID,
			},
		)
		if err != nil {
			return fmt.Errorf("unable to create post: %v", err)
		}
	}
	return nil
}

// parseToNullTime converts a date string into a sql.NullTime value.
//
// Parameters:
// - date: The date string to parse.
//
// Returns:
// - A sql.NullTime value, which is valid if the date string matches any known format.
func parseToNullTime(date string) sql.NullTime {
	// List of possible date formats for parsing.
	var dateFormats = []string{
		time.RFC1123,          // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,         // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05", // ISO 8601 without time zone
		"2006-01-02",          // Date only
	}

	// Attempt to parse the date string using each format.
	for _, format := range dateFormats {
		parsedTime, err := time.Parse(format, date)
		if err == nil {
			return sql.NullTime{Time: parsedTime, Valid: true}
		}
	}
	return sql.NullTime{Valid: false} // Return an invalid sql.NullTime if parsing fails.
}

// parseToNullString converts a string into a sql.NullString value.
//
// Parameters:
// - s: The string to convert.
//
// Returns:
// - A sql.NullString value, which is valid if the string is non-empty.
func parseToNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// getConfigFilePath constructs the full path to the configuration file.
//
// Returns:
// - The full path to the configuration file in the user's home directory.
func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir() // Retrieve the user's home directory.
	return homeDir + configFileName
}

// Read reads the configuration file and parses it into a Config struct.
//
// Returns:
// - A Config struct populated with data from the configuration file.
func Read() Config {
	// Open the configuration file for reading.
	filepath := getConfigFilePath()
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("error opening file at location: %v\n", filepath)
	}
	defer file.Close()

	var config Config

	// Read the file's content into memory.
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file data")
	}

	// Unmarshal the JSON data into the Config struct.
	if err = json.Unmarshal(data, &config); err != nil {
		fmt.Println("unable to parse data to struct")
	}
	return config
}
