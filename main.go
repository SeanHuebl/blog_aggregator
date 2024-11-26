package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/seanhuebl/blog_aggregator/internal/config"
	"github.com/seanhuebl/blog_aggregator/internal/database"
)

func main() {
	var arguments []string
	if len(os.Args) < 2 {
		fmt.Printf("not enough arguments")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		arguments = os.Args[2:]
	} else {
		arguments = nil
	}
	conf := config.Read()
	state := config.State{
		ConfigPtr: &conf,
	}
	commands := &config.Commands{
		Commands: make(map[string]func(*config.State, config.Command) error),
	}
	command := config.Command{
		Name:      os.Args[1],
		Arguments: arguments,
	}
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

	db, err := sql.Open("postgres", state.ConfigPtr.DbUrl)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries

	if err := commands.Run(&state, command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
