package config

import (
	"os"
)

type Config struct {
	GithubToken string
}

func LoadConfig() Config {
	return Config{
		GithubToken: os.Getenv("GITHUB_TOKEN"),
	}
}
