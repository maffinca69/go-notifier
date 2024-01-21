package infrastructure

import (
	"encoding/json"
	"os"
)

type RepositoryConfig struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type TelegramConfig struct {
	ChatID            int64  `json:"chat_id"`
	NewVersionMessage string `json:"new_version_message"`
}

type Config struct {
	Repository     []RepositoryConfig `json:"repositories"`
	TelegramConfig TelegramConfig     `json:"telegram"`
	CronExpression string             `json:"cron_expression"`
}

const ConfigName = "config.json"

var configInstance *Config

func GetConfig() *Config {
	if configInstance == nil {
		configInstance = setupConfig()
	}

	return configInstance
}

func setupConfig() *Config {
	content, _ := os.ReadFile(ConfigName)
	payload := &Config{}

	if err := json.Unmarshal(content, &payload); err != nil {
		panic("Error load config " + ConfigName)
	}
	return payload
}
