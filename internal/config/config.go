package config

import (
	"os"
)

type Config struct {
	GithubToken string
	TimeZone    string
}

func LoadConfig() Config {
	return Config{
		GithubToken: os.Getenv("GITHUB_TOKEN"),
		TimeZone:    os.Getenv("TIMEZONE"),
	}
}
