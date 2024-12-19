package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/mahyarmirrashed/github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/github-readme-stats/internal/github"
	"github.com/spf13/cobra"
)

const readmePath = "README.md"

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

		// Fetch user information
		rep, err := github.GetUserInfo(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		// Generate the new content to insert into README
		newContent := fmt.Sprintf("\nUser Login: %s\n", rep.Viewer.Login)

		// Update the README file
		if err := updateReadme(readmePath, newContent); err != nil {
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
	updatedContent := re.ReplaceAllString(string(data), fmt.Sprintf("<!-- README-STATS:START -->\n%s\n<!-- README-STATS:END -->", newContent))

	// Write the updated content back to the file
	err = os.WriteFile(filepath, []byte(updatedContent), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write updated README file: %w", err)
	}

	return nil
}
