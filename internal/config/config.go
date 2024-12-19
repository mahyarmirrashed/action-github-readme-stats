package config

import (
	"log"
	"os"
)

type Config struct {
	GithubToken string
}

func LoadConfig() Config {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN not provided")
	}

	return Config{
		GithubToken: token,
	}
}
