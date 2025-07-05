package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_DBNAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Database.User == "" {
		return fmt.Errorf("DB_USER environment variable is required")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("DB_DBNAME environment variable is required")
	}
	if !strings.HasPrefix(cfg.Server.Port, ":") {
		cfg.Server.Port = ":" + cfg.Server.Port
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
