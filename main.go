package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mahyarmirrashed/github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/github-readme-stats/internal/github"
	"github.com/mahyarmirrashed/github-readme-stats/internal/stats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "golang.org/x/crypto/x509roots/fallback" // CA bundle for FROM Scratch
)

func main() {
	// Initialize logger with console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // Set default log level to Info

	// Fetch arguments excluding the program name
	if len(os.Args) != 2 {
		log.Error().Msg("Invalid number of arguments provided")
		printUsage()
		os.Exit(1)
	}
	argument := os.Args[1]

	// Validate and process features
	features, err := validate(argument)
	if err != nil {
		log.Error().Err(err).Msg("Invalid arguments provided")
		printUsage()
		os.Exit(1)
	}

	log.Debug().Msgf("Features to include: %v", features)

	// Load and validate configuration
	cfg := config.LoadConfig()
	if cfg.GithubToken == "" {
		log.Error().Msg("GITHUB_TOKEN not provided in configuration")
		os.Exit(1)
	}

	log.Debug().Msgf("Timezone: %s", cfg.TimeZone)

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current working directory")
		os.Exit(1)
	}
	readmePath := filepath.Join(cwd, "README.md")

	log.Debug().Msgf("Path for README is: %s", readmePath)

	// Create GitHub client
	client := github.NewClient(cfg.GithubToken)
	ctx := context.Background()

	// Fetch repositories from user
	repositories, err := github.FetchRepositories(ctx, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get repositories from GitHub")
		os.Exit(1)
	}

	// Fetch commits from all repositories
	commits, err := github.FetchCommitsFromRepositories(ctx, client, repositories)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get commits from repositories")
		os.Exit(1)
	}

	// Fetch languages from all repositories
	languages, err := github.FetchLanguagesFromRepositories(ctx, client, repositories)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get languages from repositories")
		os.Exit(1)
	}

	// Build the output content based on the order of features
	var contentBuilder strings.Builder
	codeBlock := func(content string) string { return "\n```\n" + content + "\n```\n" }

	for _, feature := range features {
		switch feature {
		case "DAY_STATS":
			log.Info().Msg("Calculating commit statistics based on time of day")
			dailyStats, err := stats.GetDailyCommitData(cfg, commits)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get daily commit stats")
				os.Exit(1)
			}
			contentBuilder.WriteString(codeBlock(dailyStats))

		case "WEEK_STATS":
			log.Info().Msg("Calculating commit statistics based on day of week")
			weeklyStats, err := stats.GetWeeklyCommitData(cfg, commits)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get weekly commit stats")
				os.Exit(1)
			}
			contentBuilder.WriteString(codeBlock(weeklyStats))

		case "LANGUAGE_STATS":
			log.Info().Msg("Calculating language statistics")
			languageStats, err := stats.GetLanguageData(cfg, languages)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get language stats")
				os.Exit(1)
			}
			contentBuilder.WriteString(codeBlock(languageStats))

		default:
			// Unknown feature, skip or handle error
			log.Warn().Msgf("Unknown feature: %s", feature)
			contentBuilder.WriteString(fmt.Sprintf("\n\nUnknown feature: %s\n", feature))
		}
	}

	// Append with newline
	contentBuilder.WriteString("\n")

	// Update the README file
	if err := updateReadme(readmePath, contentBuilder.String()); err != nil {
		log.Error().Err(err).Msg("Failed to update README.md")
		os.Exit(1)
	}

	log.Info().Msg("README.md successfully updated")
}

// Validate the provided features
func validate(arg string) ([]string, error) {
	validFeatures := map[string]bool{
		"DAY_STATS":      true,
		"WEEK_STATS":     true,
		"LANGUAGE_STATS": true,
	}

	var features []string
	for _, feature := range strings.Split(arg, ",") {
		if !validFeatures[strings.ToUpper(strings.TrimSpace(feature))] {
			return nil, fmt.Errorf("invalid feature argument: %s", feature)
		}
		features = append(features, feature)
	}
	return features, nil
}

// Replace README-STATS block with new content
func updateReadme(filepath string, newContent string) error {
	// Read the file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read README.md file: %w", err)
	}

	// Define the regex to find the block between <!-- README-STATS:START --> and <!-- README-STATS:END -->
	re := regexp.MustCompile(`(?s)<!--\s*README-STATS:START\s*-->(.*?)<!--\s*README-STATS:END\s*-->`)
	if !re.Match(data) {
		return fmt.Errorf("could not find README-STATS block in README.md")
	}

	// Replace the block content
	updatedContent := re.ReplaceAllString(string(data), fmt.Sprintf("<!-- README-STATS:START -->\n%s<!-- README-STATS:END -->", newContent))

	// Write the updated content back to the file
	err = os.WriteFile(filepath, []byte(updatedContent), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write updated README.md file: %w", err)
	}

	return nil
}

// Displays the correct usage of the program
func printUsage() {
	usage := `
Usage: github-readme-stats <FEATURE_LIST>

FEATURE_LIST: Comma-separated list of features to include (e.g., DAY_STATS,WEEK_STATS,LANGUAGE_STATS)

Example:
  github-readme-stats DAY_STATS,WEEK_STATS,LANGUAGE_STATS
`
	fmt.Println(usage)
}
