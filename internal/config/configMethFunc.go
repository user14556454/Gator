package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser() error {
	newConfig, err := json.Marshal(c)
	if err != nil {
		return err
	}
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, newConfig, 0644); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(homeDir, ".gatorconfig.json")
	return filePath, nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonSlice, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	var jsonData Config
	if err := json.Unmarshal(jsonSlice, &jsonData); err != nil {
		return Config{}, err
	}
	return jsonData, nil
}
