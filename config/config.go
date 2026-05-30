package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration.
// Loaded once at startup in main.go and passed everywhere via dependency injection.
// os.Getenv is called only in this package — nothing else reaches for env vars directly.
type Config struct {
	Port      string
	Env       string
	LogLevel  string
	LogFormat string // "json" | "text"
	RedisURL  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:      getEnv("SERVER_PORT", "8080"),
		Env:       getEnv("APP_ENV", "development"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		RedisURL:  getEnv("REDIS_URL", "redis://localhost:6379"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
