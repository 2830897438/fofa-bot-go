package config

import (
	"encoding/json"
	"os"
)

// Config represents the bot configuration
type Config struct {
	BotToken   string   `json:"bot_token"`
	APIs       []string `json:"apis"`
	Admins     []int64  `json:"admins"`
	Proxy      string   `json:"proxy"`
	FullMode   bool     `json:"full_mode"`
	PublicMode bool     `json:"public_mode"`
	Presets    []Preset `json:"presets"`
	UpdateURL  string   `json:"update_url"`
}

// Preset represents a query preset
type Preset struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Load loads configuration from file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves configuration to file
func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// IsAdmin checks if a user is an admin
func (c *Config) IsAdmin(userID int64) bool {
	for _, admin := range c.Admins {
		if admin == userID {
			return true
		}
	}
	return false
}
