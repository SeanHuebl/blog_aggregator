package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/seanhuebl/blog_aggregator/internal/database"
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
	fmt.Println(id)
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
