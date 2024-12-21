package stats

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/mahyarmirrashed/github-readme-stats/internal/config"
	"github.com/mahyarmirrashed/github-readme-stats/internal/github"
)

type Language struct {
	Name  string
	Count int
}

func GetLanguageData(
	cfg config.Config,
	languages []github.Language,
) (string, error) {
	languageUsage := make(map[string]int)

	// Count language usage by name
	for _, lang := range languages {
		languageName := lang.Node.Name
		languageUsage[languageName]++
	}

	// Convert map to slice for sorting
	var aggregatedLanguages []Language
	var languageUsageOverAllReposCount int
	for name, count := range languageUsage {
		languageUsageOverAllReposCount += count

		aggregatedLanguages = append(aggregatedLanguages, Language{
			Name:  name,
			Count: count,
		})
	}

	if len(aggregatedLanguages) == 0 || languageUsageOverAllReposCount == 0 {
		return "no languages found", nil
	}

	// Sort aggregated languages
	sort.Slice(aggregatedLanguages, func(i, j int) bool {
		if aggregatedLanguages[i].Count != aggregatedLanguages[j].Count {
			// Sort by count descending if counts are not equal
			return aggregatedLanguages[i].Count > aggregatedLanguages[j].Count
		}
		// Sort by name ascending, otherwise
		return aggregatedLanguages[i].Name < aggregatedLanguages[j].Name
	})

	// Build output
	var lines []string
	topLanguage := aggregatedLanguages[0].Name
	lines = append(lines, fmt.Sprintf("🧪 %s for the win!", topLanguage))
	lines = append(lines, "")

	if len(aggregatedLanguages) > 5 {
		aggregatedLanguages = aggregatedLanguages[:5]
	}
	for _, lang := range aggregatedLanguages {
		absolutePercentageOfLanguageUsage := math.Round((float64(lang.Count)/float64(languageUsageOverAllReposCount))*10000) / 100
		relativePercentageOfLanguageUsage := math.Round((float64(lang.Count)/float64(aggregatedLanguages[0].Count))*10000) / 100

		// Graph creation
		done := int((relativePercentageOfLanguageUsage / 100) * float64(width))
		if done > width {
			done = width
		}
		graph := strings.Repeat("█", done) + strings.Repeat("░", width-done)

		line := fmt.Sprintf("%-12s\t%-6d repos\t%s\t%.2f%%",
			lang.Name, lang.Count, graph, absolutePercentageOfLanguageUsage)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}
