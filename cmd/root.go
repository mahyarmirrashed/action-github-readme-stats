package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mahyarmirrashed/github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/github-readme-stats/internal/github"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github-readme-stats",
	Short: "Update GitHub readme statistics",
	Long:  "Update your GitHub README with various statistics such as, what time of day you code, what days of the week you code, and more!",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.LoadConfig()
		if cfg.GithubToken == "" {
			return fmt.Errorf("GITHUB_TOKEN not provided")
		}

		client := github.NewClient(cfg.GithubToken)
		ctx := context.Background()

		rep, err := github.GetUserInfo(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		fmt.Printf("User Login: %s\n", rep.Viewer.Login)
		fmt.Printf("User ID: %s\n", rep.Viewer.Id)
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
