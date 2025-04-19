package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting os home dir: %w", err)
	}

	return fmt.Sprintf("%s/%s", homedir, configFileName), nil
}

func write(c Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// read write perm
	return os.WriteFile(path, jsonData, 0644)

}

func (c *Config) SetUser(userName string) error {
	fmt.Println(*c)

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
