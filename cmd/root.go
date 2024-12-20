package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mahyarmirrashed/github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/github-readme-stats/internal/github"
	"github.com/mahyarmirrashed/github-readme-stats/internal/stats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const readmePath = "README.md"

var includes []string

var rootCmd = &cobra.Command{
	Use:   "github-readme-stats",
	Short: "Update GitHub readme statistics",
	Long:  "Update your GitHub README with various statistics such as, what time of day you code, what days of the week you code, and more!",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load and validate configuration
		cfg := config.LoadConfig()
		if cfg.GithubToken == "" {
			return fmt.Errorf("GITHUB_TOKEN not provided")
		}

		// Create GraphQL client
		client := github.NewClient(cfg.GithubToken)
		ctx := context.Background()

		// Fetch repositories from user
		repositories, err := github.FetchRepositories(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to get repositories: %w", err)
		}

		// Fetch commits from all repositories
		commits, err := github.FetchCommitsFromRepositories(ctx, client, repositories)
		if err != nil {
			return fmt.Errorf("failed to get commits: %w", err)
		}

		// Build the output content based on the order of `includes`
		var contentBuilder strings.Builder

		for _, item := range includes {
			switch item {
			case "DAY_STATS":
				contentBuilder.WriteString("\nDay Stats Placeholder\n")
			case "WEEK_STATS":
				// Compute the weekly stats from the repositories and user info
				weeklyData, err := stats.GetWeeklyCommitData(cfg, commits)
				if err != nil {
					return fmt.Errorf("failed to get weekly commit data: %w", err)
				}
				contentBuilder.WriteString("\n```\n" + weeklyData + "\n```\n")

			case "TOP_LANGUAGES":
				contentBuilder.WriteString("\nTop Languages Placeholder\n")
			default:
				// Unknown item, skip or handle error
				contentBuilder.WriteString(fmt.Sprintf("Unknown item: %s\n", item))
			}
		}

		// Append with newline
		contentBuilder.WriteString("\n")

		// Update the README file
		if err := updateReadme(readmePath, contentBuilder.String()); err != nil {
			return fmt.Errorf("failed to update README: %w", err)
		}

		return nil
	},
}

func Execute() {
	rootCmd.SilenceUsage = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rootCmd.Flags().StringArrayVar(&includes, "include", []string{}, "Ordered list of stats to include (e.g. DAY_STATS, WEEK_STATS, TOP_LANGUAGES)")
}

func updateReadme(filepath string, newContent string) error {
	// Read the file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	// Find the block between <!-- README-STATS:START --> and <!-- README-STATS:END -->
	re := regexp.MustCompile("(?s)<!--( ?)README-STATS:START( ?)-->(.*?)<!--( ?)README-STATS:END( ?)-->")
	matches := re.FindSubmatch(data)
	if matches == nil {
		return fmt.Errorf("could not find README-STATS block")
	}

	// Replace the block content
	updatedContent := re.ReplaceAllString(string(data), fmt.Sprintf("<!-- README-STATS:START -->\n%s<!-- README-STATS:END -->", newContent))

	// Write the updated content back to the file
	err = os.WriteFile(filepath, []byte(updatedContent), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write updated README file: %w", err)
	}

	return nil
}
