package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting os home dir: %w", err)
	}

	fullPath := filepath.Join(homedir, configFileName)
	return fullPath, nil
}

func write(c Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(c)
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(*c)
}

// must be called first as it initializes the file path
func (c *Config) Read() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("err on file read: %w", err)
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return fmt.Errorf("err on unmarshal: %w", err)
	}

	return nil
}
