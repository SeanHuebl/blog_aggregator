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
	conf := config.Read()
	state := config.State{
		ConfigPtr: &conf,
	}
	commands := config.Commands{
		Commands: make(map[string]func(*config.State, config.Command) error),
	}
	commands.Register("login", config.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("cannot run with fewer than two arguments")
		os.Exit(1)
	}
	if len(args) == 2 && (args[1] == "login" || args[1] == "register") {
		fmt.Println("username required")
		os.Exit(1)
	}
	// need some way to distinguish between input commands maybe struct?
	state.ConfigPtr.SetUser(args[2])
	db, err := sql.Open("postgres", state.ConfigPtr.DbUrl)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries

	c := config.Read()
	fmt.Println(c)

}
