package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLDatabase     string `mapstructure:"MYSQL_DATABASE"`
	MySQLUser         string `mapstructure:"MYSQL_USER"`
	MySQLPassword     string `mapstructure:"MYSQL_PASSWORD"`
	MySQLRootPassword string `mapstructure:"MYSQL_ROOT_PASSWORD"`
	DSN               string `mapstructure:"DSN"`
	AutoMigrate       bool   `mapstructure:"AUTO_MIGRATE"`
}

func New(path string) *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env not found, reading configuration from ENV")
	}

	config := loadConfig()
	return &config
}

func loadConfig() Config {
	cfg := Config{
		MySQLDatabase:     GetString("MYSQL_DATABASE", ""),
		MySQLUser:         GetString("MYSQL_USER", ""),
		MySQLPassword:     GetString("MYSQL_PASSWORD", ""),
		MySQLRootPassword: GetString("MYSQL_ROOT_PASSWORD", ""),
		DSN:               GetString("DSN", ""),
		AutoMigrate:       GetBool("AUTO_MIGRATE", false),
	}
	return cfg
}
