package stats

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mahyarmirrashed/action-github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/action-github-readme-stats/internal/github"
)

func GetDailyCommitData(
	cfg config.Config,
	commits []github.Commit,
) (string, error) {
	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return "", fmt.Errorf("failed to load timezone %s: %w", cfg.TimeZone, err)
	}

	var (
		morningCount int // 4am–10am
		daytimeCount int // 10am–4pm
		eveningCount int // 4pm–10pm
		nightCount   int // 10pm–4am
	)

	// Classify each commit into a time-of-day category
	for _, commit := range commits {
		commitTime := commit.CommittedDate.In(loc)
		hour := commitTime.Hour()

		switch {
		case hour >= 22 || hour < 4:
			nightCount++
		case hour >= 4 && hour < 10:
			morningCount++
		case hour >= 10 && hour < 16:
			daytimeCount++
		case hour >= 16 && hour < 22:
			eveningCount++
		}
	}

	// Calculate the total commits and find the category with the most commits
	commitsOverAllWorkingHoursCount := morningCount + daytimeCount + eveningCount + nightCount
	if commitsOverAllWorkingHoursCount == 0 {
		return "no commits found", nil
	}

	// Determine which category has the most commits
	preferredWorkingHourCount := nightCount
	preferredWorkingHour := "Night"
	if morningCount > preferredWorkingHourCount {
		preferredWorkingHourCount = morningCount
		preferredWorkingHour = "Morning"
	}
	if daytimeCount > preferredWorkingHourCount {
		preferredWorkingHourCount = daytimeCount
		preferredWorkingHour = "Daytime"
	}
	if eveningCount > preferredWorkingHourCount {
		preferredWorkingHourCount = eveningCount
		preferredWorkingHour = "Evening"
	}

	// Prepare data for output
	categories := []struct {
		Icon  string
		Name  string
		Count int
	}{
		{"🌞", "Morning", morningCount},
		{"🌆", "Daytime", daytimeCount},
		{"🌃", "Evening", eveningCount},
		{"🌙", "Night", nightCount},
	}

	// Build the output
	var lines []string
	lines = append(lines, fmt.Sprintf("🕰️ I get my jam on during the %s!", strings.ToLower(preferredWorkingHour)))
	lines = append(lines, "")

	width := 30
	for _, category := range categories {
		// Percent calculation
		absolutePercentageOfCommits := math.Round((float64(category.Count)/float64(commitsOverAllWorkingHoursCount))*10000) / 100
		relativePercentageOfCommits := math.Round((float64(category.Count)/float64(preferredWorkingHourCount))*10000) / 100

		// Graph creation
		done := int((relativePercentageOfCommits / 100) * float64(width))
		if done > width {
			done = width
		}
		graph := strings.Repeat("█", done) + strings.Repeat("░", width-done)

		line := fmt.Sprintf("%s %-9s\t%-6d commits\t%s\t%.2f%%",
			category.Icon, category.Name, category.Count, graph, absolutePercentageOfCommits)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}
