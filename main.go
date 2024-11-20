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
	if len(os.Args) < 2 {
		fmt.Printf("not enough arguments")
		os.Exit(1)
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
		Arguments: os.Args[2:],
	}
	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)
	commands.Register("reset", config.HandlerReset)
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
