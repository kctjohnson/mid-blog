package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePath string `mapstructure:"DATABASE_PATH"`
	AutoMigrate  bool   `mapstructure:"AUTO_MIGRATE"`
	OpenAIKey    string `mapstructure:"OPENAI_KEY"`
}

func New(path string) *Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("File .env not found, reading configuration from ENV")
	}

	config := loadConfig()
	return &config
}

func loadConfig() Config {
	cfg := Config{
		DatabasePath: GetString("DATABASE_PATH", "data/midblog.db"),
		AutoMigrate:  GetBool("AUTO_MIGRATE", true),
		OpenAIKey:    GetString("OPENAI_KEY", ""),
	}

	return cfg
}
