package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	config := Read()
	config.CurrentUserName = username
	jdata, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}
	filepath := getConfigFilePath()
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file at location %v", filepath)
	}
	defer file.Close()
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
	f, err := c.Commands[cmd.name]
	if !err {
		f(s, cmd)
		return nil
	}

	return fmt.Errorf("command not found")
}

type Command struct {
	name      string
	arguments []string
}

type State struct {
	ConfigPtr *Config
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("login expects one argument: username")
	}

	err := s.ConfigPtr.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("%s", err)
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
