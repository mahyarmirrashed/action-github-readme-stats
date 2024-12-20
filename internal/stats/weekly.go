package stats

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mahyarmirrashed/github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/github-readme-stats/internal/github"
)

const width = 30

var weekdayNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

func GetWeeklyCommitData(
	cfg config.Config,
	commits []github.Commit,
) (string, error) {
	// Load the user's configured timezone
	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return "", fmt.Errorf("failed to load timezone %s: %w", cfg.TimeZone, err)
	}

	// Initialize weekday counters
	weekdays := map[time.Weekday]int{
		time.Monday:    0,
		time.Tuesday:   0,
		time.Wednesday: 0,
		time.Thursday:  0,
		time.Friday:    0,
		time.Saturday:  0,
		time.Sunday:    0,
	}

	for _, commit := range commits {
		commitTime := commit.CommittedDate.In(loc)
		weekdays[commitTime.Weekday()]++
	}

	// Compute total and maximum commits
	var commitsOverAllWeekdaysCount int
	var weekdayWithMostCommitsCount int
	var weekdayWithMostCommits time.Weekday
	for weekday, count := range weekdays {
		commitsOverAllWeekdaysCount += count

		if count > weekdayWithMostCommitsCount {
			weekdayWithMostCommits = weekday
			weekdayWithMostCommitsCount = count
		}
	}

	if commitsOverAllWeekdaysCount == 0 {
		return "no commits found", nil
	}

	var lines []string

	lines = append(lines, fmt.Sprintf("📅 I'm most productive on %ss", weekdayNames[weekdayWithMostCommits]))
	lines = append(lines, "")

	// Generate graph for weekly commit stats
	for i, name := range weekdayNames {
		dayCount := weekdays[time.Weekday(i)]

		// Percent calculation
		absolutePercentageOfCommits := math.Round((float64(dayCount)/float64(commitsOverAllWeekdaysCount))*10000) / 100
		relativePercentageOfCommits := math.Round((float64(dayCount)/float64(weekdayWithMostCommitsCount))*10000) / 100

		// Graph creation
		done := int((relativePercentageOfCommits / 100) * float64(width))
		if done > width {
			done = width
		}
		graph := strings.Repeat("█", done) + strings.Repeat("░", width-done)

		line := fmt.Sprintf("%-12s\t%-4d commits\t%s\t%.2f%%",
			name, dayCount, graph, absolutePercentageOfCommits)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}
