package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Config struct {
	APIKey string `json:"api_key"`
}

var AppConfig Config

func LoadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		color.New(color.FgRed).Printf("Configuration file not found: %v\n", err)
		return fmt.Errorf("configuration file not found: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		color.New(color.FgRed).Printf("Error decoding configuration file: %v\n", err)
		return fmt.Errorf("error decoding configuration file: %v", err)
	}

	color.New(color.FgGreen).Println("Configuration loaded successfully.")
	return nil
}
