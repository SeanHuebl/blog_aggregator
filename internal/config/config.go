package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + configFileName
}

func (c *Config) SetUser(username string) {
	config := Read()
	config.CurrentUserName = username
	// need to marshal data to []byte then write to the json file

}

func Read() Config {

	filepath := getConfigFilePath()
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening file at location: %v\n", filepath)
	}
	defer file.Close()
	var config Config
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file data")
	}

	if err = json.Unmarshal(data, &config); err != nil {
		fmt.Println("Unable to parse data to struct")
	}
	return config
}
