package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	AppPort   string
	RedisURL  string
	RateLimit string
}

func (c *Config) GetRateLimit() int {
	limit, err := strconv.Atoi(c.RateLimit)
	if err != nil {
		log.Fatalf("invalid RATE_LIMIT: %v", err)
	}

	return limit
}

func Load() *Config {
	cfg := &Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		RedisURL:  getEnv("REDIS_URL", "localhost:6379"),
		RateLimit: getEnv("RATE_LIMIT", "10"),
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		if fallback == "" {
			log.Fatalf("missing required env: %s", key)
		}

		return fallback
	}

	return value
}
