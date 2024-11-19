package main

import (
	"fmt"
	"os"

	"github.com/seanhuebl/blog_aggregator/internal/config"
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
	if len(args) == 2 && args[1] == "login" {
		fmt.Println("username required")
		os.Exit(1)
	}
	state.ConfigPtr.SetUser(args[2])
	c := config.Read()
	fmt.Println(c)

}
