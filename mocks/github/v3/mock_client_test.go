package repository

import "testing"

func TestGitHubMockGeneration(t *testing.T) {
	t.Run("GitHub Client Repositories.Get(owner, repo)",
		func(t *testing.T) {
			MockGitHubClientRepositoryGet(t)
		})
	t.Run("GitHub Client Repositories.List(owner)",
		func(t *testing.T) {
			MockGitHubClientRepositoryList(t)
		})
	t.Run("GitHub Client Repository.ListReleases(owner, repo)",
		func(t *testing.T) {
			MockGitHubClientReleaseList(t)
		})
}
