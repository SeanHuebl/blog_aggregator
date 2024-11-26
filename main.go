package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver for database interaction
	"github.com/seanhuebl/blog_aggregator/internal/config"
	"github.com/seanhuebl/blog_aggregator/internal/database"
)

// main is the entry point of the blog aggregator application.
// It processes command-line arguments, initializes the application state,
// establishes a database connection, and executes the requested command.
func main() {
	var arguments []string

	// Ensure at least one command-line argument is provided
	if len(os.Args) < 2 {
		fmt.Printf("not enough arguments")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		// Capture additional arguments beyond the first
		arguments = os.Args[2:]
	} else {
		arguments = nil
	}

	// Load application configuration from the environment or config file
	conf := config.Read()

	// Initialize application state
	state := config.State{
		ConfigPtr: &conf, // Link configuration to state
	}

	// Initialize a command registry
	commands := &config.Commands{
		Commands: make(map[string]func(*config.State, config.Command) error),
	}

	// Parse the command from the first argument
	command := config.Command{
		Name:      os.Args[1],
		Arguments: arguments,
	}

	// Register available commands and their handlers
	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)
	commands.Register("reset", config.HandlerReset)
	commands.Register("users", config.HandlerGetUsers)
	commands.Register("agg", config.HandlerAgg)
	commands.Register("feeds", config.HandlerFeeds)
	commands.Register("addfeed", config.MiddlewareLoggedIn(config.HandlerAddFeed))
	commands.Register("follow", config.MiddlewareLoggedIn(config.HandlerFollow))
	commands.Register("following", config.MiddlewareLoggedIn(config.HandlerFollowing))
	commands.Register("unfollow", config.MiddlewareLoggedIn(config.HandlerUnfollow))
	commands.Register("browse", config.MiddlewareLoggedIn(config.HandlerBrowse))

	// Establish a database connection using the provided configuration
	db, err := sql.Open("postgres", state.ConfigPtr.DbUrl)
	if err != nil {
		// Exit if the database connection cannot be established
		println(err)
		os.Exit(1)
	}
	// Initialize database queries for application use
	dbQueries := database.New(db)
	state.Db = dbQueries

	// Execute the requested command
	if err := commands.Run(&state, command); err != nil {
		fmt.Println(err)
		os.Exit(1) // Exit with error if the command fails
	}
}
