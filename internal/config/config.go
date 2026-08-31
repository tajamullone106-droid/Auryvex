package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	APIID       int
	APIHash     string
	DatabaseURL string
	RedisURL    string
	WebAppURL   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found, using environment variables")
	}

	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		log.Fatal("invalid API_ID")
	}

	cfg := &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		APIID:       apiID,
		APIHash:     os.Getenv("API_HASH"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		WebAppURL:   os.Getenv("WEBAPP_URL"),
	}

	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN is missing")
	}

	return cfg
}
