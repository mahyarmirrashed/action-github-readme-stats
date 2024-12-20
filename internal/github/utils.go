package github

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/rs/zerolog/log"
)

type (
	Repository = GetRepositoriesViewerUserRepositoriesRepositoryConnectionNodesRepository
	Commit     = GetCommitsViewerUserRepositoryDefaultBranchRefTargetCommitHistoryCommitHistoryConnectionNodesCommit
)

func FetchRepositories(ctx context.Context, client graphql.Client) ([]Repository, error) {
	var repositories []Repository

	var repositoryPaginationCursor string
	repositoryPaginationLeft := true

	for repositoryPaginationLeft {
		repositoryQuery, err := GetRepositories(ctx, client, repositoryPaginationCursor)
		if err != nil {
			return nil, err
		}

		for _, repository := range repositoryQuery.Viewer.Repositories.Nodes {
			if !repository.IsEmpty {
				repositories = append(repositories, repository)
			}
		}

		repositoryPaginationLeft = repositoryQuery.Viewer.Repositories.PageInfo.HasNextPage
		repositoryPaginationCursor = repositoryQuery.Viewer.Repositories.PageInfo.EndCursor
	}

	log.Debug().Msgf("Found %d repositories", len(repositories))

	return repositories, nil
}

func FetchCommitsFromRepositories(ctx context.Context, client graphql.Client, repositories []Repository) ([]Commit, error) {
	var commits []Commit

	for _, repository := range repositories {
		log.Info().Msgf("Fetching commits for %s", repository.Name)

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

	log.Debug().Msgf("Found %d commits across all repositories", len(commits))

	return commits, nil
}
