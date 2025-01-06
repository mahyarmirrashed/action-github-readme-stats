package github

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

type (
	Repository = GetRepositoriesViewerUserRepositoriesRepositoryConnectionNodesRepository
	Commit     = GetCommitsViewerUserRepositoryDefaultBranchRefTargetCommitHistoryCommitHistoryConnectionNodesCommit
	Language   = GetLanguagesViewerUserRepositoryLanguagesLanguageConnectionEdgesLanguageEdge
)

func FetchRepositories(ctx context.Context, client graphql.Client) ([]Repository, error) {
	var repositories []Repository

	log.Info().Msg("Fetching all repositories...")

	var repositoryPaginationCursor string
	repositoryPaginationLeft := true

	for repositoryPaginationLeft {
		repositoryQuery, err := GetRepositories(ctx, client, repositoryPaginationCursor)
		if err != nil {
			return nil, err
		}

		before_ := len(repositories)

		for _, repository := range repositoryQuery.Viewer.Repositories.Nodes {
			if !repository.IsEmpty {
				repositories = append(repositories, repository)
			}
		}

		log.Debug().Msgf("Added %d repositories", len(repositories)-before_)

		repositoryPaginationLeft = repositoryQuery.Viewer.Repositories.PageInfo.HasNextPage
		repositoryPaginationCursor = repositoryQuery.Viewer.Repositories.PageInfo.EndCursor
	}

	log.Info().Msgf("Found %d repositories", len(repositories))

	return repositories, nil
}

func FetchCommitsFromRepositories(ctx context.Context, client graphql.Client, repositories []Repository) ([]Commit, error) {
	var commits []Commit

	log.Info().Msg("Fetching commits across all repositories...")

	for _, repository := range repositories {
		log.Debug().Msgf("Fetching commits for %s", repository.Name)

		var commitPaginationCursor *string
		commitPaginationLeft := true

		for commitPaginationLeft {
			commitQuery, err := GetCommits(ctx, client, repository.Name, commitPaginationCursor)
			if err != nil {
				return nil, err
			}

			commitTarget, ok := commitQuery.Viewer.Repository.DefaultBranchRef.GetTarget().(*GetCommitsViewerUserRepositoryDefaultBranchRefTargetCommit)
			if !ok {
				log.Warn().Msg("Empty repository found.")
				commitPaginationLeft = false

				continue
			}

			commits = append(commits, commitTarget.History.Nodes...)

			log.Debug().Msgf("Added %d commits from %s", len(commitTarget.History.Nodes), repository.Name)

			commitPaginationLeft = commitTarget.History.PageInfo.HasNextPage
			commitPaginationCursor = &commitTarget.History.PageInfo.EndCursor
		}
	}

	log.Info().Msgf("Found %d commits across all repositories", len(commits))

	return commits, nil
}

func FetchLanguagesFromRepositories(ctx context.Context, client graphql.Client, repositories []Repository) ([]Language, error) {
	var languages []Language

	log.Info().Msg("Fetching used languages across all repositories...")

	for _, repository := range repositories {
		languageQuery, err := GetLanguages(ctx, client, repository.Name)
		if err != nil {
			return nil, err
		}

		languages = append(languages, languageQuery.Viewer.Repository.Languages.Edges...)

		log.Debug().Msgf("Added %d languages from %s", len(languageQuery.Viewer.Repository.Languages.Edges), repository.Name)
	}

	log.Info().Msgf("Found %d languages", len(languages))

	return languages, nil
}
