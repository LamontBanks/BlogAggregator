// Handles reading/writing the config file for handling the logged in user and database connection details
package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

const (
	userConfigFileName = ".gatorconfig.json"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read tyeh JSON from the ~/.gatorconfig.json
// (must be located at the home directory)
func ReadConfig() (Config, error) {
	var config Config

	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	log.Printf("Reading config data from: %v", configFile)

	// Read bytes
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if len(configFileBytes) == 0 {
		return Config{}, errors.New("Config file empty?: " + configFile)
	}

	// Unmarshal bytes into struct
	if err := json.Unmarshal(configFileBytes, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// Saves the config file
func (cfg *Config) SetConfig() error {
	// Get file path
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Marshal to convert from JSON to bytes
	configFileBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// Wtite bytes
	err = os.WriteFile(configFile, configFileBytes, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	var homeDir string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return homeDir, err
	}
	homeDir += "/" + userConfigFileName

	return homeDir, nil
}
